package providers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type DotaProvider struct {
	TeamID string
}

func (d *DotaProvider) GetMatches() (error, []Match) {
	err, htmlResponse := fetchMatchesHTML(d.TeamID)
	if err != nil {
		return err, nil
	}

	err, matches := parseMatchesHTML(htmlResponse)
	if err != nil {
		return err, nil
	}

	return nil, matches
}

func fetchMatchesHTML(teamId string) (error, string) {
	form := url.Values{
		"ajax":     {"block_matches_search"},
		"rid":      {"matches"},
		"data[t1]": {teamId},
		"game":     {"dota-2"},
	}

	req, err := http.PostForm("https://ggscore.com/en/dota-2", form)

	if err != nil {
		return fmt.Errorf("error while sending request: %w", err), ""
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		return fmt.Errorf("error while reading response body: %w", err), ""
	}

	return nil, string(body)
}

func parseMatchesHTML(htmlResponse string) (error, []Match) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlResponse))

	if err != nil {
		return fmt.Errorf("error while parsing HTML: %w", err), nil
	}

	var matches []Match

	doc.Find("tr.m-item[data-bm='gg']").Each(func(i int, s *goquery.Selection) {
		match := Match{}

		match.Id = s.AttrOr("data-id", string(time.Now().Unix()))
		matchTimeStr := s.Find("td.tdate time").AttrOr("data-time", "")
		match.Time, err = time.Parse("2006-01-02 15:04:05", matchTimeStr)
		if err != nil {
			match.Time = time.Now().UTC()
			match.IsLive = true
		}

		match.Tournament, _ = s.Find("td.tname a.tour-pop").Attr("title")
		match.Score = s.AttrOr("data-score", "N/A")
		match.Team1 = s.Find("span.tn1").Text()
		match.Team2 = s.Find("span.tn2").Text()

		matches = append(matches, match)
	})

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Time.Before(matches[j].Time)
	})

	return nil, matches
}

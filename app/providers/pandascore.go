package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const BaseProviderUrl = "https://api.pandascore.co/dota2/matches"

type PandaScoreProvider struct {
	TeamID string
	ApiKey string
}

type TeamMatch struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Opponents []struct {
		Opponent struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"opponent"`
	} `json:"opponents"`
	League struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"league"`
	BeginAt time.Time `json:"begin_at"`
	Results []struct {
		Score  int `json:"score"`
		TeamID int `json:"team_id"`
	}
	StreamsList []struct {
		Language string `json:"language"`
		RawURL   string `json:"raw_url"`
	} `json:"streams_list"`
}

func (p *PandaScoreProvider) GetMatches() (error, []Match) {
	url := fmt.Sprintf("%s?filter[opponent_id]=%s", BaseProviderUrl, p.TeamID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err), nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.ApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while sending request: %w", err), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while reading response body: %w", err), nil
	}

	var teamMatches []TeamMatch
	err = json.Unmarshal(body, &teamMatches)
	if err != nil {
		return fmt.Errorf("error while parsing response body: %w", err), nil
	}

	return nil, ProcessMatches(teamMatches)
}

func ProcessMatches(providerMatches []TeamMatch) []Match {
	matches := make([]Match, 0)

	for _, providerMatch := range providerMatches {
		var match Match

		match.Id = strconv.Itoa(providerMatch.Id)
		match.Tournament = providerMatch.League.Name
		match.Team1 = providerMatch.Opponents[0].Opponent.Name
		match.Team2 = providerMatch.Opponents[1].Opponent.Name
		match.Time = providerMatch.BeginAt
		match.IsLive = providerMatch.Status == "running"
		match.Score = fmt.Sprintf("(%d-%d)", providerMatch.Results[0].Score, providerMatch.Results[1].Score)
		match.URL = providerMatch.StreamsList[0].RawURL
		for _, stream := range providerMatch.StreamsList {
			if stream.Language == "ru" {
				match.URL = stream.RawURL
			}
		}

		matches = append(matches, match)
	}

	return matches
}

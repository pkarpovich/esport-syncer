package providers

import (
	"encoding/json"
	"fmt"
	"github.com/ybbus/httpretry"
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

type Stream struct {
	Language string `json:"language"`
	RawURL   string `json:"raw_url"`
}

type Result struct {
	Score  int `json:"score"`
	TeamID int `json:"team_id"`
}

type Opponent struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type OpponentList struct {
	Opponent Opponent `json:"opponent"`
}

type League struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TeamMatch struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Status       string         `json:"status"`
	OpponentList []OpponentList `json:"opponents"`
	League       League         `json:"league"`
	BeginAt      time.Time      `json:"begin_at"`
	ScheduledAt  time.Time      `json:"scheduled_at"`
	ModifiedAt   time.Time      `json:"modified_at"`
	BestOf       int            `json:"number_of_games"`
	Results      []Result       `json:"results"`
	StreamsList  []Stream       `json:"streams_list"`
}

func (p *PandaScoreProvider) GetMatches() ([]Match, error) {
	url := fmt.Sprintf("%s?filter[opponent_id]=%s", BaseProviderUrl, p.TeamID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.ApiKey))

	client := httpretry.NewDefaultClient(
		httpretry.WithMaxRetryCount(5),
		httpretry.WithBackoffPolicy(func(attemptCount int) time.Duration {
			return time.Duration(attemptCount) * time.Minute
		}),
		httpretry.WithRetryPolicy(func(statusCode int, err error) bool {
			return err != nil || statusCode != http.StatusOK
		}),
	)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while sending request: status code: %d", resp.StatusCode)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("error while closing response body: %v", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}

	var teamMatches []TeamMatch
	err = json.Unmarshal(body, &teamMatches)
	if err != nil {
		bodyStr := string(body)
		return nil, fmt.Errorf("error while unmarshalling response body: %w, body: %s", err, bodyStr)
	}

	return ProcessMatches(teamMatches), nil
}

func ProcessMatches(providerMatches []TeamMatch) []Match {
	matches := make([]Match, 0)

	for _, providerMatch := range providerMatches {
		var match Match

		team1, team2 := getTeams(providerMatch.OpponentList)

		match.Id = strconv.Itoa(providerMatch.Id)
		match.Tournament = providerMatch.League.Name
		match.Team1 = team1
		match.Team2 = team2
		match.BestOf = providerMatch.BestOf
		match.Time = providerMatch.ScheduledAt
		match.ModifiedAt = providerMatch.ModifiedAt
		match.IsLive = providerMatch.Status == "running"
		match.Score = getScore(providerMatch.Results)
		match.URL = getStreamURL(providerMatch.StreamsList)

		matches = append(matches, match)
	}

	return matches
}

func getStreamURL(streams []Stream) string {
	for _, stream := range streams {
		if stream.Language == "ru" {
			return stream.RawURL
		}
	}

	if len(streams) > 0 {
		return streams[0].RawURL
	}

	return ""
}

func getScore(results []Result) string {
	if len(results) == 2 {
		return fmt.Sprintf("(%d-%d)", results[0].Score, results[1].Score)
	}

	return "(0-0)"
}

func getTeams(opponents []OpponentList) (string, string) {
	if len(opponents) == 2 {
		return opponents[0].Opponent.Name, opponents[1].Opponent.Name
	}

	return opponents[0].Opponent.Name, "TBD"
}

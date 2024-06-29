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
	BeginAt     time.Time `json:"begin_at"`
	ScheduledAt time.Time `json:"scheduled_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	BestOf      int       `json:"number_of_games"`
	Results     []struct {
		Score  int `json:"score"`
		TeamID int `json:"team_id"`
	}
	StreamsList []struct {
		Language string `json:"language"`
		RawURL   string `json:"raw_url"`
	} `json:"streams_list"`
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

	defer resp.Body.Close()
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

		match.Id = strconv.Itoa(providerMatch.Id)
		match.Tournament = providerMatch.League.Name
		match.Team1 = providerMatch.Opponents[0].Opponent.Name

		match.Team2 = "TBD"
		if len(providerMatch.Opponents) == 2 {
			match.Team2 = providerMatch.Opponents[1].Opponent.Name
		}

		match.BestOf = providerMatch.BestOf
		match.Time = providerMatch.ScheduledAt
		match.ModifiedAt = providerMatch.ModifiedAt
		match.IsLive = providerMatch.Status == "running"

		match.Score = "(0-0)"
		if len(providerMatch.Results) == 2 {
			match.Score = fmt.Sprintf("(%d-%d)", providerMatch.Results[0].Score, providerMatch.Results[1].Score)
		}

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

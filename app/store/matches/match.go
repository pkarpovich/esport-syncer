package match

import (
	"github.com/pkarpovich/esport-syncer/app/database"
	"log"
	"time"
)

type Match struct {
	Id         string    `json:"id"`
	Tournament string    `json:"tournament"`
	Team1      string    `json:"team1"`
	Team2      string    `json:"team2"`
	Score      string    `json:"score"`
	Time       time.Time `json:"time"`
	BestOf     int       `json:"bestOf"`
	Location   string    `json:"location"`
	URL        string    `json:"url"`
	IsLive     bool      `json:"isLive"`
	TeamId     int       `json:"targetTeamId"`
	GameType   string    `json:"gameType"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type Repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) (*Repository, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS events (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		tournament TEXT,
    		team1 TEXT,
    		team2 TEXT,
    		score TEXT,
    		best_of INTEGER,
    		time TIMESTAMP,
    		location TEXT,
    		url TEXT,
    		is_live BOOLEAN,
    		team_id INTEGER
            game_type TEXT
            modified_at TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CreateOrReplace(event Match) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO events (
                    id,
                    tournament,
                    team1,
                    team2,
                    score,
                    best_of,
                    time,
                    location,
                    url,
                    is_live,
					team_id,
					game_type,
					modified_at
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Id,
		event.Tournament,
		event.Team1,
		event.Team2,
		event.Score,
		event.BestOf,
		event.Time,
		event.Location,
		event.URL,
		event.IsLive,
		event.TeamId,
		event.GameType,
		event.ModifiedAt,
	)

	return err
}

func (r *Repository) GetAll() ([]Match, error) {
	rows, err := r.db.Query(`SELECT
					id,
					tournament,
					team1,
					team2,
					score,
					best_of,
					time,
					location,
					url,
					is_live,
    				team_id,
    				game_type,
					modified_at
					FROM events`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	var events []Match
	for rows.Next() {
		var event Match
		err := rows.Scan(
			&event.Id,
			&event.Tournament,
			&event.Team1,
			&event.Team2,
			&event.Score,
			&event.BestOf,
			&event.Time,
			&event.Location,
			&event.URL,
			&event.IsLive,
			&event.TeamId,
			&event.GameType,
			&event.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repository) GetByTeamId(teamID int, gameType string) ([]Match, error) {
	rows, err := r.db.Query(`SELECT
					id,
					tournament,
					team1,
					team2,
					score,
					best_of,
					time,
					location,
					url,
					is_live,
					team_id,
					game_type,
					modified_at
					FROM events WHERE team_id = ? AND game_type = ?`, teamID, gameType)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("[ERROR] error while closing rows: %v", err)
		}
	}()

	events := make([]Match, 0)

	for rows.Next() {
		var event Match

		err := rows.Scan(
			&event.Id,
			&event.Tournament,
			&event.Team1,
			&event.Team2,
			&event.Score,
			&event.BestOf,
			&event.Time,
			&event.Location,
			&event.URL,
			&event.IsLive,
			&event.TeamId,
			&event.GameType,
			&event.ModifiedAt,
		)
		if err != nil {
			log.Printf("[ERROR] error while scanning rows: %v", err)
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] error while iterating rows: %v", err)
		return nil, err
	}

	return events, nil
}

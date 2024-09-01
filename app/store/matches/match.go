package match

import (
	"github.com/pkarpovich/esport-syncer/app/database"
	"github.com/pkarpovich/esport-syncer/app/utils"
	"log"
	"time"
)

type Match struct {
	Id         string    `json:"id"`
	Tournament string    `json:"tournament"`
	Team1      Team      `json:"team1"`
	Team1Id    int       `json:"team1_id"`
	Team2      Team      `json:"team2"`
	Team2Id    int       `json:"team2_id"`
	Score      string    `json:"score"`
	Time       time.Time `json:"time"`
	BestOf     int       `json:"bestOf"`
	Location   string    `json:"location"`
	URL        string    `json:"url"`
	IsLive     bool      `json:"isLive"`
	GameType   string    `json:"gameType"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type Team struct {
	Id      int    `json:"id"`
	Acronym string `json:"acronym"`
	Name    string `json:"name"`
	Logo    string `json:"logo"`
}

type Repository struct {
	db *database.Client
}

func NewRepository(db *database.Client) (*Repository, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS teams (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			acronym TEXT,
			name TEXT,
			logo TEXT
	)`)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS events (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
			tournament TEXT,
			team1_id INTEGER,
			team2_id INTEGER,
			score TEXT,
			best_of INTEGER,
			time TIMESTAMP,
			location TEXT,
			url TEXT,
			is_live BOOLEAN,
			game_type TEXT,
			modified_at TIMESTAMP,
			FOREIGN KEY(team1_id) REFERENCES teams(id),
			FOREIGN KEY(team2_id) REFERENCES teams(id)
	)`)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CreateOrReplace(event Match) error {
	err := r.CreateOrReplaceTeam(event.Team1)
	if err != nil {
		return err
	}

	err = r.CreateOrReplaceTeam(event.Team2)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`INSERT OR REPLACE INTO events (
                    id,
                    tournament,
                    team1_id,
                    team2_id,
                    score,
                    best_of,
                    time,
                    location,
                    url,
                    is_live,
					game_type,
					modified_at
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.Id,
		event.Tournament,
		event.Team1.Id,
		event.Team2.Id,
		event.Score,
		event.BestOf,
		event.Time,
		event.Location,
		event.URL,
		event.IsLive,
		event.GameType,
		event.ModifiedAt,
	)

	return err
}

func (r *Repository) CreateOrReplaceTeam(team Team) error {
	_, err := r.db.Exec(`INSERT OR REPLACE INTO teams (
			id,
			acronym,
			name,
			logo
		) VALUES (?, ?, ?, ?)`,
		team.Id,
		team.Acronym,
		team.Name,
		team.Logo,
	)

	return err
}

func (r *Repository) GetAll() ([]Match, error) {
	rows, err := r.db.Query(`SELECT
					id,
					tournament,
					team1_id,
					team2_id,
					score,
					best_of,
					time,
					location,
					url,
					is_live,
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
			&event.Team1Id,
			&event.Team2Id,
			&event.Score,
			&event.BestOf,
			&event.Time,
			&event.Location,
			&event.URL,
			&event.IsLive,
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

func (r *Repository) GetAllTeams() ([]Team, error) {
	rows, err := r.db.Query(`SELECT
			id,
			acronym,
			name,
			logo
		FROM teams`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	var teams []Team
	for rows.Next() {
		var team Team
		err := rows.Scan(
			&team.Id,
			&team.Acronym,
			&team.Name,
			&team.Logo,
		)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *Repository) GetByTeamId(teamID int, gameType string, after time.Time) ([]Match, error) {
	teams, err := r.GetAllTeams()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`SELECT
					id,
					tournament,
					team1_id,
					team2_id,
					score,
					best_of,
					time,
					location,
					url,
					is_live,
					game_type,
					modified_at
					FROM events WHERE (team1_id = ? OR team2_id = ?) AND game_type = ? AND time > ?`,
		teamID, teamID, gameType, after,
	)
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
			&event.Team1Id,
			&event.Team2Id,
			&event.Score,
			&event.BestOf,
			&event.Time,
			&event.Location,
			&event.URL,
			&event.IsLive,
			&event.GameType,
			&event.ModifiedAt,
		)
		if err != nil {
			log.Printf("[ERROR] error while scanning rows: %v", err)
			return nil, err
		}

		team1 := utils.FirstOrDefault[Team](teams, func(team *Team) bool {
			return team.Id == event.Team1Id
		})
		if team1 != nil {
			event.Team1 = *team1
		}

		team2 := utils.FirstOrDefault[Team](teams, func(team *Team) bool {
			return team.Id == event.Team2Id
		})
		if team2 != nil {
			event.Team2 = *team2
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] error while iterating rows: %v", err)
		return nil, err
	}

	return events, nil
}

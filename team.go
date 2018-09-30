package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// Team is a representation of a Team entity in the database
type Team struct {
	ID      int64  `json:"teamId,omitempty"`
	Name    string `json:"name,omitempty"`
	Captain int64  `json:"captain,omitempty"`
}

// TeamRequest is a representation of request to team routes
type TeamRequest struct {
	*Team
	ProtectedID int64  `json:"userId"`
	Action      string `json:"action,omitempty"`
	Invitee     int64  `json:"invitee,omitempty"`
}

// TeamResponse is representation of response from team routes
type TeamResponse struct {
	*Team
}

// Bind allows for preprocessing of requests
func (tr *TeamRequest) Bind(r *http.Request) error {
	if tr.Team == nil {
		return errors.New("missing team fields")
	}
	//tr.ProtectedId = -1
	return nil
}

// NewTeamResponse creates a TeamResponse
func NewTeamResponse(t *Team) *TeamResponse {
	resp := &TeamResponse{Team: t}
	if resp.Team == nil {
		// not sure resp.Team will be nil
	}
	return resp
}

// Render allows for preprocessing of TeamResponse
func (tr *TeamResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// make changes to team response before sent
	return nil
}

// TeamListResponse is a list of TeamResponses
type TeamListResponse []*TeamResponse

// NewTeamListResponse creates a list of team responses
func NewTeamListResponse(teams []*Team) []render.Renderer {
	list := []render.Renderer{}
	for _, team := range teams {
		list = append(list, NewTeamResponse(team))
	}
	return list
}

func dbNewTeam(team *Team) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO team(name,captain) VALUES(?,?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(team.Name, team.Captain)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func dbAddToRoster(userID, teamID int64) error {
	stmt, err := db.Prepare("INSERT INTO roster(teamID,userID) VALUES(?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(teamID, userID)
	if err != nil {
		return err
	}
	return nil
}

func dbRemoveFromRoster(userID, teamID int64) error {
	stmt, err := db.Prepare("DELETE FROM roster WHERE teamID=? AND userID=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(teamID, userID)
	if err != nil {
		return err
	}
	return nil
}

func dbSearchTeamName(searchValue, offset string) ([]*Team, error) {
	var teams []*Team
	rows, err := db.Query("SELECT id, name, captain FROM team WHERE name LIKE CONCAT(?, '%') LIMIT 10 OFFSET ?", searchValue, offset)
	if err != nil {
		return teams, err
	}
	defer rows.Close()
	for rows.Next() {
		var team Team
		err := rows.Scan(&team.ID, &team.Name, &team.Captain)
		if err != nil {
			return teams, err
		}
		teams = append(teams, &team)
	}
	err = rows.Err()
	if err != nil {
		return teams, err
	}
	return teams, nil
}

func dbGetUserTeams(userID string) ([]*Team, error) {
	var teams []*Team
	rows, err := db.Query("SELECT id, name, captain FROM team INNER JOIN roster WHERE team.ID=roster.teamID AND roster.userID = ?", userID)
	if err != nil {
		return teams, err
	}
	defer rows.Close()
	for rows.Next() {
		var team Team
		err := rows.Scan(&team.ID, &team.Name, &team.Captain)
		if err != nil {
			return teams, err
		}
		teams = append(teams, &team)
	}
	err = rows.Err()
	if err != nil {
		return teams, err
	}
	return teams, nil
}

func dbEditRoster(action string, userID, teamID int64) error {
	switch action {
	case "add":
		return dbAddToRoster(userID, teamID)
	case "remove":
		var value bool
		value, err := isCaptain(userID, teamID)
		if err != nil {
			return err
		}
		if value {
			return errors.New("captain cannot leave team")
		}
		return dbRemoveFromRoster(userID, teamID)
	default:
		return errors.New("action is not valid")
	}
}

func isCaptain(userID, teamID int64) (bool, error) {
	var captain string
	err := db.QueryRow("SELECT captain FROM team WHERE captain=? AND id=?", userID, teamID).Scan(&captain)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

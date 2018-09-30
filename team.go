package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type Team struct {
	Id      int64  `json:"teamId,omitempty"`
	Name    string `json:"name,omitempty"`
	Captain int64  `json:"captain,omitempty"`
}

type TeamRequest struct {
	*Team
	ProtectedId int64  `json:"userId"`
	Action      string `json:"action,omitempty"`
	Invitee     int64  `json:invitee,omitempty`
}

type TeamResponse struct {
	*Team
}

func (tr *TeamRequest) Bind(r *http.Request) error {
	if tr.Team == nil {
		return errors.New("missing team fields")
	}
	//tr.ProtectedId = -1
	return nil
}

func NewTeamResponse(t *Team) *TeamResponse {
	resp := &TeamResponse{Team: t}
	if resp.Team == nil {
		// not sure resp.Team will be nil
	}
	return resp
}

func (tr *TeamResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// make changes to team response before sent
	return nil
}

type TeamListResponse []*TeamResponse

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

func dbAddToRoster(userId, teamId int64) error {
	stmt, err := db.Prepare("INSERT INTO roster(teamId,userId) VALUES(?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(teamId, userId)
	if err != nil {
		return err
	}
	return nil
}

func dbRemoveFromRoster(userId, teamId int64) error {
	stmt, err := db.Prepare("DELETE FROM roster WHERE teamId=? AND userId=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(teamId, userId)
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
		err := rows.Scan(&team.Id, &team.Name, &team.Captain)
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

func dbGetUserTeams(userId string) ([]*Team, error) {
	var teams []*Team
	rows, err := db.Query("SELECT id, name, captain FROM team INNER JOIN roster WHERE team.id=roster.teamId AND roster.userId = ?", userId)
	if err != nil {
		return teams, err
	}
	defer rows.Close()
	for rows.Next() {
		var team Team
		err := rows.Scan(&team.Id, &team.Name, &team.Captain)
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

func dbEditRoster(action string, userId, teamId int64) error {
	switch action {
	case "add":
		return dbAddToRoster(userId, teamId)
	case "remove":
		var value bool
		value, err := isCaptain(userId, teamId)
		if err != nil {
			return err
		}
		if value {
			return errors.New("captain cannot leave team")
		} else {
			return dbRemoveFromRoster(userId, teamId)
		}
	default:
		return errors.New("action is not valid")
	}
}

func isCaptain(userId, teamId int64) (bool, error) {
	var captain string
	err := db.QueryRow("SELECT captain FROM team WHERE captain=? AND id=?", userId, teamId).Scan(&captain)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

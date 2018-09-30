package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// TeamInvite represents an invite to a team
// Name is used only for responses
type TeamInvite struct {
	Team    int64  `json:"teamId"`
	Invitee int64  `json:"invitee"`
	Name    string `json:"teamName,omitempty"`
}

// TeamInviteRequest represents a request to teaminvite routes
type TeamInviteRequest struct {
	*TeamInvite
	ProtectedID int64 `json:"userId"`
}

// TeamInviteResponse represents a response from teaminvite routes
type TeamInviteResponse struct {
	*TeamInvite
}

// Bind allows for preprocessing of requests
func (i *TeamInviteRequest) Bind(r *http.Request) error {
	if i.TeamInvite == nil {
		return errors.New("missing invite fields")
	}
	//i.ProtectedID = -1
	return nil
}

// NewTeamInviteResponse creates a response from and invite
func NewTeamInviteResponse(invite *TeamInvite) *TeamInviteResponse {
	resp := &TeamInviteResponse{TeamInvite: invite}
	if resp.TeamInvite == nil {
		// shouldn't be
	}
	return resp
}

// Render allows for preprocessing of responses
func (ir *TeamInviteResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// make changes to TeamInviteResponse
	return nil
}

// TeamInviteListResponse  is a list of teaminvite responses
type TeamInviteListResponse []*TeamInviteResponse

//NewTeamInviteListResponse creates a list of teaminvite reponses for invites
func NewTeamInviteListResponse(invites []*TeamInvite) []render.Renderer {
	list := []render.Renderer{}
	for _, invite := range invites {
		list = append(list, NewTeamInviteResponse(invite))
	}
	return list
}

func dbNewTeamInvite(invite *TeamInvite, userID int64) error {
	value, err := isCaptain(userID, invite.Team)
	if err != nil {
		return err
	}
	if !value {
		return errors.New("only captain can invite")
	}
	stmt, err := db.Prepare("INSERT INTO team_invite(teamId,invitee) VALUES(?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(invite.Team, invite.Invitee)
	if err != nil {
		return err
	}
	return nil
}

func dbGetUserTeamInvites(userID string) ([]*TeamInvite, error) {
	var invites []*TeamInvite
	rows, err := db.Query("SELECT team.name,team_invite.teamId,team_invite.invitee FROM team INNER JOIN team_invite WHERE team.id=team_invite.teamId AND team_invite.invitee = ?", userID)
	if err != nil {
		return invites, err
	}
	defer rows.Close()
	for rows.Next() {
		var invite TeamInvite
		err := rows.Scan(&invite.Name, &invite.Team, &invite.Invitee)
		if err != nil {
			return invites, err
		}
		invites = append(invites, &invite)
	}
	err = rows.Err()
	if err != nil {
		return invites, err
	}
	return invites, nil
}

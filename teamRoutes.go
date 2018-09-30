package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
)

// TeamRoutes returns a router with the team routes to be mounted in routes.go
func TeamRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", CreateTeam)
	r.Get("/search/{value}/{offset}", SearchTeam)
	r.Get("/by-user/{userID}", GetUserTeams)
	r.Put("/", ModifyRoster)
	return r
}

// ModifyRoster allows for adding and removing users to and from a team
func ModifyRoster(w http.ResponseWriter, r *http.Request) {
	data := &TeamRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	// TODO check to make sure team is not in active tournaments
	if data.Action == "add" && data.Invitee != data.ProtectedID {
		render.Render(w, r, ErrUnauthorized(errors.New("not intended user")))
		return
	}
	err := dbEditRoster(data.Action, data.ProtectedID, data.Team.ID)
	if err != nil {
		_, ok := err.(*mysql.MySQLError)
		if !ok {
			render.Render(w, r, ErrForbidden(err))
			return
		}
		render.Render(w, r, ErrDB(err))
		return
	}
	render.Render(w, r, NewTeamResponse(data.Team))
}

// GetUserTeams renders a list of all teams of which the userID is a member
func GetUserTeams(w http.ResponseWriter, r *http.Request) {
	var teamList []*Team
	var err error
	if userID := chi.URLParam(r, "userID"); userID != "" {
		teamList, err = dbGetUserTeams(userID)
	} else {
		render.Render(w, r, ErrBadRequest(errors.New("user id not valid")))
		return
	}
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	if err := render.RenderList(w, r, NewTeamListResponse(teamList)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// SearchTeam searches for teams with name values starting with the search value
func SearchTeam(w http.ResponseWriter, r *http.Request) {
	var err error
	var teamList []*Team
	offset := chi.URLParam(r, "offset")
	if searchValue := chi.URLParam(r, "value"); searchValue != "" {
		teamList, err = dbSearchTeamName(searchValue, offset)
	} else {
		render.Render(w, r, ErrBadRequest(errors.New("search value empty")))
		return
	}
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	if err := render.RenderList(w, r, NewTeamListResponse(teamList)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreateTeam creates a team in the database
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	data := &TeamRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	team := data.Team
	team.Captain = data.ProtectedID
	id, err := dbNewTeam(team)
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	err = dbAddToRoster(data.ProtectedID, id)
	if err != nil {
		// need to retry adding to roster or delete created team and give error
	}
	team.ID = id
	render.Render(w, r, NewTeamResponse(team))
}

package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
)

func TeamRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", CreateTeam)
	r.Get("/search/{value}/{offset}", SearchTeam)
	r.Get("/by-user/{userId}", GetUserTeams)
	r.Put("/", ModifyRoster)
	return r
}

func ModifyRoster(w http.ResponseWriter, r *http.Request) {
	data := &TeamRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	// TODO check to make sure team is not in active tournaments
	if data.Action == "add" && data.Invitee != data.ProtectedId {
		render.Render(w, r, ErrUnauthorized(errors.New("not intended user")))
		return
	}
	err := dbEditRoster(data.Action, data.ProtectedId, data.Team.Id)
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

func GetUserTeams(w http.ResponseWriter, r *http.Request) {
	var teamList []*Team
	var err error
	if userId := chi.URLParam(r, "userId"); userId != "" {
		teamList, err = dbGetUserTeams(userId)
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

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	data := &TeamRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	team := data.Team
	team.Captain = data.ProtectedId
	id, err := dbNewTeam(team)
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	err = dbAddToRoster(data.ProtectedId, id)
	if err != nil {
		// need to retry adding to roster or delete created team and give error
	}
	team.Id = id
	render.Render(w, r, NewTeamResponse(team))
}

package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
)

// TeamInviteRoutes returns a router with team invite routes to be mounted in routes.go
func TeamInviteRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", CreateTeamInvite)
	r.Get("by-user/{userID}", GetUserTeamInvites)
	return r
}

// CreateTeamInvite creates a team invite in the database
func CreateTeamInvite(w http.ResponseWriter, r *http.Request) {
	data := &TeamInviteRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	invite := data.TeamInvite
	err := dbNewTeamInvite(invite, data.ProtectedID)
	if err != nil {
		_, ok := err.(*mysql.MySQLError)
		if !ok {
			render.Render(w, r, ErrForbidden(err))
			return
		}
		render.Render(w, r, ErrDB(err))
		return
	}
	render.Render(w, r, NewTeamInviteResponse(invite))
}

// GetUserTeamInvites renders all team invites for a given userID
func GetUserTeamInvites(w http.ResponseWriter, r *http.Request) {
	var inviteList []*TeamInvite
	var err error
	if userID := chi.URLParam(r, "userID"); userID != "" {
		inviteList, err = dbGetUserTeamInvites(userID)
	} else {
		render.Render(w, r, ErrBadRequest(errors.New("user id not valid")))
		return
	}
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	if err := render.RenderList(w, r, NewTeamInviteListResponse(inviteList)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

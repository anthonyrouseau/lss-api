package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
)

func TeamInviteRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", CreateTeamInvite)
	r.Get("by-user/{userId}", GetUserTeamInvites)
	return r
}

func CreateTeamInvite(w http.ResponseWriter, r *http.Request) {
	data := &TeamInviteRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	invite := data.TeamInvite
	err := dbNewTeamInvite(invite, data.ProtectedId)
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

func GetUserTeamInvites(w http.ResponseWriter, r *http.Request) {
	var inviteList []*TeamInvite
	var err error
	if userId := chi.URLParam(r, "userId"); userId != "" {
		inviteList, err = dbGetUserTeamInvites(userId)
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

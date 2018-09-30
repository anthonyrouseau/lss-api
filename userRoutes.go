package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// UserRoutes returns router with user routes to be mounted in routes.go
func UserRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", CreateUser)
	r.Get("/search/{value}/{offset}", SearchUser)
	return r
}

// SearchUser searches for a user with username starting with given value
func SearchUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var userList []*User
	offset := chi.URLParam(r, "offset")
	if searchValue := chi.URLParam(r, "value"); searchValue != "" {
		userList, err = dbSearchUsername(searchValue, offset)
	} else {
		render.Render(w, r, ErrBadRequest(errors.New("search value empty")))
		return
	}
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	if err := render.RenderList(w, r, NewUserListResponse(userList)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreateUser creates a user in the databse
func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	user := data.User
	summonerID, err := checkSummonerID(data.SummonerName, data.Code)
	if err != nil {
		render.Render(w, r, ErrAccountLink(err))
		return
	}
	user.SummonerID = summonerID
	id, err := dbNewUser(user)
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	user.ID = id
	render.Render(w, r, NewUserResponse(user))
}

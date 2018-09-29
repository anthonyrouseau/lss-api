package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func UserRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", CreateUser)
	r.Get("/search/{value}/{offset}", SearchUser)
	return r
}

func SearchUser(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	user := data.User
	id, err := dbNewUser(user)
	if err != nil {
		render.Render(w, r, ErrDB(err))
		return
	}
	user.Id = id
	render.Render(w, r, NewUserResponse(user))
}

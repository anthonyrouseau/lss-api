package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.DefaultCompress,
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("api running"))
	})
	r.Mount("/user", UserRoutes())
	return r
}

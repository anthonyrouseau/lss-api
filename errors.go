package main

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
)

type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:status`
	AppCode        uint16 `json:code,omitempty`
	ErrorText      string `json:error,omitempty`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Bad Request",
		ErrorText:      err.Error(),
	}
}

func ErrDB(err error) render.Renderer {
	e := err.(*mysql.MySQLError)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Database Error",
		AppCode:        e.Number,
		ErrorText:      e.Message,
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

func ErrAccountLink(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Could Not Link Account",
		ErrorText:      err.Error(),
	}
}

package main

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-sql-driver/mysql"
)

// ErrResponse represents an error to be sent to the client
type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	AppCode        uint16 `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}

// Render provides preprocessing before error response sent to client
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrBadRequest is a Bad Request error
func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Bad Request",
		ErrorText:      err.Error(),
	}
}

// ErrDB is and error resulting from a database query or process
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

// ErrRender is an error that occurs when a render method is called
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

// ErrAccountLink is an error that occurs when trying to link summoner to account
func ErrAccountLink(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Could Not Link Account",
		ErrorText:      err.Error(),
	}
}

// ErrUnauthorized is an error that occurs when a user is not authorized
func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 401,
		StatusText:     "Unauthorized",
		ErrorText:      err.Error(),
	}
}

// ErrForbidden is an error that occurs when something is not allowed by the app
func ErrForbidden(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 403,
		StatusText:     "Forbidden",
		ErrorText:      err.Error(),
	}
}

package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// ////
// Generic Responses
// ////
type SimpleIdResponse struct {
	Id int `json:"id"`
}

// ////
// Error Responses
// ////
type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 400,
		StatusText:     "invalid_request",
	}
}

func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 500,
		StatusText:     "internal_server_error",
	}
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 404,
		StatusText:     "not_found",
	}
}

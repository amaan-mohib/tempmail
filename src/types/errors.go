package types

import (
	"fmt"
	"net/http"
	"tempgalias/src/utils"

	"github.com/go-chi/render"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"`      // low-level runtime error
	HTTPStatusCode int   `json:"status"` // http response status code

	StatusText string `json:"message"`        // user-level status message
	AppCode    int64  `json:"code,omitempty"` // application-specific error code
	ErrorText  string `json:"-"`              // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	fmt.Println(err)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
func ErrInvalidRequestWithMessage(err error, message string) render.Renderer {
	fmt.Println(err)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     utils.InlineIf(message != "", message, "Invalid request."),
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	fmt.Println(err)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServer(err error) render.Renderer {
	fmt.Println(err)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

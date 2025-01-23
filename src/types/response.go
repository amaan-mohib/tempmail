package types

import (
	"net/http"

	"github.com/go-chi/render"
)

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Render implements render.Renderer.
func (b *BaseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, b.Status)
	return nil
}

func Response(response BaseResponse) render.Renderer {
	return &response
}

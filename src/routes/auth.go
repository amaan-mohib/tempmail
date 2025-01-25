package routes

import (
	"net/http"
	"tempgalias/src/config"
	"tempgalias/src/types"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Bind LoginRequest
func (l *LoginRequest) Bind(r *http.Request) error {
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := &LoginRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	//check user in db

	tokenAuth := jwtauth.New("HS256", []byte(config.Config.Authentication.JWTSecret), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"email": data.Email})
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err))
	}

	response := &types.BaseResponse{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    map[string]string{"token": tokenString},
	}

	render.Render(w, r, types.Response(*response))
}

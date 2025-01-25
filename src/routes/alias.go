package routes

import (
	"fmt"
	"net/http"
	"tempgalias/src/config"
	"tempgalias/src/database"
	"tempgalias/src/types"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/stephenafamo/scan"
)

type JWTAlias struct {
	Alias     string `json:"alias"`
	CreatedAt int64  `json:"createdAt"`
	ExpiryAt  int64  `json:"expiryAt"`
}

type GetAliasRequest struct {
	Alias string `json:"alias,omitempty"`
}

// Bind GetAliasRequest
func (l *GetAliasRequest) Bind(r *http.Request) error {
	return nil
}

func getAliasObject(w http.ResponseWriter, r *http.Request) *types.BaseResponse {
	data, err := GenerateAlias()
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err))
	}
	response := &types.BaseResponse{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    data,
	}

	return response
}

func GetAliasHandler(w http.ResponseWriter, r *http.Request) {
	data := &GetAliasRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	x, _ := database.ScanOne(scan.SingleColumnMapper[int], "select id from users")
	fmt.Println(x)
	if data.Alias == "" {
		response := getAliasObject(w, r)

		render.Render(w, r, types.Response(*response))
		return
	}

	result, err := DecodeAliasToken(data.Alias)
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err))
		return
	}

	//check alias in db, if not then throw error else check expiry and generate new

	expiryAt := result.ExpiryAt
	if expiryAt <= time.Now().UnixMilli() {
		response := getAliasObject(w, r)
		render.Render(w, r, types.Response(*response))
		return
	}

	response := &types.BaseResponse{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    map[string]string{},
	}

	render.Render(w, r, types.Response(*response))
}

func GenerateAlias() (map[string]interface{}, error) {
	// tokenAuth := jwtauth.New("HS256", []byte(Config.Authentication.JWTSecret), nil)
	baseEmail := config.Config.BaseEmail
	currentTime := time.Now()
	createdAt := currentTime.UnixMilli()
	expiryAt := currentTime.Add(time.Millisecond * time.Duration(config.Config.DefaultExpiry)).UnixMilli()
	uuid, _ := gonanoid.New(5)
	emailAlias := fmt.Sprintf("%s+%s@gmail.com", baseEmail, uuid)

	value := map[string]interface{}{
		"alias":     emailAlias,
		"createdAt": createdAt,
		"expiryAt":  expiryAt,
	}
	return value, nil
	// _, tokenString, err := tokenAuth.Encode(value)

	// return tokenString, err
}

func DecodeAliasToken(token string) (JWTAlias, error) {
	tokenAuth := jwtauth.New("HS256", []byte(config.Config.Authentication.JWTSecret), nil)
	jwtToken, err := jwtauth.VerifyToken(tokenAuth, token)
	if err != nil {
		return JWTAlias{}, err
	}
	alias, _ := jwtToken.Get("alias")
	expiryAt, _ := jwtToken.Get("expiryAt")
	createdAt, _ := jwtToken.Get("createdAt")

	return JWTAlias{
		Alias:     alias.(string),
		CreatedAt: int64(createdAt.(float64)),
		ExpiryAt:  int64(expiryAt.(float64)),
	}, nil
}

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

type JSONAlias struct {
	Alias     string `json:"alias"`
	CreatedAt int64  `json:"createdAt"`
	ExpiryAt  int64  `json:"expiryAt"`
}
type DBAlias struct {
	Alias     string `db:"alias"`
	CreatedAt int64  `db:"created_at_ts"`
	ExpiryAt  int64  `db:"expiry_at"`
}

type GetAliasRequest struct {
	Alias string `json:"alias,omitempty"`
}

// Bind GetAliasRequest
func (l *GetAliasRequest) Bind(r *http.Request) error {
	return nil
}

func getAliasObject() (*types.BaseResponse, error) {
	data, err := GenerateAlias()
	if err != nil {
		return &types.BaseResponse{}, err
	}
	response := &types.BaseResponse{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    data,
	}

	return response, nil
}

func GetAliasHandler(w http.ResponseWriter, r *http.Request) {
	data := &GetAliasRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	if data.Alias == "" {
		response, err := getAliasObject()
		if err != nil {
			render.Render(w, r, types.ErrInternalServer(err))
			return
		}

		render.Render(w, r, types.Response(*response))
		return
	}

	//check alias in db, if not then throw error else check expiry and generate new
	alias := data.Alias

	result, err := database.ScanOne(
		scan.StructMapper[DBAlias](),
		"select alias, expiry_at, created_at_ts from aliases where alias = $1",
		alias,
	)

	if err != nil {
		render.Render(w, r, types.ErrInvalidRequestWithMessage(err, "No such alias exist"))
		return
	}

	expiryAt := result.ExpiryAt
	if expiryAt <= time.Now().UnixMilli() {
		response, err := getAliasObject()
		if err != nil {
			render.Render(w, r, types.ErrInternalServer(err))
			return
		}

		render.Render(w, r, types.Response(*response))
		return
	}

	response := &types.BaseResponse{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    JSONAlias(result),
	}

	render.Render(w, r, types.Response(*response))
}

func GenerateAlias() (map[string]interface{}, error) {
	baseEmail := config.Config.BaseEmail
	currentTime := time.Now()
	createdAt := currentTime.UnixMilli()
	expiryAt := currentTime.Add(time.Millisecond * time.Duration(config.Config.DefaultExpiry)).UnixMilli()
	uuid, err := gonanoid.New(5)
	if err != nil {
		return map[string]interface{}{}, err
	}
	emailAlias := fmt.Sprintf("%s+%s@gmail.com", baseEmail, uuid)

	_, err = database.Query(
		`insert into aliases (alias, created_at_ts, expiry_at) values ($1, $2, $3)`,
		emailAlias,
		createdAt,
		expiryAt,
	)
	if err != nil {
		return map[string]interface{}{}, err
	}

	value := map[string]interface{}{
		"alias":     emailAlias,
		"createdAt": createdAt,
		"expiryAt":  expiryAt,
	}
	return value, nil
}

func DecodeAliasToken(token string) (JSONAlias, error) {
	tokenAuth := jwtauth.New("HS256", []byte(config.Config.Authentication.JWTSecret), nil)
	jwtToken, err := jwtauth.VerifyToken(tokenAuth, token)
	if err != nil {
		return JSONAlias{}, err
	}
	alias, _ := jwtToken.Get("alias")
	expiryAt, _ := jwtToken.Get("expiryAt")
	createdAt, _ := jwtToken.Get("createdAt")

	return JSONAlias{
		Alias:     alias.(string),
		CreatedAt: int64(createdAt.(float64)),
		ExpiryAt:  int64(expiryAt.(float64)),
	}, nil
}

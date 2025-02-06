package error

import (
	oerror "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/http"
)

var CustomCodeToHttpCodeMapping = map[oerror.Code]http.StatusCode{
	RequestInvalid:               http.StatusBadRequest,
	NotFound:                     http.StatusBadRequest,
	RequestNotValid:              http.StatusForbidden,
	SqlCreateError:               http.StatusInternalServerError,
	CreateJwtTokenError:          http.StatusUnauthorized,
	CreateOauthRefreshTokenError: http.StatusUnauthorized,
	CreateAccessToken:            http.StatusUnauthorized,
	LogIn:                        http.StatusUnauthorized,
	UpdateOauthRefreshTokenError: http.StatusUnauthorized,
	BadRequest:                   http.StatusBadRequest,
	AccessTokenExpire:            http.StatusUnauthorized,
	RefreshTokenExpire:           http.StatusUnauthorized,
}

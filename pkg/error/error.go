package error

import (
	"context"
	"github.com/gin-gonic/gin"
	oerror "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/response"
)

const (
	BadRequest                   oerror.Code = "BAD_REQUEST"
	NotFound                     oerror.Code = "NOT_FOUND"
	RequestNotValid              oerror.Code = "REQUEST_NOT_VALID"
	RequestInvalid               oerror.Code = "REQUEST_INVALID"
	RedisError                   oerror.Code = "REDIS_ERROR"
	UnmarshalError               oerror.Code = "UNMARSHAl_ERROR"
	MarshalError                 oerror.Code = "MARSHAL_ERR"
	ParseIntError                oerror.Code = "PARSE_INT_ERROR"
	SqlCreateError               oerror.Code = "SQL_CREATE_ERROR"
	SqlUpdateError               oerror.Code = "SQL_UPDATE_ERROR"
	SqlFetchError                oerror.Code = "SQL_FETCH_ERROR"
	SqlDeleteError               oerror.Code = "SQL_DELETE_ERROR"
	NoRowsAffectedError          oerror.Code = "NO_ROWS_AFFECTED_ERROR"
	SomethingWentWrong           oerror.Code = "SOMETHING_WENT_WRONG"
	CacheGetError                oerror.Code = "CACHE_GET_ERROR"
	CacheSetError                oerror.Code = "CACHE_SET_ERROR"
	DataNotFoundDbError          oerror.Code = "DATA_NOT_FOUND_DB_ERROR"
	SqlUpsertError               oerror.Code = "SQL_UPSERT_ERROR"
	GoroutineError               oerror.Code = "GOROUTINE_ERROR"
	CachePurgeError              oerror.Code = "CACHE_PURGE_ERROR"
	CreateJwtTokenError          oerror.Code = "CREATE_JWT_TOKEN_ERROR"
	CreateOauthRefreshTokenError oerror.Code = "CREATE_OAUTH_REFRESH_JWT_TOKEN_ERROR"
	CreateAccessToken            oerror.Code = "CREATE_ACCESS_TOKEN"
	LogIn                        oerror.Code = "LOG_IN_ERROR"
	UpdateOauthRefreshTokenError oerror.Code = "UPDATE_AUTH_REFRESH_TOKEN_ERROR"
	ParseFilesError              oerror.Code = "PARSE_FILES_ERROR"
	SendEMailError               oerror.Code = "SEND_EMAIL_ERROR"
	NotFoundMapError             oerror.Code = "NOT_FOUND_MAP_ERROR"
	UrlError                     oerror.Code = "URL_ERROR"
	BicryptError                 oerror.Code = "BICRYPT_ERROR"
	SqsPublishErr                oerror.Code = "SQS_PUBLISH_MESSAGE"
	SqsInitializeErr             oerror.Code = "SQS_INITIALIZE_MESSAGE"
	InternalServerError          oerror.Code = "INTERNAL_SERVER_ERROR"
	AccessTokenExpire            oerror.Code = "ACCESS_TOKEN_ERROR"
	RefreshTokenExpire           oerror.Code = "REFRESH_TOKEN_EXPIRE"
	S3CopyObjectError            oerror.Code = "S3_COPY_OBJECT_ERROR"
)

func NewErrorResponse(ctx *gin.Context, cusErr oerror.CustomError) {
	response.NewErrorResponse(ctx, cusErr, CustomCodeToHttpCodeMapping)
}

func NewErrorResponseWithData(ctx *gin.Context, cusErr oerror.CustomErrorWithData) {
	response.NewErrorResponseWithData(ctx, cusErr, CustomCodeToHttpCodeMapping)
}

func InvalidRequest(ctx context.Context, key string) oerror.CustomError {
	message := i18n.Translate(ctx, key)
	return oerror.NewCustomError(oerror.RequestInvalid, message)
}

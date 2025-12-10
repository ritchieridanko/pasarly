package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Internal error codes
const (
	CodeAuthNotFound       errCode = "AUTH_NOT_FOUND_ERR"
	CodeCacheQueryExec     errCode = "CACHE_QUERY_EXEC_ERR"
	CodeCacheScriptExec    errCode = "CACHE_SCRIPT_EXEC_ERR"
	CodeCookieNotFound     errCode = "COOKIE_NOT_FOUND_ERR"
	CodeCtxValueNotFound   errCode = "CTX_VALUE_NOT_FOUND_ERR"
	CodeDataConflict       errCode = "DATA_CONFLICT_ERR"
	CodeDBQueryExec        errCode = "DB_QUERY_EXEC_ERR"
	CodeDBTx               errCode = "DB_TX_ERR"
	CodeJWTCreationFailed  errCode = "JWT_CREATION_FAILED_ERR"
	CodeHashingFailed      errCode = "HASHING_FAILED_ERR"
	CodeInternal           errCode = "INTERNAL_ERR"
	CodeInvalidCredentials errCode = "INVALID_CREDENTIALS_ERR"
	CodeInvalidParams      errCode = "INVALID_PARAMS_ERR"
	CodeInvalidPayload     errCode = "INVALID_PAYLOAD_ERR"
	CodeInvalidToken       errCode = "INVALID_TOKEN_ERR"
	CodeSessionNotFound    errCode = "SESSION_NOT_FOUND_ERR"
	CodeTokenExpired       errCode = "TOKEN_EXPIRED_ERR"
	CodeTokenMalformed     errCode = "TOKEN_MALFORMED_ERR"
	CodeUnauthenticated    errCode = "UNAUTHENTICATED_ERR"
	CodeUnknown            errCode = "UNKNOWN_ERR"
	CodeWrongSignInMethod  errCode = "WRONG_SIGN_IN_METHOD_ERR"
)

// External error messages
const (
	MsgEmailAlreadyRegistered string = "Email is already registered"
	MsgInternalServer         string = "Internal server error"
	MsgInvalidCredentials     string = "Invalid credentials"
	MsgInvalidParams          string = "Invalid params"
	MsgInvalidPayload         string = "Invalid payload"
	MsgUnauthenticated        string = "Unauthenticated"
)

// Internal errors
var (
	ErrCacheNil               error = redis.Nil
	ErrDBAffectNoRows         error = errors.New("no rows affected")
	ErrDBReturnNoRows         error = pgx.ErrNoRows
	ErrEmailAlreadyRegistered error = errors.New("email already registered")
	ErrEmailReserved          error = errors.New("email reserved")
	ErrEventOnProcess         error = errors.New("message is being processed on another instance")
	ErrInvalidToken           error = errors.New("invalid token")
	ErrWrongSignInMethod      error = errors.New("wrong sign in method")
)

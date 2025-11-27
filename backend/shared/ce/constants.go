package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal error codes
const (
	CodeAuthNotFound       errCode = "AUTH_NOT_FOUND_ERR"
	CodeCacheQueryExec     errCode = "CACHE_QUERY_EXEC_ERR"
	CodeCacheScriptExec    errCode = "CACHE_SCRIPT_EXEC_ERR"
	CodeCookieNotFound     errCode = "COOKIE_NOT_FOUND_ERR"
	CodeDataConflict       errCode = "DATA_CONFLICT_ERR"
	CodeDBQueryExec        errCode = "DB_QUERY_EXEC_ERR"
	CodeDBTx               errCode = "DB_TX_ERR"
	CodeJWTCreationFailed  errCode = "JWT_CREATION_FAILED_ERR"
	CodeHashingFailed      errCode = "HASHING_FAILED_ERR"
	CodeInternal           errCode = "INTERNAL_ERR"
	CodeInvalidCredentials errCode = "INVALID_CREDENTIALS_ERR"
	CodeInvalidPayload     errCode = "INVALID_PAYLOAD_ERR"
	CodeSessionNotFound    errCode = "SESSION_NOT_FOUND_ERR"
	CodeUnauthenticated    errCode = "UNAUTHENTICATED_ERR"
	CodeUnknown            errCode = "UNKNOWN_ERR"
	CodeWrongSignInMethod  errCode = "WRONG_SIGN_IN_METHOD_ERR"
)

// External error messages
const (
	MsgEmailAlreadyRegistered string = "Email is already registered"
	MsgInternalServer         string = "Internal server error"
	MsgInvalidCredentials     string = "Invalid credentials"
	MsgInvalidPayload         string = "Invalid payload"
	MsgUnauthenticated        string = "Unauthenticated"
)

// Internal errors
var (
	ErrDBAffectNoRows         error = errors.New("no rows affected")
	ErrDBReturnNoRows         error = pgx.ErrNoRows
	ErrEmailAlreadyRegistered error = errors.New("email already registered")
	ErrEmailReserved          error = errors.New("email reserved")
	ErrWrongSignInMethod      error = errors.New("wrong sign in method")
)

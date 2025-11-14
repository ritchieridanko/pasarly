package ce

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Internal error codes
const (
	CodeCacheQueryExec    errCode = "CACHE_QUERY_EXEC_ERR"
	CodeCacheScriptExec   errCode = "CACHE_SCRIPT_EXEC_ERR"
	CodeDataConflict      errCode = "DATA_CONFLICT_ERR"
	CodeDBQueryExec       errCode = "DB_QUERY_EXEC_ERR"
	CodeDBTx              errCode = "DB_TX_ERR"
	CodeJWTCreationFailed errCode = "JWT_CREATION_FAILED_ERR"
	CodeHashingFailed     errCode = "HASHING_FAILED_ERR"
	CodeInvalidPayload    errCode = "INVALID_PAYLOAD_ERR"
)

// External error messages
const (
	MsgEmailAlreadyRegistered string = "Email is already registered"
	MsgInternalServer         string = "Internal server error"
	MsgInvalidPayload         string = "Invalid payload"
)

// Internal errors
var (
	ErrDBAffectNoRows         error = errors.New("no rows affected")
	ErrDBReturnNoRows         error = pgx.ErrNoRows
	ErrEmailAlreadyRegistered error = errors.New("email already registered")
	ErrEmailReserved          error = errors.New("email reserved")
)

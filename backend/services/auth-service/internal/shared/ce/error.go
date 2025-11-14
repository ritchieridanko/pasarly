package ce

import (
	"net/http"

	otc "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errCode string

type Error struct {
	Code    errCode
	Message string
	Err     error
}

func NewError(s trace.Span, c errCode, msg string, err error) *Error {
	if s != nil {
		s.RecordError(err)
		s.SetStatus(otc.Error, msg)
	}

	return &Error{Code: c, Message: msg, Err: err}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) GRPCStatus() error {
	switch e.Code {
	case CodeInvalidPayload:
		return status.Error(codes.InvalidArgument, e.Message)
	case CodeAuthNotFound, CodeInvalidCredentials, CodeWrongSignInMethod:
		return status.Error(codes.Unauthenticated, e.Message)
	case CodeDataConflict:
		return status.Error(codes.AlreadyExists, e.Message)
	case
		CodeCacheQueryExec, CodeCacheScriptExec, CodeDBQueryExec,
		CodeDBTx, CodeHashingFailed, CodeJWTCreationFailed:
		return status.Error(codes.Internal, e.Message)
	default:
		return status.Error(codes.Internal, e.Message)
	}
}

func (e *Error) HTTPStatus() int {
	switch e.Code {
	case CodeInvalidPayload:
		return http.StatusBadRequest
	case CodeAuthNotFound, CodeInvalidCredentials, CodeWrongSignInMethod:
		return http.StatusUnauthorized
	case CodeDataConflict:
		return http.StatusConflict
	case
		CodeCacheQueryExec, CodeCacheScriptExec, CodeDBQueryExec,
		CodeDBTx, CodeHashingFailed, CodeJWTCreationFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

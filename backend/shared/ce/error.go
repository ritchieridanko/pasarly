package ce

import (
	"net/http"

	otc "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	gc "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errCode string

type Error struct {
	Code    errCode
	Message string
	Err     error
}

func NewError(s trace.Span, ec errCode, msg string, err error) *Error {
	if s != nil {
		s.RecordError(err)
		s.SetStatus(otc.Error, msg)
	}

	return &Error{Code: ec, Message: msg, Err: err}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) ToGRPCStatus() error {
	switch e.Code {
	case CodeInvalidPayload:
		return status.Error(gc.InvalidArgument, e.Message)
	case CodeAuthNotFound, CodeInvalidCredentials, CodeSessionNotFound, CodeWrongSignInMethod:
		return status.Error(gc.Unauthenticated, e.Message)
	case CodeDataConflict:
		return status.Error(gc.AlreadyExists, e.Message)
	case
		CodeCacheQueryExec, CodeCacheScriptExec, CodeDBQueryExec,
		CodeDBTx, CodeHashingFailed, CodeJWTCreationFailed:
		return status.Error(gc.Internal, e.Message)
	default:
		return status.Error(gc.Internal, e.Message)
	}
}

func (e *Error) ToHTTPStatus() int {
	switch e.Code {
	case CodeInvalidPayload:
		return http.StatusBadRequest
	case CodeUnauthenticated:
		return http.StatusUnauthorized
	case CodeDataConflict:
		return http.StatusConflict
	case CodeInternal, CodeUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

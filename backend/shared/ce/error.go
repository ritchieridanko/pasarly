package ce

import (
	"fmt"
	"net/http"
	"time"

	otc "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	gc "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errCode string

type Error struct {
	Code      errCode
	Message   string
	Err       error
	Timestamp time.Time
}

func NewError(s trace.Span, ec errCode, msg string, err error) *Error {
	if s != nil {
		s.RecordError(err)
		s.SetStatus(otc.Error, msg)
	}

	return &Error{Code: ec, Message: msg, Err: err, Timestamp: time.Now().UTC()}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s\t[%s]\t%v", e.Timestamp.Format("2006-01-02 15:04:05"), e.Code, e.Err)
}

func (e *Error) ToGRPCStatus() error {
	switch e.Code {
	case CodeInvalidPayload:
		return status.Error(gc.InvalidArgument, e.Message)
	case CodeAuthNotFound, CodeInvalidCredentials, CodeSessionNotFound, CodeWrongSignInMethod:
		return status.Error(gc.Unauthenticated, e.Message)
	case CodeUserNotFound:
		return status.Error(gc.NotFound, e.Message)
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
	case CodeInvalidParams, CodeInvalidPayload:
		return http.StatusBadRequest
	case
		CodeCookieNotFound,
		CodeInvalidToken,
		CodeTokenExpired,
		CodeTokenMalformed,
		CodeUnauthenticated,
		CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeNotFound, CodeUserNotFound:
		return http.StatusNotFound
	case CodeDataConflict:
		return http.StatusConflict
	case CodeCtxValueNotFound, CodeInternal, CodeUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

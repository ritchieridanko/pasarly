package ce

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromGRPCErr(s trace.Span, err error) *Error {
	st, ok := status.FromError(err)
	e := fmt.Errorf("%s", st.Message())
	if !ok {
		return NewError(s, CodeUnknown, MsgInternalServer, e)
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return NewError(s, CodeDataConflict, st.Message(), e)
	case codes.InvalidArgument:
		return NewError(s, CodeInvalidPayload, st.Message(), e)
	case codes.NotFound:
		return NewError(s, CodeNotFound, st.Message(), e)
	case codes.Unauthenticated:
		return NewError(s, CodeUnauthenticated, st.Message(), e)
	case codes.Internal:
		return NewError(s, CodeInternal, st.Message(), e)
	default:
		return NewError(s, CodeInternal, st.Message(), e)
	}
}

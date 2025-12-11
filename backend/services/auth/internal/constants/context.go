package constants

type ctxKey string

const (
	CtxKeyRequestID ctxKey = "x-request-id"
)

const (
	CtxKeyIPAddress string = "x-ip-address"
	CtxKeyUserAgent string = "x-user-agent"
)

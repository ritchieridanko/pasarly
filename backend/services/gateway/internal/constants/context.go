package constants

type ctxKey string

const (
	CtxKeyAuthID     ctxKey = "x-auth-id"
	CtxKeyIsVerified ctxKey = "x-is-verified"
	CtxKeyRequestID  ctxKey = "x-request-id"
	CtxKeyRole       ctxKey = "x-role"
)

const (
	CtxKeyIPAddress string = "x-ip-address"
	CtxKeyUserAgent string = "x-user-agent"
)

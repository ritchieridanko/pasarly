package constants

type ctxKey string

const (
	CtxKeyAuthID     ctxKey = "auth-id"
	CtxKeyIsVerified ctxKey = "is-verified"
	CtxKeyRequestID  ctxKey = "request-id"
	CtxKeyRole       ctxKey = "role"
)

const (
	CtxKeyIPAddress string = "ip-address"
	CtxKeyUserAgent string = "user-agent"
)

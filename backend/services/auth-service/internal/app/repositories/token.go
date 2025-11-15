package repositories

import (
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/auth-service/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/infra/cache"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/shared/ce"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/internal/shared/constants"
	"go.opentelemetry.io/otel"
	"golang.org/x/net/context"
)

const tokenErrTracer string = "repository.token"

type TokenRepository interface {
	CreateVerificationToken(ctx context.Context, authID int64, token string) (err *ce.Error)
}

type tokenRepository struct {
	config *configs.Auth
	cache  *cache.Cache
}

func NewTokenRepository(cfg *configs.Auth, c *cache.Cache) TokenRepository {
	return &tokenRepository{config: cfg, cache: c}
}

func (r *tokenRepository) CreateVerificationToken(ctx context.Context, authID int64, token string) *ce.Error {
	ctx, span := otel.Tracer(tokenErrTracer).Start(ctx, "CreateVerificationToken")
	defer span.End()

	prefix := constants.CachePrefixVerification
	authKey := fmt.Sprintf("%s:%d", prefix, authID)
	tokenKey := fmt.Sprintf("%s:%s", prefix, token)
	duration := int(r.config.Token.Duration.Verification.Seconds())

	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
		redis.call("SET", KEYS[2], ARGV[2], "EX", ARGV[3])
		return 1
	`

	_, err := r.cache.Evaluate(
		ctx, "hs:cvt", script,
		[]string{authKey, tokenKey, prefix}, token, authID, duration,
	)
	if err != nil {
		e := fmt.Errorf("failed to create verification token: %w", err)
		return ce.NewError(span, ce.CodeCacheScriptExec, ce.MsgInternalServer, e)
	}

	return nil
}

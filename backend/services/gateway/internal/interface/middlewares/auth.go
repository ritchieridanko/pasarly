package middlewares

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const authErrTracer string = "middleware.auth"

func Authenticate(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, span := otel.Tracer(authErrTracer).Start(ctx.Request.Context(), "Authenticate")
		defer span.End()

		authorization := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if len(authorization) == 0 {
			e := fmt.Errorf("failed to authenticate: %w", errors.New("authorization is not provided"))
			ctx.Error(ce.NewError(span, ce.CodeUnauthenticated, ce.MsgUnauthenticated, e))
			ctx.Abort()
			return
		}

		auth := strings.Split(authorization, " ")
		if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
			e := fmt.Errorf("failed to authenticate: %w", errors.New("invalid authorization format"))
			ctx.Error(ce.NewError(span, ce.CodeTokenMalformed, ce.MsgUnauthenticated, e))
			ctx.Abort()
			return
		}

		claim, err := utils.JWTParse(auth[1], secret)
		if err != nil {
			e := fmt.Errorf("failed to authenticate: %w", err)

			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				err = ce.NewError(span, ce.CodeTokenExpired, ce.MsgUnauthenticated, e)
			case errors.Is(err, jwt.ErrTokenMalformed):
				err = ce.NewError(span, ce.CodeTokenMalformed, ce.MsgUnauthenticated, e)
			case errors.Is(err, ce.ErrInvalidToken):
				err = ce.NewError(span, ce.CodeInvalidToken, ce.MsgUnauthenticated, e)
			default:
				err = ce.NewError(span, ce.CodeUnknown, ce.MsgInternalServer, e)
			}

			ctx.Error(err)
			ctx.Abort()
			return
		}

		c = context.WithValue(c, constants.CtxKeyAuthID, claim.AuthID)
		c = context.WithValue(c, constants.CtxKeyRole, claim.Role)
		c = context.WithValue(c, constants.CtxKeyIsVerified, claim.IsVerified)

		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

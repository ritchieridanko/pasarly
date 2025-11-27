package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/utils"
)

func NewRequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if strings.TrimSpace(requestID) == "" {
			requestID = utils.NewUUID().String()
		}

		ctx.Writer.Header().Set("X-Request-ID", requestID)
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), constants.CtxKeyRequestID, requestID),
		)

		ctx.Next()
	}
}

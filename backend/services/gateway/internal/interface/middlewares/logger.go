package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Logger(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now().UTC()
		ctx.Next()

		requestID, _ := ctx.Request.Context().Value(constants.CtxKeyRequestID).(string)
		traceID := trace.SpanFromContext(ctx.Request.Context()).SpanContext().TraceID().String()
		status := ctx.Writer.Status()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("trace_id", traceID),
			zap.String("ip_address", ctx.ClientIP()),
			zap.String("user_agent", ctx.Request.UserAgent()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", status),
			zap.String("latency", time.Since(start).String()),
		}

		if status < http.StatusBadRequest {
			l.Base().Info("REQUEST SUCCEEDED", fields...)
		} else {
			l.Base().Warn("REQUEST FAILED", fields...)
		}
	}
}

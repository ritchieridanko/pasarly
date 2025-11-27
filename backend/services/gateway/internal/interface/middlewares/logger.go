package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Logger(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now().UTC()
		ctx.Next()

		requestID, _ := ctx.Request.Context().Value(constants.CtxKeyRequestID).(string)
		traceID := trace.SpanFromContext(ctx.Request.Context()).SpanContext().TraceID().String()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("trace_id", traceID),
			zap.String("ip_address", ctx.ClientIP()),
			zap.String("user_agent", ctx.Request.UserAgent()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("latency", time.Since(start).String()),
		}

		errs := ctx.Errors
		if len(errs) == 0 {
			fields = append(fields, zap.Int("status", ctx.Writer.Status()))

			l.Base().Info("REQUEST_SUCCEEDED", fields...)
			return
		}

		var e *ce.Error
		if errors.As(errs[0].Err, &e) {
			fields = append(
				fields,
				zap.Int("status", e.ToHTTPStatus()),
				zap.String("error_code", string(e.Code)),
				zap.String("error_message", e.Message),
				zap.String("error_detail", e.Error()),
			)

			l.Base().Error("REQUEST_FAILED", fields...)
			utils.SendResponse[any](ctx, e.ToHTTPStatus(), e.Message, nil)
			return
		}

		fields = append(
			fields,
			zap.Int("status", http.StatusInternalServerError),
			zap.String("error_code", string(ce.CodeUnknown)),
			zap.String("error_message", ce.MsgInternalServer),
			zap.String("error_detail", errs[0].Err.Error()),
		)

		l.Base().Error("REQUEST_FAILED", fields...)
		utils.SendResponse[any](ctx, http.StatusInternalServerError, ce.MsgInternalServer, nil)
	}
}

package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/dtos"
)

func SendResponse[T any](ctx *gin.Context, status int, message string, data T) {
	requestID, _ := ctx.Value(constants.CtxKeyRequestID).(string)

	resp := dtos.Response[T]{
		Status:  status,
		Message: message,
		Data:    data,
		Meta: &dtos.Meta{
			RequestID: requestID,
			Timestamp: time.Now().UTC(),
		},
	}

	ctx.JSON(status, resp)
}

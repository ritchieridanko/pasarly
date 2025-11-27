package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/dtos"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/metadata"
)

const authErrTracer string = "handler.auth"

type AuthHandler struct {
	as apis.AuthServiceClient
}

func NewAuthHandler(as apis.AuthServiceClient) *AuthHandler {
	return &AuthHandler{as: as}
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	c, span := otel.Tracer(authErrTracer).Start(ctx.Request.Context(), "SignUp")
	defer span.End()

	var payload dtos.SignUpRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		e := fmt.Errorf("failed to sign up: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, e))
		return
	}

	req := apis.SignUpRequest{
		Email:    payload.Email,
		Password: payload.Password,
	}

	oc := metadata.NewOutgoingContext(c, metadata.Pairs(
		constants.CtxKeyUserAgent, ctx.Request.UserAgent(),
		constants.CtxKeyIPAddress, ctx.ClientIP(),
	))

	resp, err := h.as.SignUp(oc, &req)
	if err != nil {
		ctx.Error(ce.FromGRPCErr(span, err))
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusCreated,
		"Signed up successfully",
		dtos.SignUpResponse{
			AccessToken: resp.GetToken().GetAccess(),
			Auth: dtos.Auth{
				ID:         resp.GetAuth().GetId(),
				Email:      resp.GetAuth().GetEmail(),
				Role:       resp.GetAuth().GetRole(),
				IsVerified: resp.GetAuth().GetIsVerified(),
				CreatedAt:  resp.GetAuth().GetCreatedAt().AsTime(),
				UpdatedAt:  resp.GetAuth().GetUpdatedAt().AsTime(),
			},
		},
	)
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	c, span := otel.Tracer(authErrTracer).Start(ctx.Request.Context(), "SignIn")
	defer span.End()

	var payload dtos.SignInRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		e := fmt.Errorf("failed to sign in: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, e))
		return
	}

	req := apis.SignInRequest{
		Email:    payload.Email,
		Password: payload.Password,
	}

	oc := metadata.NewOutgoingContext(c, metadata.Pairs(
		constants.CtxKeyUserAgent, ctx.Request.UserAgent(),
		constants.CtxKeyIPAddress, ctx.ClientIP(),
	))

	resp, err := h.as.SignIn(oc, &req)
	if err != nil {
		ctx.Error(ce.FromGRPCErr(span, err))
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		"Signed in successfully",
		dtos.SignInResponse{
			AccessToken: resp.GetToken().GetAccess(),
			Auth: dtos.Auth{
				ID:         resp.GetAuth().GetId(),
				Email:      resp.GetAuth().GetEmail(),
				Role:       resp.GetAuth().GetRole(),
				IsVerified: resp.GetAuth().GetIsVerified(),
				CreatedAt:  resp.GetAuth().GetCreatedAt().AsTime(),
				UpdatedAt:  resp.GetAuth().GetUpdatedAt().AsTime(),
			},
		},
	)
}

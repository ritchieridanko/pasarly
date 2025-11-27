package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
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
	config *configs.Config
	as     apis.AuthServiceClient
	cookie *utils.Cookie
}

func NewAuthHandler(cfg *configs.Config, as apis.AuthServiceClient, c *utils.Cookie) *AuthHandler {
	return &AuthHandler{config: cfg, as: as, cookie: c}
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

	h.cookie.Set(
		ctx,
		constants.CookieKeySession,
		resp.GetToken().Session,
		h.config.Duration.Session,
		"/",
	)

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

	h.cookie.Set(
		ctx,
		constants.CookieKeySession,
		resp.GetToken().Session,
		h.config.Duration.Session,
		"/",
	)

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

func (h *AuthHandler) SignOut(ctx *gin.Context) {
	c, span := otel.Tracer(authErrTracer).Start(ctx.Request.Context(), "SignOut")
	defer span.End()

	session, err := ctx.Cookie(constants.CookieKeySession)
	e := fmt.Errorf("failed to sign out: %w", err)

	if errors.Is(err, http.ErrNoCookie) {
		ctx.Error(ce.NewError(span, ce.CodeCookieNotFound, ce.MsgUnauthenticated, e))
		return
	}
	if err != nil {
		ctx.Error(ce.NewError(span, ce.CodeInternal, ce.MsgInternalServer, e))
		return
	}
	if session == "" {
		e := fmt.Errorf("failed to sign out: %w", http.ErrNoCookie)
		ctx.Error(ce.NewError(span, ce.CodeCookieNotFound, ce.MsgUnauthenticated, e))
		return
	}

	_, err = h.as.SignOut(c, &apis.SignOutRequest{Session: session})
	if err != nil {
		ctx.Error(ce.FromGRPCErr(span, err))
		return
	}

	h.cookie.Unset(ctx, constants.CookieKeySession, "/")
	utils.SendResponse[any](ctx, http.StatusNoContent, "", nil)
}

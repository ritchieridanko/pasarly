package handlers

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/usecases"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const authErrTracer string = "handler.auth"

type AuthHandler struct {
	apis.UnimplementedAuthServiceServer
	au     usecases.AuthUsecase
	su     usecases.SessionUsecase
	logger *logger.Logger
}

func NewAuthHandler(au usecases.AuthUsecase, su usecases.SessionUsecase, l *logger.Logger) *AuthHandler {
	return &AuthHandler{au: au, su: su, logger: l}
}

func (h *AuthHandler) SignUp(ctx context.Context, req *apis.SignUpRequest) (*apis.SignUpResponse, error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "SignUp")
	defer span.End()

	pass := req.GetPassword()
	data := models.CreateAuth{
		Email:    req.GetEmail(),
		Password: &pass,
	}

	auth, err := h.au.SignUp(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	resp := apis.SignUpResponse{
		Auth: &apis.Auth{
			Id:         auth.ID,
			Email:      auth.Email,
			Role:       auth.Role,
			IsVerified: auth.IsVerified,
			CreatedAt:  timestamppb.New(auth.CreatedAt),
			UpdatedAt:  timestamppb.New(auth.UpdatedAt),
		},
	}

	ua, ip := utils.RequestMeta(ctx)
	if ua == "" || ip == "" {
		w := fmt.Sprintf("invalid request metadata (user_agent=%s, ip_address=%s)", ua, ip)
		h.logger.Sugar().Warnf("failed to create session: %s", w)
		return &resp, nil
	}

	rm := models.RequestMeta{
		UserAgent: ua,
		IPAddress: ip,
	}

	authToken, err := h.su.CreateSession(ctx, auth, &rm)
	if err != nil {
		h.logger.Sugar().Warnln(err.Error())
		return &resp, nil
	}

	resp.Token = &apis.AuthToken{
		Session: authToken.Session,
		Access:  authToken.Access,
	}

	return &resp, nil
}

func (h *AuthHandler) SignIn(ctx context.Context, req *apis.SignInRequest) (*apis.SignInResponse, error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "SignIn")
	defer span.End()

	data := models.GetAuth{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	auth, err := h.au.SignIn(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	ua, ip := utils.RequestMeta(ctx)
	if ua == "" || ip == "" {
		w := fmt.Sprintf("invalid request metadata (user_agent=%s, ip_address=%s)", ua, ip)
		h.logger.Sugar().Errorf("failed to create session: %s", w)
		return nil, status.Error(codes.Internal, ce.MsgInternalServer)
	}

	rm := models.RequestMeta{
		UserAgent: ua,
		IPAddress: ip,
	}

	authToken, err := h.su.CreateSession(ctx, auth, &rm)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	return &apis.SignInResponse{
		Token: &apis.AuthToken{
			Session: authToken.Session,
			Access:  authToken.Access,
		},
		Auth: &apis.Auth{
			Id:         auth.ID,
			Email:      auth.Email,
			Role:       auth.Role,
			IsVerified: auth.IsVerified,
			CreatedAt:  timestamppb.New(auth.CreatedAt),
			UpdatedAt:  timestamppb.New(auth.UpdatedAt),
		},
	}, nil
}

func (h *AuthHandler) SignOut(ctx context.Context, req *apis.SignOutRequest) (*emptypb.Empty, error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "SignOut")
	defer span.End()

	if err := h.su.RevokeSession(ctx, req.GetSession()); err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	return &emptypb.Empty{}, nil
}

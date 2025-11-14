package handlers

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/pasarly/auth-service/internal/app/models"
	"github.com/ritchieridanko/pasarly/auth-service/internal/app/usecases"
	"github.com/ritchieridanko/pasarly/auth-service/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/auth-service/internal/interface/grpc/protobufs/v1"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/ce"
	"github.com/ritchieridanko/pasarly/auth-service/internal/shared/constants"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const authErrTracer string = "grpc.handler.auth"

type AuthGRPCHandler struct {
	protobufs.UnimplementedAuthServiceServer
	au     usecases.AuthUsecase
	su     usecases.SessionUsecase
	logger *logger.Logger
}

func NewAuthGRPCHandler(au usecases.AuthUsecase, su usecases.SessionUsecase, l *logger.Logger) *AuthGRPCHandler {
	return &AuthGRPCHandler{au: au, su: su, logger: l}
}

func (h *AuthGRPCHandler) SignUp(ctx context.Context, req *protobufs.SignUpRequest) (*protobufs.SignUpResponse, error) {
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
		return nil, err.GRPCStatus()
	}

	resp := protobufs.SignUpResponse{
		Auth: &protobufs.Auth{
			Id:         auth.ID,
			Email:      auth.Email,
			Role:       auth.Role,
			IsVerified: auth.IsVerified,
			CreatedAt:  timestamppb.New(auth.CreatedAt),
			UpdatedAt:  timestamppb.New(auth.UpdatedAt),
		},
	}

	ua, _ := ctx.Value(constants.CtxKeyUserAgent).(string)
	ip, _ := ctx.Value(constants.CtxKeyIPAddress).(string)
	if ua == "" || ip == "" {
		w := fmt.Sprintf("invalid request metadata (user_agent=%s, ip_address=%s)", ua, ip)
		h.logger.Sugar().Warnf("failed to sign up: %s", w)
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

	resp.Token = &protobufs.AuthToken{
		Session: authToken.Session,
		Access:  authToken.Access,
	}

	return &resp, nil
}

func (h *AuthGRPCHandler) SignIn(ctx context.Context, req *protobufs.SignInRequest) (*protobufs.SignInResponse, error) {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "SignIn")
	defer span.End()

	data := models.GetAuth{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	auth, err := h.au.SignIn(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.GRPCStatus()
	}

	ua, _ := ctx.Value(constants.CtxKeyUserAgent).(string)
	ip, _ := ctx.Value(constants.CtxKeyIPAddress).(string)
	if ua == "" || ip == "" {
		w := fmt.Sprintf("invalid request metadata (user_agent=%s, ip_address=%s)", ua, ip)
		h.logger.Sugar().Errorf("failed to sign in: %s", w)
		return nil, status.Error(codes.Internal, ce.MsgInternalServer)
	}

	rm := models.RequestMeta{
		UserAgent: ua,
		IPAddress: ip,
	}

	authToken, err := h.su.CreateSession(ctx, auth, &rm)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.GRPCStatus()
	}

	return &protobufs.SignInResponse{
		Token: &protobufs.AuthToken{
			Session: authToken.Session,
			Access:  authToken.Access,
		},
		Auth: &protobufs.Auth{
			Id:         auth.ID,
			Email:      auth.Email,
			Role:       auth.Role,
			IsVerified: auth.IsVerified,
			CreatedAt:  timestamppb.New(auth.CreatedAt),
			UpdatedAt:  timestamppb.New(auth.UpdatedAt),
		},
	}, nil
}

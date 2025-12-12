package handlers

import (
	"context"

	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/usecases"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const userErrTracer string = "handler.user"

type UserHandler struct {
	apis.UnimplementedUserServiceServer
	uu     usecases.UserUsecase
	logger *logger.Logger
}

func NewUserHandler(uu usecases.UserUsecase, l *logger.Logger) *UserHandler {
	return &UserHandler{uu: uu, logger: l}
}

func (h *UserHandler) UpsertUser(ctx context.Context, req *apis.UpsertUserRequest) (*apis.UpsertUserResponse, error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "UpsertUser")
	defer span.End()

	data := models.UpsertUser{
		AuthID:    req.GetAuthId(),
		Name:      req.GetName(),
		Bio:       utils.UnwrapString(req.GetBio()),
		Sex:       utils.UnwrapString(req.GetSex()),
		Birthdate: utils.UnwrapTimestamp(req.GetBirthdate()),
		Phone:     utils.UnwrapString(req.GetPhone()),
	}

	user, err := h.uu.UpsertUser(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	return &apis.UpsertUserResponse{User: h.toUser(user)}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *apis.GetUserRequest) (*apis.GetUserResponse, error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "GetUser")
	defer span.End()

	user, err := h.uu.GetUser(ctx, req.GetAuthId())
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	return &apis.GetUserResponse{User: h.toUser(user)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *apis.UpdateUserRequest) (*apis.UpdateUserResponse, error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "UpdateUser")
	defer span.End()

	data := models.UpdateUser{
		AuthID:    req.GetAuthId(),
		Name:      utils.UnwrapString(req.GetName()),
		Bio:       utils.UnwrapString(req.GetBio()),
		Sex:       utils.UnwrapString(req.GetSex()),
		Birthdate: utils.UnwrapTimestamp(req.GetBirthdate()),
		Phone:     utils.UnwrapString(req.GetPhone()),
	}

	user, err := h.uu.UpdateUser(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	return &apis.UpdateUserResponse{User: h.toUser(user)}, nil
}

func (h *UserHandler) UpdateProfilePicture(ctx context.Context, req *apis.UpdateProfilePictureRequest) (*apis.UpdateProfilePictureResponse, error) {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "UpdateProfilePicture")
	defer span.End()

	data := models.UpdateProfilePicture{
		AuthID:         req.GetAuthId(),
		ProfilePicture: req.GetProfilePicture(),
	}

	profilePicture, err := h.uu.UpdateProfilePicture(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	return &apis.UpdateProfilePictureResponse{ProfilePicture: profilePicture}, nil
}

func (h *UserHandler) toUser(u *models.User) *apis.User {
	user := apis.User{
		Id:             u.ID,
		Name:           u.Name,
		Bio:            utils.WrapString(u.Bio),
		Sex:            utils.WrapString(u.Sex),
		Birthdate:      utils.WrapTime(u.Birthdate),
		Phone:          utils.WrapString(u.Phone),
		ProfilePicture: utils.WrapString(u.ProfilePicture),
		CreatedAt:      timestamppb.New(u.CreatedAt),
		UpdatedAt:      timestamppb.New(u.UpdatedAt),
	}
	return &user
}

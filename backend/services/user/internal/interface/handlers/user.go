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
		Bio:       utils.UnwrapString(req.Bio),
		Sex:       utils.UnwrapString(req.Sex),
		Birthdate: utils.UnwrapTimestamp(req.Birthdate),
		Phone:     utils.UnwrapString(req.Phone),
	}

	user, err := h.uu.UpsertUser(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	return &apis.UpsertUserResponse{
		User: &apis.User{
			Id:             user.ID,
			Name:           user.Name,
			Bio:            utils.WrapString(user.Bio),
			Sex:            utils.WrapString(user.Sex),
			Birthdate:      utils.WrapTime(user.Birthdate),
			Phone:          utils.WrapString(user.Phone),
			ProfilePicture: utils.WrapString(user.ProfilePicture),
			CreatedAt:      timestamppb.New(user.CreatedAt),
			UpdatedAt:      timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

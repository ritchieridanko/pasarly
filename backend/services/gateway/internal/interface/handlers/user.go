package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/interface/dtos"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const userErrTracer string = "handler.user"

type UserHandler struct {
	us apis.UserServiceClient
}

func NewUserHandler(us apis.UserServiceClient) *UserHandler {
	return &UserHandler{us: us}
}

func (h *UserHandler) UpsertUser(ctx *gin.Context) {
	c, span := otel.Tracer(userErrTracer).Start(ctx.Request.Context(), "UpsertUser")
	defer span.End()

	var payload dtos.UpsertUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		e := fmt.Errorf("failed to upsert user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, e))
		return
	}

	authID, err := utils.CtxAuthID(c)
	if err != nil {
		e := fmt.Errorf("failed to upsert user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeCtxValueNotFound, ce.MsgInternalServer, e))
		return
	}

	req := apis.UpsertUserRequest{
		AuthId:    authID,
		Name:      payload.Name,
		Bio:       utils.WrapString(payload.Bio),
		Sex:       utils.WrapString(payload.Sex),
		Birthdate: utils.WrapTime(payload.Birthdate),
		Phone:     utils.WrapString(payload.Phone),
	}

	resp, err := h.us.UpsertUser(c, &req)
	if err != nil {
		ctx.Error(ce.FromGRPCErr(span, err))
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		"User created successfully",
		dtos.UpsertUserResponse{
			User: dtos.User{
				ID:             resp.GetUser().GetId(),
				Name:           resp.GetUser().GetName(),
				Bio:            utils.UnwrapString(resp.GetUser().GetBio()),
				Sex:            utils.UnwrapString(resp.GetUser().GetSex()),
				Birthdate:      utils.UnwrapTimestamp(resp.GetUser().GetBirthdate()),
				Phone:          utils.UnwrapString(resp.GetUser().GetPhone()),
				ProfilePicture: utils.UnwrapString(resp.GetUser().GetProfilePicture()),
				CreatedAt:      resp.GetUser().GetCreatedAt().AsTime(),
				UpdatedAt:      resp.GetUser().GetUpdatedAt().AsTime(),
			},
		},
	)
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	c, span := otel.Tracer(userErrTracer).Start(ctx.Request.Context(), "GetUser")
	defer span.End()

	authID, err := utils.CtxAuthID(c)
	if err != nil {
		e := fmt.Errorf("failed to fetch user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeCtxValueNotFound, ce.MsgInternalServer, e))
		return
	}

	resp, err := h.us.GetUser(c, &apis.GetUserRequest{AuthId: authID})
	if err != nil {
		ctx.Error(ce.FromGRPCErr(span, err))
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		"OK",
		dtos.GetUserResponse{
			User: dtos.User{
				ID:             resp.GetUser().GetId(),
				Name:           resp.GetUser().GetName(),
				Bio:            utils.UnwrapString(resp.GetUser().GetBio()),
				Sex:            utils.UnwrapString(resp.GetUser().GetSex()),
				Birthdate:      utils.UnwrapTimestamp(resp.GetUser().GetBirthdate()),
				Phone:          utils.UnwrapString(resp.GetUser().GetPhone()),
				ProfilePicture: utils.UnwrapString(resp.GetUser().GetProfilePicture()),
				CreatedAt:      resp.GetUser().GetCreatedAt().AsTime(),
				UpdatedAt:      resp.GetUser().GetUpdatedAt().AsTime(),
			},
		},
	)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	c, span := otel.Tracer(userErrTracer).Start(ctx.Request.Context(), "UpdateUser")
	defer span.End()

	var payload dtos.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		e := fmt.Errorf("failed to update user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, e))
		return
	}

	authID, err := utils.CtxAuthID(c)
	if err != nil {
		e := fmt.Errorf("failed to update user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeCtxValueNotFound, ce.MsgInternalServer, e))
		return
	}

	req := apis.UpdateUserRequest{
		AuthId:    authID,
		Name:      utils.WrapString(payload.Name),
		Bio:       utils.WrapString(payload.Bio),
		Sex:       utils.WrapString(payload.Sex),
		Birthdate: utils.WrapTime(payload.Birthdate),
		Phone:     utils.WrapString(payload.Phone),
	}

	resp, err := h.us.UpdateUser(c, &req)
	if err != nil {
		ctx.Error(ce.FromGRPCErr(span, err))
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		"User updated successfully",
		dtos.UpdateUserResponse{
			User: dtos.User{
				ID:             resp.GetUser().GetId(),
				Name:           resp.GetUser().GetName(),
				Bio:            utils.UnwrapString(resp.GetUser().GetBio()),
				Sex:            utils.UnwrapString(resp.GetUser().GetSex()),
				Birthdate:      utils.UnwrapTimestamp(resp.GetUser().GetBirthdate()),
				Phone:          utils.UnwrapString(resp.GetUser().GetPhone()),
				ProfilePicture: utils.UnwrapString(resp.GetUser().GetProfilePicture()),
				CreatedAt:      resp.GetUser().GetCreatedAt().AsTime(),
				UpdatedAt:      resp.GetUser().GetUpdatedAt().AsTime(),
			},
		},
	)
}

func (h *UserHandler) UpdateProfilePicture(ctx *gin.Context) {
	c, span := otel.Tracer(userErrTracer).Start(ctx.Request.Context(), "UpdateProfilePicture")
	defer span.End()

	var payload dtos.UpdateProfilePictureRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		e := fmt.Errorf("failed to update profile picture: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, e))
		return
	}

	authID, err := utils.CtxAuthID(c)
	if err != nil {
		e := fmt.Errorf("failed to update profile picture: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeCtxValueNotFound, ce.MsgInternalServer, e))
		return
	}

	req := apis.UpdateProfilePictureRequest{
		AuthId:         authID,
		ProfilePicture: payload.ProfilePicture,
	}

	resp, err := h.us.UpdateProfilePicture(c, &req)
	if err != nil {
		ctx.Error(ce.FromGRPCErr(span, err))
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		"Profile picture updated successfully",
		dtos.UpdateProfilePictureResponse{
			ProfilePicture: resp.GetProfilePicture(),
		},
	)
}

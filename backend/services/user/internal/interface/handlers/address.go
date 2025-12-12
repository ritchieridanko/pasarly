package handlers

import (
	"context"
	"strings"

	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/logger"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/usecases"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/apis/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const addressErrTracer string = "handler.address"

type AddressHandler struct {
	apis.UnimplementedUserAddressServiceServer
	au     usecases.AddressUsecase
	logger *logger.Logger
}

func NewAddressHandler(au usecases.AddressUsecase, l *logger.Logger) *AddressHandler {
	return &AddressHandler{au: au, logger: l}
}

func (h *AddressHandler) CreateAddress(ctx context.Context, req *apis.CreateUserAddressRequest) (*apis.CreateUserAddressResponse, error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "CreateAddress")
	defer span.End()

	data := models.CreateAddress{
		AuthID:       req.GetAuthId(),
		Recipient:    strings.TrimSpace(req.GetRecipient()),
		Phone:        strings.TrimSpace(req.GetPhone()),
		Label:        strings.TrimSpace(req.GetLabel()),
		Notes:        utils.UnwrapString(req.GetNotes()),
		IsPrimary:    req.GetIsPrimary(),
		Country:      utils.NormalizeString(req.GetCountry()),
		Subdivision1: utils.NormalizeStringPtr(utils.UnwrapString(req.GetSubdivision_1())),
		Subdivision2: utils.NormalizeStringPtr(utils.UnwrapString(req.GetSubdivision_2())),
		Subdivision3: utils.NormalizeStringPtr(utils.UnwrapString(req.GetSubdivision_3())),
		Subdivision4: utils.NormalizeStringPtr(utils.UnwrapString(req.GetSubdivision_4())),
		Street:       strings.TrimSpace(req.GetStreet()),
		Postcode:     strings.TrimSpace(req.GetPostcode()),
		Latitude:     req.GetLatitude(),
		Longitude:    req.GetLongitude(),
	}

	address, opa, err := h.au.CreateAddress(ctx, &data)
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	var oldPrimaryAddress *apis.UserAddress
	if opa != nil {
		oldPrimaryAddress = h.toAddress(opa)
	}

	return &apis.CreateUserAddressResponse{
		Address:           h.toAddress(address),
		OldPrimaryAddress: oldPrimaryAddress,
	}, nil
}

func (h *AddressHandler) GetAllAddresses(ctx context.Context, req *apis.GetAllUserAddressesRequest) (*apis.GetAllUserAddressesResponse, error) {
	ctx, span := otel.Tracer(addressErrTracer).Start(ctx, "GetAllAddresses")
	defer span.End()

	addresses, err := h.au.GetAllAddresses(ctx, req.GetAuthId())
	if err != nil {
		h.logger.Sugar().Errorln(err.Error())
		return nil, err.ToGRPCStatus()
	}

	addrs := make([]*apis.UserAddress, 0, len(addresses))
	for _, address := range addresses {
		addr := h.toAddress(&address)
		addrs = append(addrs, addr)
	}

	return &apis.GetAllUserAddressesResponse{Addresses: addrs}, nil
}

func (h *AddressHandler) toAddress(a *models.Address) *apis.UserAddress {
	address := apis.UserAddress{
		Id:            a.ID,
		Recipient:     a.Recipient,
		Phone:         a.Phone,
		Label:         a.Label,
		Notes:         utils.WrapString(a.Notes),
		IsPrimary:     a.IsPrimary,
		Country:       utils.ToTitlecase(a.Country),
		Subdivision_1: utils.WrapString(utils.ToTitlecasePtr(a.Subdivision1)),
		Subdivision_2: utils.WrapString(utils.ToTitlecasePtr(a.Subdivision2)),
		Subdivision_3: utils.WrapString(utils.ToTitlecasePtr(a.Subdivision3)),
		Subdivision_4: utils.WrapString(utils.ToTitlecasePtr(a.Subdivision4)),
		Street:        a.Street,
		Postcode:      a.Postcode,
		Latitude:      a.Latitude,
		Longitude:     a.Longitude,
		CreatedAt:     timestamppb.New(a.CreatedAt),
		UpdatedAt:     timestamppb.New(a.UpdatedAt),
	}
	return &address
}

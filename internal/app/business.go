package app

import (
	"context"
	"fmt"

	"ownify_api/internal/constants"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	desc "ownify_api/pkg"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer) CreateBusiness(ctx context.Context, req *desc.CreateBusinessRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}

	business := dto.BriefBusiness{
		UserId:      *uid,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Pin:         utils.Hash(req.Password),
		Business:    req.Business,
		PhoneNumber: req.PhoneNumber,
		Location:    req.Location,
	}
	if !business.Valid() {
		return nil, fmt.Errorf(constants.ErrInvalidBusinessInfo, "raw message: %s", req.Email)
	}

	isRegistered := m.authService.VerifyBusinessByUserId(*uid)
	if isRegistered {
		err := m.businessService.UpdateBusiness(&business)
		if err != nil {
			return nil, err
		}
		return &desc.NetWorkResponse{Success: true, Msg: "Successfully updated!"}, nil
	}

	customerId, err := m.paymentService.CreateCustomer(req.Email)
	if err != nil {
		return nil, err
	}
	err = m.businessService.CreateBusiness(&business)
	if err != nil {
		return nil, err
	}
	type CreateBusinessRes struct {
		CustomerId string `json:"customer_id"`
	}
	res := CreateBusinessRes{CustomerId: customerId}
	return BuildRes(res, "Successfully created!", true)
}

func (m *MicroserviceServer) DeleteBusiness(ctx context.Context, req *desc.DeleteBusinessRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ERR] no permission to user: %s", req.Email)
	}

	if !m.authService.ValidBusiness(*uid, req.Email) {
		return nil, fmt.Errorf("[ERR] no permission to user: %s", req.Email)
	}

	err = m.businessService.DeleteBusiness(req.Email, *uid)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{Success: true, Msg: "Successfully deleted."}, nil
}

func (m *MicroserviceServer) GetBusiness(ctx context.Context, req *desc.GetBusinessRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	err = utils.IsEmail(req.Email)
	if err != nil {
		return nil, err
	}

	data, err := m.businessService.GetBusiness(req.Email)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrNotFoundBusiness, "raw message:%s", err)
	}

	// check subscription status
	isSubscription := m.paymentService.VerifySubscriptionStatus(req.Email)
	data.IsSubscribed = isSubscription
	return BuildRes(data, "Here is your business info", true)
}

func (m *MicroserviceServer) GetBusinessByPubAddr(ctx context.Context, req *desc.GetBusinessWithPubAddrRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}

	err = utils.IsPubKey(req.PubAddr)
	if err != nil {
		return nil, err
	}

	data, err := m.businessService.GetBusinessByWalletAddress(req.PubAddr)
	if err != nil {
		return nil, err
	}

	return BuildRes(data, "Here is your business info", true)
}

func (m *MicroserviceServer) GetBusinessByUserId(ctx context.Context, req *emptypb.Empty) (*desc.NetWorkResponse, error) {

	// validate token.
	userId, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}

	data, err := m.businessService.GetBusinessByUserId(*userId)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrNotFoundBusiness, "raw message:%s", err)
	}

	return BuildRes(data, "Here is your business info", true)
}

func (m *MicroserviceServer) GenerateNewAPIKey(ctx context.Context, req *desc.NewLicenseRequest) (*desc.NetWorkResponse, error) {
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}
	apiKey, err := m.licenseService.GenerateAPIKey(req.Email, *uid)
	if err != nil {
		return nil, err
	}
	return BuildRes(apiKey, "successfully generated", true)
}

func (m *MicroserviceServer) GetApiKey(ctx context.Context, req *desc.LicenseRequest) (*desc.NetWorkResponse, error) {
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}
	apiKeys, err := m.licenseService.GetAPIKey(req.Email, *uid)
	if err != nil {
		return nil, fmt.Errorf(constants.WarningNotFoundApiKey, "raw message:%s", err)
	}
	type ApiKeyRes struct {
		ApiKeyRes []string `json:"api_keys"`
	}
	res := ApiKeyRes{
		ApiKeyRes: apiKeys,
	}
	return BuildRes(res, "you have api key already.", true)
}

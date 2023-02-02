package app

import (
	"context"
	"fmt"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	desc "ownify_api/pkg"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer) CreateWallet(ctx context.Context, req *desc.CreateWalletRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	if utils.IsEmpty(req.Email) {
		return nil, fmt.Errorf("[ERR] invalid request: %v", req)
	}
	if !m.authService.ValidBusiness(*uid, req.Email) {
		return nil, err
	}
	// create new wallet.
	pubKey, err := m.walletService.AddNewAccount(
		req.UserRole,
		req.Email,
	)
	if err != nil {
		return nil, err
	}

	return BuildRes(dto.Wallet{
		UserRole: req.UserRole,
		Email:    req.Email,
		PubKey:   *pubKey,
	}, "Successfully created", true)

}

func (m *MicroserviceServer) GetMyAccounts(ctx context.Context, req *desc.SignInRequest) (*emptypb.Empty, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	if utils.IsEmpty(req.Email) {
		return nil, fmt.Errorf("[ERR] invalid request: %v", req)
	}
	if !m.authService.ValidBusiness(*uid, req.Email) {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (m *MicroserviceServer) MintOwnify(ctx context.Context, req *desc.MintOwnifyRequest) (*desc.MintOwnifyResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.MintOwnifyResponse{}, err
	}
	if !m.authService.ValidBusiness(*uid, req.Email) {
		return &desc.MintOwnifyResponse{}, err
	}

	ownifyProducts := []dto.BriefProduct{}
	for _, item := range req.Products {
		product := dto.BriefProduct{
			Barcode:        item.Barcode,
			BrandName:      item.BrandName,
			ItemName:       item.ItemName,
			AdditionalData: item.AdditionalData,
			Location:       item.Location,
			IssuedDate:     item.IssuedDate,
		}
		ownifyProducts = append(ownifyProducts, product)
	}

	// mint ownify assets.
	assetIds, err := m.walletService.MintOwnify(req.Email, req.PubKey, ownifyProducts, req.Net)
	if err != nil {
		return &desc.MintOwnifyResponse{}, nil
	}
	return &desc.MintOwnifyResponse{AssetIds: assetIds}, nil
}

func (m *MicroserviceServer) MakeTransaction(ctx context.Context, req *desc.MakeTransactionRequest) (*desc.MakeTransactionResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.MakeTransactionResponse{}, err
	}

	if !m.authService.ValidBusiness(*uid, req.Email) {
		return &desc.MakeTransactionResponse{}, err
	}

	txId, err := m.walletService.MakeTransaction(
		req.Role,
		req.Email,
		req.PubKey,
		req.RawTx,
		req.Net,
	)
	if err != nil {
		return &desc.MakeTransactionResponse{}, err
	}

	return &desc.MakeTransactionResponse{TxId: *txId}, nil
}

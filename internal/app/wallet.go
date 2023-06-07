package app

import (
	"context"
	"fmt"
	"ownify_api/internal/constants"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	desc "ownify_api/pkg"
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

	acc, err := m.walletService.GetMyAccounts(req.Email, *uid)
	if len(acc) != 0 || err != nil {
		return BuildRes(dto.Wallet{
			UserRole: req.UserRole,
			Email:    req.Email,
			PubKey:   acc[0],
		}, "Successfully created", true)
	}
	// create new wallet.
	pubKey, err := m.walletService.AddNewAccount(
		req.UserRole,
		*uid,
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

func (m *MicroserviceServer) RegisterWallet(ctx context.Context, req *desc.RegisterWalletRequest) (*desc.NetWorkResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	if utils.IsPubKey(req.WalletAddress) != nil {
		return nil, fmt.Errorf("[ERR] invalid request: %v", req)
	}

	if !utils.VerifySignature([]byte(req.WalletAddress), req.WalletAddress, req.Sig) {
		return nil, fmt.Errorf("[ERR] can't register other's wallet: %s", req.WalletAddress)
	}

	// register new wallet.
	pubKey, err := m.walletService.RegisterNewAccount(req.WalletAddress, *uid)
	if err != nil {
		return nil, err
	}
	return BuildRes(dto.Wallet{
		PubKey: *pubKey,
	}, "Successfully registered", true)

}

func (m *MicroserviceServer) SendOwnify(ctx context.Context, req *desc.SendOwnifyRequest) (*desc.NetWorkResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	if utils.IsPubKey(req.Sender) != nil || utils.IsPubKey(req.Receiver) != nil {
		return nil, fmt.Errorf("[ERR] invalid request(invalid address): %v", req)
	}

	if len(req.AssetIds) == 0 {
		return nil, fmt.Errorf("[ERR] invalid request(empty asset list): %v", req)
	}

	email, err := m.authService.GetEmailByUserId(*uid)
	if err != nil {
		return nil, fmt.Errorf("[ERR] invalid request(business doesn't exist): %v", err)
	}

	// send assets.
	txId, err := m.walletService.SendOwnify(*email, req.AssetIds, req.Sender, req.Receiver, req.Net)
	if err != nil {
		return nil, err
	}
	return BuildRes(txId, "Successfully sent", true)
}

func (m *MicroserviceServer) DeleteOwnify(ctx context.Context, req *desc.DeleteOwnifyRequest) (*desc.NetWorkResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	if utils.IsPubKey(req.Owner) != nil {
		return nil, fmt.Errorf("[ERR] invalid request(invalid address): %v", req)
	}

	if len(req.AssetIds) == 0 {
		return nil, fmt.Errorf("[ERR] invalid request(empty asset list): %v", req)
	}

	email, err := m.authService.GetEmailByUserId(*uid)
	if err != nil {
		return nil, fmt.Errorf("[ERR] invalid request(business doesn't exist): %v", err)
	}

	// delete assets.
	txId, err := m.walletService.DeleteOwnify(*email, req.AssetIds, req.Owner, req.Net)
	if err != nil {
		return nil, err
	}
	return BuildRes(txId, "Successfully deleted", true)
}

func (m *MicroserviceServer) GetMyAccounts(ctx context.Context, req *desc.SignInRequest) (*desc.GetAccountResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}
	if utils.IsEmpty(req.Email) {
		return nil, fmt.Errorf(constants.ErrInvalidRequest, "raw message: %v", req)
	}
	if !m.authService.ValidBusiness(*uid, req.Email) {
		return nil, err
	}

	return &desc.GetAccountResponse{}, nil
}

func (m *MicroserviceServer) GetAccounts(ctx context.Context, req *desc.GetWalletsRequest) (*desc.NetWorkResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}
	if utils.IsEmpty(req.Email) {
		return nil, fmt.Errorf(constants.ErrInvalidRequest, "raw message: %v", req)
	}
	if !m.authService.ValidBusiness(*uid, req.Email) {
		return nil, fmt.Errorf(constants.ErrNotFoundBusiness, "raw message:%s", err)
	}

	wallets, err := m.walletService.GetMyAccounts(req.Email, *uid)
	if err != nil {
		return nil, fmt.Errorf(constants.WarningNotFoundWallet, "raw message:%s", err)
	}

	type AccountRes struct {
		Wallets []string
	}
	res := AccountRes{
		Wallets: wallets,
	}
	return BuildRes(res, "Here is your wallets", true)
}

func (m *MicroserviceServer) MintOwnify(ctx context.Context, req *desc.MintOwnifyRequest) (*desc.MintOwnifyResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}
	if !m.authService.ValidBusiness(*uid, req.Email) {
		return &desc.MintOwnifyResponse{}, err
	}

	ownifyProducts := []dto.BriefProduct{}
	for _, item := range req.Products {
		product := dto.BriefProduct{
			ChainId:        int(item.ChainId),
			Barcode:        item.Barcode,
			BrandName:      item.BrandName,
			ItemName:       item.ItemName,
			AdditionalData: item.AdditionalData,
			Location:       item.Location,
			IssuedDate:     item.IssuedDate,
		}
		// set default value Algorand.
		if item.ChainId == 0 {
			item.ChainId = 1
		}
		ownifyProducts = append(ownifyProducts, product)
	}

	// mint ownify assets.
	assetIds, err := m.walletService.MintOwnify(req.Email, req.PubKey, ownifyProducts, req.Net)
	if err != nil {
		return &desc.MintOwnifyResponse{}, err
	}
	return &desc.MintOwnifyResponse{AssetIds: assetIds}, nil
}

func (m *MicroserviceServer) MakeTransaction(ctx context.Context, req *desc.MakeTransactionRequest) (*desc.MakeTransactionResponse, error) {
	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUser, "raw message:%s", err)
	}

	if !m.authService.VerifyBusinessByUserId(*uid) {
		return &desc.MakeTransactionResponse{}, err
	}

	txId, _, err := m.walletService.MakeTx(
		req.RawTx,
		req.Net,
	)
	if err != nil {
		return &desc.MakeTransactionResponse{}, err
	}

	return &desc.MakeTransactionResponse{TxId: *txId}, nil
}

package app

import (
	"context"
	"fmt"
	"ownify_api/internal/dto"
	desc "ownify_api/pkg"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer) CreateWallet(ctx context.Context, req *desc.CreateWalletRequest) (*desc.CreateWalletResponse, error) {

	// validate token.
	err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.CreateWalletResponse{}, err
	}

	// create new wallet.
	pubKey, err := m.walletService.AddNewAccount(
		req.Role,
		req.Email,
	)
	if err != nil {
		return &desc.CreateWalletResponse{}, nil
	}
	return &desc.CreateWalletResponse{Wallet: *pubKey}, nil
}

func (m *MicroserviceServer) GetMyAccounts(ctx context.Context, req *desc.SignInRequest) (*emptypb.Empty, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	// fmt.Println(md.Get("test"))
	// token, err := m.authService.SignIn(req.GetEmail(), req.GetPassword())
	// if err != nil {
	// 	return nil, err
	// }
	// err = grpc.SendHeader(ctx, metadata.New(map[string]string{
	// 	"Token": *token,
	// }))

	// if err != nil {
	// 	return nil, err
	// }
	return &emptypb.Empty{}, nil
}

func (m *MicroserviceServer) MintOwnify(ctx context.Context, req *desc.MintOwnifyRequest) (*desc.MintOwnifyResponse, error) {
	// validate token.
	err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.MintOwnifyResponse{}, err
	}

	ownifyProducts := []dto.BriefProduct{}
	for _, item := range req.Products {
		product := dto.BriefProduct{
			BrandName:      item.BrandName,
			ItemName:       item.ItemName,
			AdditionalData: item.AdditionalData,
			Location:       item.Location,
			IssueDate:      fmt.Sprintf("%v", item.IssueDate),
		}
		ownifyProducts = append(ownifyProducts, product)
	}

	// mint ownify assets.
	assetIds, err := m.walletService.MintOwnify(req.PubKey, ownifyProducts, req.Net)
	if err != nil {
		return &desc.MintOwnifyResponse{}, nil
	}
	return &desc.MintOwnifyResponse{AssetIds: assetIds}, nil
}

func (m *MicroserviceServer) MakeTransaction(ctx context.Context, req *desc.MakeTransactionRequest) (*desc.MakeTransactionResponse, error) {
	// validate token.
	err := m.TokenInterceptor(ctx)
	if err != nil {
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

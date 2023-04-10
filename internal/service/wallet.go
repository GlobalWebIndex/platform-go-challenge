package service

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type WalletService interface {
	AddNewAccount(role string, userId string, email string) (*string, error)
	RegisterNewAccount(walletAddress string, userId string) (*string, error)
	GetMyAccounts(role string, userId string) ([]string, error)
	MintOwnify(email string, pubKey string, products []dto.BriefProduct, net string) ([]uint64, error)
	UpdatePinCode(role string, userId string, newPinCode string) error
	MakeTx(rawTx []byte, net string) (*string, *uint64, error)

	SendOwnify(email string, assetIds []uint64, sender string, receiver string, net string) (*string, error)

	DeleteOwnify(email string, assetIds []uint64, owner string, net string) (*string, error)
}

type walletService struct {
	wallet repository.AlgoHandler
}

func NewWalletService(wallet repository.AlgoHandler) WalletService {
	return &walletService{wallet: wallet}
}

func (w *walletService) AddNewAccount(
	role string,
	userId string,
	email string,
) (*string, error) {
	return w.wallet.NewWalletQuery().AddNewAccount(role, userId, email)
}

func (w *walletService) GetMyAccounts(
	role string,
	userId string,
) ([]string, error) {
	return w.wallet.NewWalletQuery().GetMyAccounts(role, userId)
}

func (w *walletService) MintOwnify(email string, pubKey string, products []dto.BriefProduct, net string) ([]uint64, error) {
	mintedIds, err := w.wallet.NewWalletQuery().MintOwnify(email, pubKey, products, net)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(mintedIds); i++ {
		products[i].AssetId = int64(mintedIds[i])
		products[i].Owner = pubKey
	}
	err = w.wallet.NewProductQuery().AddProducts(products, net, false)
	if err != nil {
		return nil, err
	}
	return mintedIds, nil
}

func (w *walletService) UpdatePinCode(role string, userId string, newPinCode string) error {
	return w.wallet.NewWalletQuery().UpdatePinCode(
		role, userId, newPinCode,
	)
}

func (w *walletService) MakeTx(rawTx []byte, net string) (*string, *uint64, error) {
	return w.wallet.NewWalletQuery().MakeTx(rawTx, net)
}

func (w *walletService) RegisterNewAccount(walletAddress string, userId string) (*string, error) {
	return w.wallet.NewWalletQuery().RegisterNewAccount(walletAddress, userId)
}

func (w *walletService) SendOwnify(email string, assetIds []uint64, sender string, receiver string, net string) (*string, error) {
	return w.wallet.NewWalletQuery().SendOwnify(email, assetIds, sender, receiver, net)
}

func (w *walletService) DeleteOwnify(email string, assetIds []uint64, owner string, net string) (*string, error) {
	return w.wallet.NewWalletQuery().DeleteOwnify(email, assetIds, owner, net)
}

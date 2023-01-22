package service

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type WalletService interface {
	AddNewAccount(role string, userId string) (*string, error)
	GetMyAccounts(role string, userId string, net string) ([]string, error)
	MintOwnify(pubKey string, products []dto.BriefProduct, net string) ([]string, error)
	UpdatePinCode(role string, userId string, newPinCode string) error
	MakeTransaction(role string, userId string, pubKey string, rawTx []byte, net string) (*string, error)
}

type walletService struct {
	wallet repository.AlgoHandler
}

func NewWalletService(wallet repository.AlgoHandler) WalletService {
	return &walletService{wallet}
}

func (w *walletService) AddNewAccount(
	role string,
	userId string,
) (*string, error) {
	return w.wallet.NewWalletQuery().AddNewAccount(role, userId)
}

func (w *walletService) GetMyAccounts(
	role string,
	userId string,
	net string,
) ([]string, error) {
	return w.wallet.NewWalletQuery().GetMyAccounts(role, userId, net)
}

func (w *walletService) MintOwnify(pubKey string, products []dto.BriefProduct, net string) ([]string, error) {
	return w.wallet.NewWalletQuery().MintOwnify(pubKey, products, net)
}

func (w *walletService) UpdatePinCode(role string, userId string, newPinCode string) error {
	return w.wallet.NewWalletQuery().UpdatePinCode(
		role, userId, newPinCode,
	)
}

func (w *walletService) MakeTransaction(role string, userId string, pubKey string, rawTx []byte, net string) (*string, error) {
	return w.wallet.NewWalletQuery().MakeTransaction(
		role, userId, pubKey, rawTx, net,
	)
}

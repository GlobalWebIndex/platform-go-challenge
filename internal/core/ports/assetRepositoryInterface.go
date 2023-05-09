package ports

import "github.com/loukaspe/platform-go-challenge/internal/core/domain"

type AssetRepositoryInterface interface {
	GetAsset(uuid string) (domain.Asset, error)
	CreateAsset(domain.Asset) (string, error)
	UpdateAsset(string, domain.Asset) error
	DeleteAsset(string) error
}

package service

import (
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type ProductService interface {
	AddProduct(
		product dto.BriefProduct,
		net string,
	) error
	GetProduct(
		chainId string, assetId string,
		net string,
	) (*domain.Product, error)
}

type productService struct {
	dbHandler repository.DBHandler
}

func NewProductService(dbHandler repository.DBHandler) ProductService {
	return &productService{dbHandler: dbHandler}
}

// GetUser implements Product Service
func (p *productService) AddProduct(
	product dto.BriefProduct,
	net string,
) error {
	err := p.dbHandler.NewProductQuery().AddProduct(product, net)
	return err
}

func (p *productService) GetProduct(
	chainId string, assetId string,
	net string,
) (*domain.Product, error) {
	product, err := p.dbHandler.NewProductQuery().GetProduct(chainId, assetId, net)
	if err != nil {
		return &domain.Product{}, nil
	}
	return &product, nil
}

package service

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type ProductService interface {
	AddProduct(
		product dto.BriefProduct,
		net string,
	) error
	GetProduct(
		chainId int, assetId int64,
		net string,
	) (*dto.BriefProduct, error)
	GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error)
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
	chainId int, assetId int64,
	net string,
) (*dto.BriefProduct, error) {
	product, err := p.dbHandler.NewProductQuery().GetProduct(chainId, assetId, net)
	if err != nil {
		return nil, err
	}
	return product, nil
}
func (p *productService) GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error) {
	return p.dbHandler.NewProductQuery().GetProducts(
		net, page, per_page,
	)
}

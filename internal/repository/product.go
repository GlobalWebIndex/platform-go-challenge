package repository

import (
	//"fmt"

	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"

	"github.com/Masterminds/squirrel"
	//"ownify_api/internal/dto"
	//"github.com/Masterminds/squirrel"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
)

type ProductQuery interface {
	AddProduct(
		product dto.BriefProduct,
		net string,
	) error
	AddProducts(products []dto.BriefProduct, net string) error
	GetProduct(chainId string, assetId string, net string) (domain.Product, error)
}

type productQuery struct{}

func (u *productQuery) AddProduct(product dto.BriefProduct, net string) error {
	// product validation.
	if !product.Valid() {
		return fmt.Errorf("[ERR] invalid Info: %v", product)
	}

	// add to database.
	tableName := domain.MainProductTableName
	if net == domain.TestNet {
		tableName = domain.TestProductTableName
	}
	cols, values := utils.ConvertToEntity(&product)
	sqlBuilder := utils.NewSqlBuilder()
	query, err := sqlBuilder.Insert(tableName, cols, values)
	if err != nil {
		return err
	}
	_, err = DB.Exec(*query)
	if err != nil {
		return err
	}
	return nil
}

func (u *productQuery) AddProducts(products []dto.BriefProduct, net string) error {
	// product validation
	for _, product := range products {
		if !product.Valid() {
			return fmt.Errorf("[ERR] invalid Info: %v", product)
		}
	}

	// add to database.
	tableName := domain.MainProductTableName
	if net == domain.TestNet {
		tableName = domain.TestProductTableName
	}
	sqlBuilder := utils.NewSqlBuilder()
	for _, product := range products {

		cols, values := utils.ConvertToEntity(&product)
		query, err := sqlBuilder.Insert(tableName, cols, values)
		if err != nil {
			return err
		}
		_, err = DB.Exec(*query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *productQuery) GetProduct(chainId string, assetId string, net string) (domain.Product, error) {
	tableName := domain.MainProductTableName
	if net == domain.TestNet {
		tableName = domain.TestProductTableName
	}
	var result domain.Product
	err := pgQb().
		Select("*").
		From(tableName).
		Where(squirrel.Eq{"chain_id": chainId, "asset_id": assetId}).Limit(1).Scan(&result)
	if err != nil {
		return domain.Product{}, err
	}
	return result, nil
}

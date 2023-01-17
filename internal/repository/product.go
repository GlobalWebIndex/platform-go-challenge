package repository

import (
	//"fmt"

	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"reflect"

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
	GetProduct(chainId string, assetId string, net string) (domain.Product, error)
}

type productQuery struct {}

func (u *productQuery) AddProduct(product dto.BriefProduct, net string) error {

	tableName := domain.MainProductTableName
	if net == domain.TestNet {
		tableName = domain.TestProductTableName
	}
	entity := reflect.ValueOf(product).Elem()
	cols := []string{}
	values := []interface{}{}
	for i := 0; i < entity.NumField(); i++ {
		field := entity.Type().Field(i).Name
		value := entity.FieldByName(field)
		cols = append(cols, field)
		values = append(values, value)
	}

	sqlBuilder := NewSqlBuilder()
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

package repository

import (
	//"fmt"

	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	"strings"
	"time"
)

type ProductQuery interface {
	AddProduct(
		product dto.BriefProduct,
		net string,
	) error
	AddProducts(products []dto.BriefProduct, net string) error
	GetProduct(chainId int, assetId int64, net string) (*dto.BriefProduct, error)

	GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error)

	SearchProducts(filter dto.BriefProduct, net string, page int32, perPage int32) ([]dto.BriefProduct, error)
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

func (u *productQuery) GetProduct(chainId int, assetId int64, net string) (*dto.BriefProduct, error) {
	tableName := domain.MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = domain.TestProductTableName
	}

	sqlBuilder := utils.NewSqlBuilder()

	var issue_date time.Time

	var product dto.BriefProduct
	sql, err := sqlBuilder.Select(tableName, []string{
		"owner",
		"barcode",
		"item_name",
		"brand_name",
		"additional_data",
		"issue_date",
		"location",
	}, []utils.Tuple{{
		Key: "asset_id",
		Val: assetId,
	}, {Key: "chain_id", Val: chainId}}, "=", "AND")

	if err != nil {
		return nil, err
	}
	err = DB.QueryRow(*sql).Scan(
		&product.Owner,
		&product.Barcode,
		&product.ItemName,
		&product.BrandName,
		&product.AdditionalData,
		&issue_date,
		&product.Location,
	)

	if err != nil {
		return nil, err
	}
	product.ChainId = chainId
	product.AssetId = assetId
	product.IssueDate = int32(issue_date.UnixMilli())

	return &product, nil
}

func (u *productQuery) GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error) {
	tableName := domain.MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = domain.TestProductTableName
	}

	sqlBuilder := utils.NewSqlBuilder()

	var issue_date time.Time

	var products []dto.BriefProduct
	preSql, err := sqlBuilder.Select(tableName, []string{
		"chain_id",
		"asset_id",
		"owner",
		"barcode",
		"item_name",
		"brand_name",
		"additional_data",
		"issue_date",
		"location",
	}, []utils.Tuple{}, "=", "AND")

	if err != nil {
		return nil, err
	}
	//SELECT * FROM products_test ORDER BY create_time  LIMIT 4 OFFSET 1;
	sql := *preSql + fmt.Sprintf(" ORDER BY create_time LIMIT %d OFFSET %d", per_page, page)
	rows, err := DB.Query(sql)

	for rows.Next() {
		var product dto.BriefProduct
		err := rows.Scan(
			&product.ChainId,
			&product.AssetId,
			&product.Owner,
			&product.Barcode,
			&product.ItemName,
			&product.BrandName,
			&product.AdditionalData,
			&issue_date,
			&product.Location,
		)
		if err != nil {
			continue
		}
		product.IssueDate = int32(issue_date.UnixMilli())
		products = append(products, product)
	}
	return products, nil
}

func (u *productQuery) SearchProducts(filter dto.BriefProduct, net string, page int32, perPage int32) ([]dto.BriefProduct, error) {
	tableName := domain.MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = domain.TestProductTableName
	}

	sqlBuilder := utils.NewSqlBuilder()

	var issue_date time.Time
	var products []dto.BriefProduct
	cons, values := utils.ConvertToEntity(&filter)
	conds := utils.GenerateCond(cons, values)
	preSql, err := sqlBuilder.Select(tableName, []string{
		"chain_id",
		"asset_id",
		"owner",
		"barcode",
		"item_name",
		"brand_name",
		"additional_data",
		"issue_date",
		"location",
	}, conds, "LIKE", "OR")

	if err != nil {
		return nil, err
	}
	//SELECT * FROM products_test ORDER BY create_time  LIMIT 4 OFFSET 1;
	sql := *preSql + fmt.Sprintf(" ORDER BY create_time LIMIT %d OFFSET %d", perPage, page)
	rows, err := DB.Query(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product dto.BriefProduct
		err := rows.Scan(
			&product.ChainId,
			&product.AssetId,
			&product.Owner,
			&product.Barcode,
			&product.ItemName,
			&product.BrandName,
			&product.AdditionalData,
			&issue_date,
			&product.Location,
		)
		if err != nil {
			continue
		}
		product.IssueDate = int32(issue_date.UnixMilli())
		products = append(products, product)
	}
	return products, nil
}

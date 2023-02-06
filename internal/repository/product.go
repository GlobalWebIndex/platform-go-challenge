package repository

import (
	//"fmt"
	"context"
	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	"strings"
)

type ProductQuery interface {
	AddProduct(
		product dto.BriefProduct,
		net string,
		verify bool,
	) error

	AddProducts(products []dto.BriefProduct, net string, verify bool) error
	GetProduct(chainId int, assetId int64, net string) (*dto.BriefProduct, error)

	GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error)

	SearchProducts(filter dto.BriefProduct, net string, page int32, perPage int32) ([]dto.BriefProduct, error)
}

type productQuery struct{}

func (u *productQuery) AddProduct(product dto.BriefProduct, net string, verify bool) error {
	// product validation.
	if verify {
		err := validProducts([]dto.BriefProduct{product}, net)
		if err != nil {
			return err
		}
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

func (u *productQuery) AddProducts(products []dto.BriefProduct, net string, verify bool) error {

	if verify {
		err := validProducts(products, net)
		if err != nil {
			return err
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
	var product dto.BriefProduct
	sql, err := sqlBuilder.Select(tableName, []string{
		"owner",
		"barcode",
		"item_name",
		"brand_name",
		"additional_data",
		"issued_date",
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
		&product.IssuedDate,
		&product.Location,
	)

	if err != nil {
		return nil, err
	}
	product.ChainId = chainId
	product.AssetId = assetId

	return &product, nil
}

func (u *productQuery) GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error) {
	tableName := domain.MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = domain.TestProductTableName
	}

	sqlBuilder := utils.NewSqlBuilder()

	var products []dto.BriefProduct
	preSql, err := sqlBuilder.Select(tableName, []string{
		"chain_id",
		"asset_id",
		"owner",
		"barcode",
		"item_name",
		"brand_name",
		"additional_data",
		"issued_date",
		"location",
	}, []utils.Tuple{}, "=", "AND")

	if err != nil {
		return nil, err
	}
	//SELECT * FROM products_test ORDER BY create_time  LIMIT 4 OFFSET 1;
	sql := *preSql + fmt.Sprintf(" ORDER BY created_time LIMIT %d OFFSET %d", per_page, page)
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
			&product.IssuedDate,
			&product.Location,
		)
		if err != nil {
			continue
		}
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
		"issued_date",
		"location",
	}, conds, "LIKE", "OR")

	if err != nil {
		return nil, err
	}
	//SELECT * FROM products_test ORDER BY create_time  LIMIT 4 OFFSET 1;
	sql := *preSql + fmt.Sprintf(" ORDER BY created_time LIMIT %d OFFSET %d", perPage, page)
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
			&product.IssuedDate,
			&product.Location,
		)
		if err != nil {
			continue
		}

		products = append(products, product)
	}
	return products, nil
}

func validProducts(products []dto.BriefProduct, net string) error {
	client, _, _ := NewClient(net)
	counter := len(products)
	// product validation
	validChan := make(chan domain.Result[int])
	for index, product := range products {
		go func(
			product dto.BriefProduct,
			index int,
		) {
			act, err := client.AccountInformation(product.Owner).Do(context.Background())
			if !product.Valid() || err != nil {
				validChan <- domain.Result[int]{
					Err: fmt.Errorf("invalid information at %d", index),
				}
			}
			isInclude := false
			for _, assetholding := range act.Assets {
				if uint64(product.AssetId) == assetholding.AssetId {
					isInclude = true
					break
				}
			}
			if !isInclude {
				validChan <- domain.Result[int]{
					Err: fmt.Errorf("invalid information at %d", index),
				}
			}
			counter -= 1
			if counter == 0 {
				validChan <- domain.Result[int]{
					Val: counter,
				}
			}
		}(product, index)
	}
	valid := <-validChan
	if valid.Err != nil {
		return valid.Err
	}
	return nil
}

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
	DeleteProducts(assetIds []uint64, net string) error
	GetProduct(chainId int, assetId int64, net string) (*dto.BriefProduct, error)

	GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error)

	SearchProducts(filter dto.BriefProduct, net string, page int32, perPage int32) (*int64, []dto.BriefProduct, error)
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
	tableName := MainProductTableName
	if net == domain.TestNet {
		tableName = TestProductTableName
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

	if len(products) == 0 {
		return nil
	}

	// add to database.
	tableName := MainProductTableName
	if net == domain.TestNet {
		tableName = TestProductTableName
	}

	var valueStrings []string
	var valueArgs []interface{}

	cols, _ := utils.ConvertToEntity(&products[0])

	for _, product := range products {
		_, values := utils.ConvertToEntity(&product)
		var qmarks []string

		for range cols {
			qmarks = append(qmarks, "?")
		}

		valueStrings = append(valueStrings, "("+strings.Join(qmarks, ", ")+")")
		valueArgs = append(valueArgs, values...)
	}

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, strings.Join(cols, ","), strings.Join(valueStrings, ","))

	_, err := DB.Exec(stmt, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (u *productQuery) DeleteProducts(assetIds []uint64, net string) error {
	// If no products to delete, just return.
	if len(assetIds) == 0 {
		return nil
	}

	// Choose the table name based on the network.
	tableName := MainProductTableName
	if net == domain.TestNet {
		tableName = TestProductTableName
	}

	// Prepare a string of placeholders for the query.
	placeholders := make([]string, len(assetIds))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	stmt := fmt.Sprintf("DELETE FROM %s WHERE  asset_id IN (%s)", tableName, strings.Join(placeholders, ","))

	// Convert the productIds to []interface{} for the Exec function.
	args := make([]interface{}, len(assetIds))
	for i, asset_id := range assetIds {
		args[i] = asset_id
	}

	_, err := DB.Exec(stmt, args...)
	if err != nil {
		return err
	}
	return nil
}

func (u *productQuery) GetProduct(chainId int, assetId int64, net string) (*dto.BriefProduct, error) {
	tableName := MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = TestProductTableName
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
	}}, "=", "AND")

	if err != nil {
		return nil, err
	}
	DB.QueryRow(*sql).Scan(
		&product.Owner,
		&product.Barcode,
		&product.ItemName,
		&product.BrandName,
		&product.AdditionalData,
		&product.IssuedDate,
		&product.Location,
	)

	product.ChainId = chainId
	product.AssetId = assetId

	return &product, nil
}

func (u *productQuery) GetProducts(net string, page int, per_page int) ([]dto.BriefProduct, error) {
	tableName := MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = TestProductTableName
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
		_ = rows.Scan(
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
		// if err != nil {
		// 	continue
		// }
		products = append(products, product)
	}
	return products, nil
}

func (u *productQuery) SearchProducts(filter dto.BriefProduct, net string, page int32, perPage int32) (*int64, []dto.BriefProduct, error) {

	tableName := MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = TestProductTableName
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
		return nil, nil, err
	}

	countsql, err := sqlBuilder.TotalCount(tableName, conds, "LIKE", "OR")
	if err != nil {
		return nil, nil, err
	}

	countChan := make(chan domain.Result[int64])
	productChan := make(chan domain.Result[[]dto.BriefProduct])

	defer close(countChan)
	defer close(productChan)
	go func() {
		var count int64
		err := DB.QueryRow(*countsql).Scan(&count)
		if err != nil {
			countChan <- domain.Result[int64]{Err: err}
			return
		}
		countChan <- domain.Result[int64]{Val: count}
	}()

	go func() {
		//SELECT * FROM products_test ORDER BY create_time  LIMIT 4 OFFSET 1;
		sql := *preSql + fmt.Sprintf(" ORDER BY created_time LIMIT %d OFFSET %d", perPage, page*perPage)

		rows, err := DB.Query(sql)
		if err != nil {
			productChan <- domain.Result[[]dto.BriefProduct]{Err: err}
			return
		}
		for rows.Next() {
			var product dto.BriefProduct
			_ = rows.Scan(
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
			products = append(products, product)
		}
		productChan <- domain.Result[[]dto.BriefProduct]{Val: products}
	}()

	count := <-countChan
	filteredProducts := <-productChan

	if count.Err != nil {
		return nil, nil, err
	}
	if filteredProducts.Err != nil {
		return nil, nil, err
	}

	return &count.Val, filteredProducts.Val, nil
}

func (u *productQuery) SearchProductsByAssetId(filter dto.BriefProduct, net string, page int32, perPage int32) (*int64, []dto.BriefProduct, error) {
	tableName := MainProductTableName
	if net == strings.ToLower(domain.TestNet) {
		tableName = TestProductTableName
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
		return nil, nil, err
	}

	countsql, err := sqlBuilder.TotalCount(tableName, conds, "LIKE", "OR")
	if err != nil {
		return nil, nil, err
	}

	countChan := make(chan domain.Result[int64])
	productChan := make(chan domain.Result[[]dto.BriefProduct])

	defer close(countChan)
	defer close(productChan)
	go func() {
		var count int64
		err := DB.QueryRow(*countsql).Scan(&count)
		if err != nil {
			countChan <- domain.Result[int64]{Err: err}
			return
		}
		countChan <- domain.Result[int64]{Val: count}
	}()

	go func() {
		//SELECT * FROM products_test ORDER BY create_time  LIMIT 4 OFFSET 1;
		sql := *preSql + fmt.Sprintf(" ORDER BY created_time LIMIT %d OFFSET %d", perPage, page*perPage)

		rows, err := DB.Query(sql)
		if err != nil {
			productChan <- domain.Result[[]dto.BriefProduct]{Err: err}
			return
		}
		for rows.Next() {
			var product dto.BriefProduct
			_ = rows.Scan(
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
			products = append(products, product)
		}
		productChan <- domain.Result[[]dto.BriefProduct]{Val: products}
	}()

	count := <-countChan
	filteredProducts := <-productChan

	if count.Err != nil {
		return nil, nil, err
	}
	if filteredProducts.Err != nil {
		return nil, nil, err
	}

	return &count.Val, filteredProducts.Val, nil
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

package repository_test

import (
	"os"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

}
func TestAddProduct(t *testing.T) {
	db, err := DBSetup()
	require.Nil(t, err)
	dbHandler := repository.NewDBHandler(db)
	product := dto.BriefProduct{
		ChainId:        1,
		AssetId:        338333,
		Owner:          "7OKB2T72KCO2E7GZO35XWEFW7IMU4PURMHRF2MAUFXLGIZG3ZSWC3VWJ6Q",
		Barcode:        "234324",
		ItemName:       "Tobacco",
		BrandName:      "Hello",
		AdditionalData: "S'tudio service",
		Location:       "This is test",
		IssuedDate:     3392939,
	}

	require.Equal(t, true, product.Valid())
	err = dbHandler.NewProductQuery().AddProduct(product, domain.TestNet, false)

	require.Nil(t, err)
}

func TestAddDuplicatedProduct(t *testing.T) {
	db, err := DBSetup()
	require.Nil(t, err)
	dbHandler := repository.NewDBHandler(db)

	product := dto.BriefProduct{
		ChainId:        1,
		AssetId:        338333,
		Owner:          "7OKB2T72KCO2E7GZO35XWEFW7IMU4PURMHRF2MAUFXLGIZG3ZSWC3VWJ6Q",
		Barcode:        "234324",
		ItemName:       "Tobacco",
		BrandName:      "Hello",
		AdditionalData: "This is test",
		Location:       "This is test",
		IssuedDate:     3392939,
	}

	//add first product
	require.Equal(t, true, product.Valid())
	err = dbHandler.NewProductQuery().AddProduct(product, domain.TestNet, false)
	require.Nil(t, err)

	//add second product
	err = dbHandler.NewProductQuery().AddProduct(product, domain.TestNet, false)
	require.NotNil(t, err)
}

func TestAddProducts(t *testing.T) {
	db, err := DBSetup()
	require.Nil(t, err)
	dbHandler := repository.NewDBHandler(db)

	products := []dto.BriefProduct{
		{
			ChainId:        1,
			AssetId:        338333,
			Owner:          "7OKB2T72KCO2E7GZO35XWEFW7IMU4PURMHRF2MAUFXLGIZG3ZSWC3VWJ6Q",
			Barcode:        "234324",
			ItemName:       "Tobacco",
			BrandName:      "Hello",
			AdditionalData: "S'tudio service",
			Location:       "This is test",
			IssuedDate:     3392939,
		},
		{
			ChainId:        1,
			AssetId:        3383333,
			Owner:          "7OKB2T72KCO2E7GZO35XWEFW7IMU4PURMHRF2MAUFXLGIZG3ZSWC3VWJ6Q",
			Barcode:        "234324",
			ItemName:       "Tobacco",
			BrandName:      "Hello",
			AdditionalData: "S'tudio service",
			Location:       "This is test",
			IssuedDate:     3392939,
		},
	}

	err = dbHandler.NewProductQuery().AddProducts(products, domain.TestNet, false)

	require.Nil(t, err)
}

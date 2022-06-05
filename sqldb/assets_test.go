package sqldb

import (
	"context"
	"log"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupSuite(tb testing.TB) (*DB, func(tb testing.TB)) {
	log.Println("setup suite")
	db, err := NewDB("user", "user", "127.0.0.1:3306", "mydb")
	if err != nil {
		tb.Fatal(err)
	}
	db.dropTablesIfExist()
	sqldb, _ := db.db.DB()
	db.createTables()
	// Return a function to teardown the test
	return db, func(tb testing.TB) {
		db.dropTablesIfExist()
		sqldb.Close()
	}
}

func TestCRUDInsight(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	asset, err := db.AddAsset(ctx, domain.Asset{Type: domain.InsightAssetType,
		Data: &domain.Insight{
			Text:        "40% of millenials spend more than 3hours on social media daily",
			Description: "example",
		}})
	assert.NotNil(t, asset)
	assert.Nil(t, err)
	assert.Equal(t, 1, asset.ID)
	assert.Equal(t, "example", asset.Data.(*domain.Insight).Description)

	asset, err = db.UpdateAsset(ctx, domain.Asset{
		ID:   1,
		Type: domain.InsightAssetType,
		Data: &domain.Insight{
			Text:        "100% of millenials spend more than 3hours on social media daily",
			Description: "updated example",
		}})
	assert.NotNil(t, asset)
	assert.Nil(t, err)
	assert.Equal(t, 1, asset.ID)
	assert.Equal(t, "updated example", asset.Data.(*domain.Insight).Description)

	err = db.DeleteAsset(ctx, asset.ID)
	assert.Nil(t, err)

}

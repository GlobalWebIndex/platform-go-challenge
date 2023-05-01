package repository

import (
	"fmt"
	"platform-go-challenge/integrations"
	"platform-go-challenge/models"
)

func GetFavoriteByUserId(user_id string) ([]models.Asset[any], error) {
	db := integrations.DB

	qry := "SELECT a.asset_id, a.description, a.type, adata FROM tAsset a, tFavorite f where f.user_id = $1 and a.asset_id = f.asset_id"

	assets, err := db.Query(qry, user_id)

	if err != nil {
		// handle this error
		panic(err)
	}

	assetsFound := make([]models.Asset[any], 0)

	for assets.Next() {
		var asset models.Asset[any]
		err = assets.Scan(&asset.Id, &asset.Description, &asset.Type, &asset.Data)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(asset.Id, asset.Description, asset.Type, asset.Data)
		assetsFound = append(assetsFound, asset)
	}

	return assetsFound, nil
}

func AddFavoriteByUserId(user_id string, asset models.Asset[any]) error {
	db := integrations.DB

	qry := "INSERT INTO tAsset (description, type, data) VALUES ($1, $2, $3)"

	result, err := db.Exec(qry, asset.Description, asset.Type, asset.Data)

	if err != nil {
		return err
	}

	asset_id, err1 := result.LastInsertId()

	if err1 != nil {
		return err1
	}

	_, err2 := db.Exec("INSERT INTO tFavorite (user_id, asset_id) VALUES ($1, $2)", user_id, asset_id)

	return err2
}

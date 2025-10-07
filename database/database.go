package database

import (
	"encoding/json"
	asset_model "gwi/models/asset"
	favorite_model "gwi/models/favorite"
	user_model "gwi/models/user"
)

var AssetsDB = []asset_model.Asset{
	{ID: 1, UserID: 1, Description: json.RawMessage(`{
		"title": "Bitcoin price",
		"x-axe title": "date",
		"y-axe title": "Price",
		"data": {
			"price": 94500.2,
			"date": "2024"
		}
	}`)},
	{ID: 2, UserID: 1, Description: json.RawMessage(`{
		"gender": "Female",
		"birth country": "Greece",
		"hours spent": 3,
		"number of purchases": "8"
	}`)},
	{ID: 3, UserID: 2, Description: json.RawMessage(`{
		"outcome": "40% of millenials spend more than 3hours on socialmedia daily",
	}`)},
}

var UsersDB = []user_model.User{
	{ID: 1, Email: "test@test.com"},
	{ID: 2, Email: "gwi@gwi.com"},
}

var FavoritesDB = []favorite_model.Favorite{
	{ID: 1, UserID: 1, AssetID: 1},
	{ID: 2, UserID: 1, AssetID: 2},
}

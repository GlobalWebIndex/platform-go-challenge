package router

import (
	"gwiEng/internal/api/controllers"
	"net/http"
)

func addFavouriteAsset(w http.ResponseWriter, r *http.Request) {
	err := controllers.AddFavouriteAsset(w, r, AppLogger, MongoClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

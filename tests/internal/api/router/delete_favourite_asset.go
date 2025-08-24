package router

import (
	"gwiEng/internal/api/controllers"
	"net/http"
)

func delFavouriteAsset(w http.ResponseWriter, r *http.Request) {
	 err := controllers.DelFavouriteAsset(w, r, AppLogger, MongoClient)	
	 if err != nil {
	 	w.WriteHeader(http.StatusInternalServerError)
	 }
}

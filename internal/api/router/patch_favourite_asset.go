package router

import (
	"gwiEng/internal/api/controllers"
	"net/http"
)

func patchFavouriteAsset(w http.ResponseWriter, r *http.Request) {
	err:= controllers.PatchFavouriteAsset(w, r, AppLogger, MongoClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
} 

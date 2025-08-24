package router

import (
	"gwiEng/internal/api/controllers"
	"net/http"
)

func getUserFavourites(w http.ResponseWriter, r *http.Request) {
	 err := controllers.GetUserFavourites(w, r, AppLogger, MongoClient)	
	 if err != nil {
	 	w.WriteHeader(http.StatusInternalServerError)
	 }
}

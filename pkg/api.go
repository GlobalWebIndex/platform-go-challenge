package gwi

import (
	"context"
	"log"
	"net/http"
)

type Handler struct {
	Repo Repository
}

// AddNewFav get input info (userid) and retrieves the fav list associated to the user
func (h *Handler) GetFavsFromUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("GetFavsFromUser -> Retrieving results...")

		userid := GetUuid(req)
		if !ValidateUuid(userid) {
			responseUnprocessable(w, "uuid not valid")
			return
		}
		assets, err := h.Repo.RetrieveFavs(context.Background(), userid)
		if err != nil {
			responseNoContent(w, "no results")
			return
		}
		responseList(w, true, assets)
	})
}

// AddNewFav get input info (userid and payload), validate it and add the asset to the user fav list
func (h *Handler) AddNewFav() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("AddNewFav -> Adding new fav to list...")

		userid := GetUuid(req)
		if !ValidateUuid(userid) {
			responseUnprocessable(w, "uuid not valid")
			return
		}

		asset := Asset{}
		output, err := Decode(req, asset)
		if err != nil {
			responseUnprocessable(w, "malformed input data")
			return
		}
		asset = output.(Asset)

		ok, err := h.Repo.AddAssetToFavs(context.Background(), userid, &asset)
		if err != nil {
			responseError(w, err.Error())
			return
		}
		if !ok {
			responseUnprocessable(w, "duplicated entry")
			return
		}
		responseOk(w, ok)
	})
}

// UpdateFav get input info (userid and payload), validate it and update the asset (if exists) into the user fav list
func (h *Handler) UpdateFav() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("UpdateFav -> Updating fav...")

		userid := GetUuid(req)
		if !ValidateUuid(userid) {
			responseUnprocessable(w, "uuid not valid")
			return
		}

		asset := Asset{}
		output, err := Decode(req, asset)
		if err != nil {
			responseUnprocessable(w, "malformed input data")
			return
		}
		asset = output.(Asset)
		ok, err := h.Repo.UpdateFav(context.Background(), userid, &asset)
		if err != nil {
			responseNoContent(w, err.Error())
			return
		}
		responseOk(w, ok)
	})
}

// DeleteFav get input info (userid and payload), validate it and delete the asset (if exists) from the user fav list
func (h *Handler) DeleteFav() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("DeleteFav -> Deleting fav...")

		userid := GetUuid(req)
		if !ValidateUuid(userid) {
			responseUnprocessable(w, "uuid not valid")
			return
		}

		asset := Asset{}
		output, err := Decode(req, asset)
		if err != nil {
			responseUnprocessable(w, "malformed input data")
			return
		}
		asset = output.(Asset)
		ok, err := h.Repo.RemoveFav(context.Background(), userid, &asset)
		if err != nil {
			responseNoContent(w, err.Error())
			return
		}
		responseOk(w, ok)
	})
}

package httpapi

import (
	"log"
	"net/http"

	application "github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/assets/domain"
	favapp "github.com/manos/favourites/favourites/application"
	favdomain "github.com/manos/favourites/favourites/domain"
	"github.com/manos/favourites/http/auth"
)

type AssetsHandler struct {
	svc *application.AssetService
}

func NewAssetsHandler(svc *application.AssetService) *AssetsHandler {
	return &AssetsHandler{svc: svc}
}

func (h *AssetsHandler) GetV1AssetsTypeId(w http.ResponseWriter, r *http.Request, pType GetV1AssetsTypeIdParamsType, assetID string) {
	userId, _ := auth.SubjectFromContext(r.Context())

	log.Printf("assets: %s %s user=%s type=%s id=%s", r.Method, r.URL.Path, userId, pType, assetID)

	var (
		result any
		err    error
	)

	switch pType {
	case Insights:
		result, err = h.svc.GetInsight(r.Context(), userId, assetID)
	case Audiences:
		result, err = h.svc.GetAudience(r.Context(), userId, assetID)
	case Charts:
		result, err = h.svc.GetChart(r.Context(), userId, assetID)
	default:
		log.Printf("assets: user=%s type=%s id=%s unknown type", userId, pType, assetID)
		http.NotFound(w, r)
		return
	}

	if err != nil {
		log.Printf("assets: user=%s type=%s id=%s err=%v", userId, pType, assetID, err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch v := result.(type) {
	case domain.InsightAsset:
		dto := favapp.AssetDTO{
			ID:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			OwnerUserID: v.UserId,
			Type:        favdomain.FavouriteInsight,
			Text:        v.Text,
		}
		writeJSON(w, http.StatusOK, assetToAPI(dto))
	case domain.AudienceAsset:
		dto := favapp.AssetDTO{
			ID:                v.Id,
			Name:              v.Name,
			Description:       v.Description,
			OwnerUserID:       v.UserId,
			Type:              favdomain.FavouriteAudience,
			SampleSize:        v.SampleSize,
			TotalRespondents:  v.TotalRespondents,
			EstimatedReach:    v.EstimatedReach,
			PopulationPercent: v.PopulationPercent,
		}
		writeJSON(w, http.StatusOK, assetToAPI(dto))
	case domain.ChartAsset:
		dto := favapp.AssetDTO{
			ID:          v.Id,
			Name:        v.Name,
			Description: v.Description,
			OwnerUserID: v.UserId,
			Type:        favdomain.FavouriteChart,
			XAxisTitle:  v.XAxisTitle,
			YAxisTitle:  v.YAxisTitle,
			Data:        v.Data,
		}
		writeJSON(w, http.StatusOK, assetToAPI(dto))
	default:
		log.Printf("assets: user=%s type=%s id=%s unknown asset type", userId, pType, assetID)
		http.NotFound(w, r)
		return
	}

	log.Printf("assets: user=%s type=%s id=%s served", userId, pType, assetID)
}

package services

import (
	"github.com/dintzeler/platform-go-challenge/models"
	"github.com/dintzeler/platform-go-challenge/storage"
	"github.com/dintzeler/platform-go-challenge/customerrors"
	"net/http"
	"os"
)

type FavoritesResponse struct {
	Charts	 []models.Chart `json:"charts"`
	Insights []models.Insight `json:"insights"`
	Audiences []models.Audience `json:"audiences"`
}

func getUserFavorites(userID int) []models.Favorite {
	data, err := storage.LoadData(os.Getenv("DATA_FILE"))
	if err != nil {
		return []models.Favorite{}
	}

	return storage.GetUserFavorites(data, userID)
}

func buildFavoritesResponse(favorites []models.Favorite, data *storage.DataStore) FavoritesResponse {
	favoritesResponse := FavoritesResponse{
		Charts:   []models.Chart{},
		Insights: []models.Insight{},
		Audiences: []models.Audience{},
	}
    for _, fav := range favorites {
        switch fav.AssetType {
        case models.AssetTypeChart:
			chart := storage.GetChartById(data.Charts, fav.AssetID)
            if chart != nil {
                favoritesResponse.Charts = append(favoritesResponse.Charts, *chart)
            }
        case models.AssetTypeInsight:
			insight := storage.GetInsightById(data.Insights, fav.AssetID)
            if insight != nil {
                favoritesResponse.Insights = append(favoritesResponse.Insights, *insight)
            }
        case models.AssetTypeAudience:
			audience := storage.GetAudienceById(data.Audiences, fav.AssetID)
            if audience != nil {
                favoritesResponse.Audiences = append(favoritesResponse.Audiences, *audience)
            }
        }
    }
	return favoritesResponse
}

func updateAssetDescription(data *storage.DataStore, userID int, assetID int, assetType models.AssetType, description string) {
	assetExists := storage.AssetExists(data, assetID, assetType)
	if !assetExists {
		return
	}

	userFavorites := storage.GetUserFavorites(data, userID)
	favoriteExists := storage.FavoriteExists(userFavorites, assetID, assetType)
	if !favoriteExists {
		return
	}

	switch assetType {
	case models.AssetTypeChart:
		for _, chart := range data.Charts {
			if chart.GetID() == assetID {
				chart.SetDescription(description)
			}
		}
	case models.AssetTypeInsight:
		for _, insight := range data.Insights {
			if insight.GetID() == assetID {
				insight.SetDescription(description)
			}
		}
	case models.AssetTypeAudience:
		for _, audience := range data.Audiences {
			if audience.GetID() == assetID {
				audience.SetDescription(description)
			}
		}
	}

	storage.SaveData(os.Getenv("DATA_FILE"), data)


}


func GetFavorites(userID int) (*FavoritesResponse, error) {
	data, err := storage.LoadData(os.Getenv("DATA_FILE"))
	if err != nil {
		return nil, &customerrors.ValidationError{
			Message:   "Error loading data",
			ErrorCode: "DATA_LOAD_ERROR",
		}
	}
	
	favorites := storage.GetUserFavorites(data, userID)
	favoritesResponse := buildFavoritesResponse(favorites, data)
	
	return &favoritesResponse, nil
}

func AddFavorite(userID int, assetID int, assetType models.AssetType) (int, error) {
	data, err := storage.LoadData(os.Getenv("DATA_FILE"))
	if err != nil {
		return http.StatusInternalServerError, &customerrors.ValidationError{
			Message:   "Error loading data",
			ErrorCode: "DATA_LOAD_ERROR",
		}
	}

	assetExists := storage.AssetExists(data, assetID, assetType)
	if !assetExists {
		return http.StatusNotFound, &customerrors.ValidationError{
			Message:   "Asset does not exist",
			ErrorCode: "ASSET_NOT_FOUND",
		}
	}

	userFavorites := storage.GetUserFavorites(data, userID)
	favoriteExists := storage.FavoriteExists(userFavorites, assetID, assetType)
	if favoriteExists {
		return http.StatusOK, nil
	}
	


	err = storage.CreateFavorite(data, userID, assetID, assetType)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func DeleteFavorite(userID int, assetID int, assetType models.AssetType) (int, error) {
	data, err := storage.LoadData(os.Getenv("DATA_FILE"))
	if err != nil {
		return http.StatusInternalServerError, &customerrors.ValidationError{
			Message:   "Error loading data",
			ErrorCode: "DATA_LOAD_ERROR",
		}
	}

	assetExists := storage.AssetExists(data, assetID, assetType)
	if !assetExists {
		return http.StatusNotFound, &customerrors.ValidationError{
			Message:   "Asset does not exist",
			ErrorCode: "ASSET_NOT_FOUND",
		}
	}
	userFavorites := storage.GetUserFavorites(data, userID)
	favoriteExists := storage.FavoriteExists(userFavorites, assetID, assetType)
	if !favoriteExists {
		return http.StatusNotFound, &customerrors.ValidationError{
			Message:   "Favorite does not exist",
			ErrorCode: "FAVORITE_NOT_FOUND",
		}
	}
	err = storage.RemoveFavorite(data, userID, assetID, assetType)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	
	return http.StatusNoContent, nil
}

func UpdateFavorite(userID int, assetType models.AssetType, assetID int, description string) (int, error) {
	data, err := storage.LoadData(os.Getenv("DATA_FILE"))
	if err != nil {
		return http.StatusInternalServerError, &customerrors.ValidationError{
			Message:   "Error loading data",
			ErrorCode: "DATA_LOAD_ERROR",
		}
	}

	assetExists := storage.AssetExists(data, assetID, assetType)
	if !assetExists {
		return http.StatusNotFound, &customerrors.ValidationError{
			Message:   "Asset does not exist",
			ErrorCode: "ASSET_NOT_FOUND",
		}
	}

	userFavorites := storage.GetUserFavorites(data, userID)
	favoriteExists := storage.FavoriteExists(userFavorites, assetID, assetType)
	if !favoriteExists {
		return http.StatusForbidden, &customerrors.ValidationError{
			Message:   "User does not have this asset as favorite",
			ErrorCode: "ASSET_NOT_FAVORITE",
		}
	}

	err = storage.UpdateAssetDescription(data, userID, assetID, assetType, description)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
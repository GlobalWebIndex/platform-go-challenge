package services

import (
	"app/api/schema"
	"app/services"
	
	"fmt"
	"net/http"
	"log"
)

type XService interface {
	GetAssets(filter schema.AssetFilter) []schema.AssetDetailsDto
	getAsset(id string) (schema.AssetDetailsDto, error)
	UpdateAsset(assetID string, updateRequest schema.AssetUpdateRequest) error
	GetUserFavorites(userID string) []schema.UserFavoriteDto
	CreateUserFavorite(userID string, assetID string) (string, error)
	DeleteUserFavorite(userID string, userFavoriteID string) error
}

type ApiService struct {
	Svc services.GService
}

func NewXService(svc services.GService) XService {
	return &ApiService{Svc: svc}
}

func (a *ApiService) GetAssets(filter schema.AssetFilter) []schema.AssetDetailsDto {
	return a.Svc.FindAssets(filter)
}

func (a *ApiService) getAsset(id string) (schema.AssetDetailsDto, error) {
	dto, err := a.Svc.FindAsset(id)
	if err != nil {
		return schema.AssetDetailsDto{}, schema.NewApiError(http.StatusNotFound, err)
	}
	return dto, nil
}

func (a *ApiService) UpdateAsset(assetID string, updateRequest schema.AssetUpdateRequest) error {
	return a.Svc.UpdateAsset(assetID, updateRequest.Description)
}

func (a *ApiService) GetUserFavorites(userID string) []schema.UserFavoriteDto {
	ufs := a.Svc.FindUserFavorites(userID)
	if len(ufs) == 0 {
		return make([]schema.UserFavoriteDto, 0)
	}
	userFavoriteDtos := make([]schema.UserFavoriteDto, 0, len(ufs))
	for _, uf := range ufs {
		assetDto, err := a.getAsset(uf.AssetID)
		if err != nil {
			log.Printf("[!] failed to get asset details for asset %s, skipping...\n", uf.AssetID)
			continue
		} else {
		dto := schema.UserFavoriteDto{ID: uf.ID.Hex(), Details: assetDto}
		userFavoriteDtos = append(userFavoriteDtos, dto)
		}
	}
	return userFavoriteDtos
}

func (a *ApiService) CreateUserFavorite(userID string, assetID string) (string, error) {
	_, err := a.Svc.FindUser(userID)
	if err != nil {
		fmt.Printf("[X] user %s not found\n", userID)
		return "", schema.NewApiError(http.StatusNotFound, err)
	}
	if _, err := a.Svc.FindAsset(assetID); err != nil {
		fmt.Printf("[X] asset %s does not exist\n", assetID)
		return "", schema.NewApiError(http.StatusNotFound, err)
	}
	id, err := a.Svc.CreateFavoriteAsset(assetID, userID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (a *ApiService) DeleteUserFavorite(userID string, userFavoriteID string) error {
	_, err := a.Svc.FindUser(userID)
	if err != nil {
		fmt.Printf("user %s not found", userID)
		return schema.NewApiError(http.StatusNotFound, err)
	}
	return a.Svc.DeleteUserFavorite(userFavoriteID)
}

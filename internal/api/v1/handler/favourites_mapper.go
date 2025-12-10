package handler

import (
	"errors"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/api/v1/dto"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

// parseAssetType converts a DTO asset type string into the domain AssetType.
func parseAssetType(s string) (favourites.AssetType, error) {
	switch dto.AssetType(s) {
	case dto.AssetTypeChart:
		return favourites.AssetTypeChart, nil
	case dto.AssetTypeInsight:
		return favourites.AssetTypeInsight, nil
	case dto.AssetTypeAudience:
		return favourites.AssetTypeAudience, nil
	default:
		return "", errors.New("unknown asset type")
	}
}

// dtoToDomainAsset converts a CreateFavouriteRequest into a domain Asset.
func dtoToDomainAsset(req dto.CreateFavouriteRequest) (favourites.Asset, error) {
	var asset favourites.Asset

	at, err := parseAssetType(string(req.Asset.Type))
	if err != nil {
		return favourites.Asset{}, err
	}
	asset.Type = at
	asset.Description = req.Asset.Description

	switch req.Asset.Type {
	case dto.AssetTypeChart:
		if req.Asset.Chart == nil {
			return favourites.Asset{}, favourites.ErrInvalidAssetPayload
		}
		asset.Chart = &favourites.Chart{
			Title: req.Asset.Chart.Title,
			XAxis: req.Asset.Chart.XAxis,
			YAxis: req.Asset.Chart.YAxis,
			Data:  req.Asset.Chart.Data,
		}
	case dto.AssetTypeInsight:
		if req.Asset.Insight == nil {
			return favourites.Asset{}, favourites.ErrInvalidAssetPayload
		}
		asset.Insight = &favourites.Insight{
			Text: req.Asset.Insight.Text,
		}
	case dto.AssetTypeAudience:
		if req.Asset.Audience == nil {
			return favourites.Asset{}, favourites.ErrInvalidAssetPayload
		}
		asset.Audience = &favourites.Audience{
			Gender:                   req.Asset.Audience.Gender,
			BirthCountry:             req.Asset.Audience.BirthCountry,
			AgeGroup:                 req.Asset.Audience.AgeGroup,
			HoursOnSocialMediaPerDay: req.Asset.Audience.HoursOnSocialMediaPerDay,
			PurchasesLastMonth:       req.Asset.Audience.PurchasesLastMonth,
		}
	default:
		return favourites.Asset{}, favourites.ErrInvalidAssetPayload
	}

	return asset, nil
}

// favouriteEntryToCreateResponse converts a domain FavouriteEntry into a CreateFavouriteResponse.
func favouriteEntryToCreateResponse(entry favourites.FavouriteEntry) dto.CreateFavouriteResponse {
	return dto.CreateFavouriteResponse{
		UserID:    entry.Favourite.UserID,
		AssetID:   entry.Favourite.AssetID,
		CreatedAt: entry.Favourite.CreatedAt,
	}
}

// favouriteEntryToDTO converts a domain FavouriteEntry into an API DTO.
func favouriteEntryToDTO(entry favourites.FavouriteEntry) dto.FavouriteResponse {
	assetDTO := assetToDTO(entry.Asset)
	return dto.FavouriteResponse{
		Asset:     assetDTO,
		CreatedAt: entry.Favourite.CreatedAt,
	}
}

// favouriteEntryToListItemDTO converts a domain FavouriteEntry into a list item DTO.
func favouriteEntryToListItemDTO(entry favourites.FavouriteEntry) dto.FavouriteItemResponse {
	assetDTO := assetToDTO(entry.Asset)
	return dto.FavouriteItemResponse{
		UserID:    entry.Favourite.UserID,
		AssetID:   entry.Favourite.AssetID,
		CreatedAt: entry.Favourite.CreatedAt,
		Asset:     assetDTO,
	}
}

// assetToDTO converts a domain Asset into an API DTO.
func assetToDTO(asset favourites.Asset) dto.AssetResponse {
	var chart *dto.ChartPayload
	if asset.Chart != nil {
		chart = &dto.ChartPayload{
			Title: asset.Chart.Title,
			XAxis: asset.Chart.XAxis,
			YAxis: asset.Chart.YAxis,
			Data:  asset.Chart.Data,
		}
	}

	var insight *dto.InsightPayload
	if asset.Insight != nil {
		insight = &dto.InsightPayload{
			Text: asset.Insight.Text,
		}
	}

	var audience *dto.AudiencePayload
	if asset.Audience != nil {
		audience = &dto.AudiencePayload{
			Gender:                   asset.Audience.Gender,
			BirthCountry:             asset.Audience.BirthCountry,
			AgeGroup:                 asset.Audience.AgeGroup,
			HoursOnSocialMediaPerDay: asset.Audience.HoursOnSocialMediaPerDay,
			PurchasesLastMonth:       asset.Audience.PurchasesLastMonth,
		}
	}

	return dto.AssetResponse{
		ID:          asset.ID,
		Type:        dto.AssetType(asset.Type),
		Description: asset.Description,
		Chart:       chart,
		Insight:     insight,
		Audience:    audience,
	}
}

package assets_client

import (
	"context"
	"errors"

	assetsapp "github.com/manos/favourites/assets/application"
	"github.com/manos/favourites/favourites/application"
	favdomain "github.com/manos/favourites/favourites/domain"
)

type Client struct {
	assetService assetsapp.AssetsService
}

func NewAssetClient(assetService assetsapp.AssetsService) *Client {
	return &Client{assetService: assetService}
}

func (c *Client) GetInsight(ctx context.Context, userId string, assetId string) (application.AssetDTO, error) {
	a, err := c.assetService.GetInsight(ctx, userId, assetId)
	if err != nil {
		if errors.Is(err, assetsapp.ErrAssetNotFound) {
			return application.AssetDTO{}, application.ErrAssetNotFound
		}
		return application.AssetDTO{}, err
	}
	return application.AssetDTO{
		ID:          a.Id,
		Name:        a.Name,
		Description: a.Description,
		OwnerUserID: a.UserId,
		Type:        favdomain.FavouriteInsight,
		Text:        a.Text,
	}, nil
}

func (c *Client) GetAudience(ctx context.Context, userId string, assetId string) (application.AssetDTO, error) {
	a, err := c.assetService.GetAudience(ctx, userId, assetId)
	if err != nil {
		if errors.Is(err, assetsapp.ErrAssetNotFound) {
			return application.AssetDTO{}, application.ErrAssetNotFound
		}
		return application.AssetDTO{}, err
	}
	return application.AssetDTO{
		ID:                a.Id,
		Name:              a.Name,
		Description:       a.Description,
		OwnerUserID:       a.UserId,
		Type:              favdomain.FavouriteAudience,
		SampleSize:        a.SampleSize,
		TotalRespondents:  a.TotalRespondents,
		EstimatedReach:    a.EstimatedReach,
		PopulationPercent: a.PopulationPercent,
	}, nil
}

func (c *Client) GetChart(ctx context.Context, userId string, assetId string) (application.AssetDTO, error) {
	a, err := c.assetService.GetChart(ctx, userId, assetId)
	if err != nil {
		if errors.Is(err, assetsapp.ErrAssetNotFound) {
			return application.AssetDTO{}, application.ErrAssetNotFound
		}
		return application.AssetDTO{}, err
	}
	return application.AssetDTO{
		ID:          a.Id,
		Name:        a.Name,
		Description: a.Description,
		OwnerUserID: a.UserId,
		Type:        favdomain.FavouriteChart,
		XAxisTitle:  a.XAxisTitle,
		YAxisTitle:  a.YAxisTitle,
		Data:        a.Data,
	}, nil
}

func (c *Client) GetInsights(ctx context.Context, userId string, assetIds []string) (map[string]application.AssetDTO, error) {
	assets, err := c.assetService.GetInsights(ctx, userId, assetIds)
	if err != nil {
		if errors.Is(err, assetsapp.ErrAssetNotFound) {
			return nil, application.ErrAssetNotFound
		}
		return nil, err
	}
	out := make(map[string]application.AssetDTO, len(assets))
	for id, a := range assets {
		out[id] = application.AssetDTO{
			ID:          a.Id,
			Name:        a.Name,
			Description: a.Description,
			OwnerUserID: a.UserId,
			Type:        favdomain.FavouriteInsight,
			Text:        a.Text,
		}
	}
	return out, nil
}

func (c *Client) GetAudiences(ctx context.Context, userId string, assetIds []string) (map[string]application.AssetDTO, error) {
	assets, err := c.assetService.GetAudiences(ctx, userId, assetIds)
	if err != nil {
		if errors.Is(err, assetsapp.ErrAssetNotFound) {
			return nil, application.ErrAssetNotFound
		}
		return nil, err
	}
	out := make(map[string]application.AssetDTO, len(assets))
	for id, a := range assets {
		out[id] = application.AssetDTO{
			ID:                a.Id,
			Name:              a.Name,
			Description:       a.Description,
			OwnerUserID:       a.UserId,
			Type:              favdomain.FavouriteAudience,
			SampleSize:        a.SampleSize,
			TotalRespondents:  a.TotalRespondents,
			EstimatedReach:    a.EstimatedReach,
			PopulationPercent: a.PopulationPercent,
		}
	}
	return out, nil
}

func (c *Client) GetCharts(ctx context.Context, userId string, assetIds []string) (map[string]application.AssetDTO, error) {
	assets, err := c.assetService.GetCharts(ctx, userId, assetIds)
	if err != nil {
		if errors.Is(err, assetsapp.ErrAssetNotFound) {
			return nil, application.ErrAssetNotFound
		}
		return nil, err
	}
	out := make(map[string]application.AssetDTO, len(assets))
	for id, a := range assets {
		out[id] = application.AssetDTO{
			ID:          a.Id,
			Name:        a.Name,
			Description: a.Description,
			OwnerUserID: a.UserId,
			Type:        favdomain.FavouriteChart,
			XAxisTitle:  a.XAxisTitle,
			YAxisTitle:  a.YAxisTitle,
			Data:        a.Data,
		}
	}
	return out, nil
}

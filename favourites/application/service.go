package application

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/manos/favourites/favourites/domain"
)

var (
	ErrFavouriteNotFound      = errors.New("favourite not found")
	ErrFavouriteForbidden     = errors.New("favourite does not belong to user")
	ErrAssetNotFound          = errors.New("asset not found")
	ErrFavouriteAlreadyExists = errors.New("favourite already exists")
)

type FavouriteService struct {
	repo   FavouriteRepository
	assets AssetClient
}

func NewFavouriteService(repo FavouriteRepository, assets AssetClient) *FavouriteService {
	return &FavouriteService{repo: repo, assets: assets}
}

func (s *FavouriteService) GetFavouritesForUser(ctx context.Context, userId string, page int, limit int) (FavouritePageDTO, error) {
	log.Printf("svc favourites: list user=%s page=%d limit=%d", userId, page, limit)
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	favs, total, err := s.repo.FindByUserPage(ctx, userId, offset, limit)
	if err != nil {
		log.Printf("svc favourites: list user=%s err=%v", userId, err)
		return FavouritePageDTO{}, err
	}

	// collect ids per type to batch asset fetches
	insightIDs := make([]string, 0)
	audienceIDs := make([]string, 0)
	chartIDs := make([]string, 0)
	for _, fav := range favs {
		switch fav.Type {
		case domain.FavouriteInsight:
			insightIDs = append(insightIDs, fav.AssetID)
		case domain.FavouriteAudience:
			audienceIDs = append(audienceIDs, fav.AssetID)
		case domain.FavouriteChart:
			chartIDs = append(chartIDs, fav.AssetID)
		default:
			return FavouritePageDTO{}, fmt.Errorf("unknown favourite type: %s", fav.Type)
		}
	}

	insights := make(map[string]AssetDTO)
	audiences := make(map[string]AssetDTO)
	charts := make(map[string]AssetDTO)

	if len(insightIDs) > 0 {
		insights, err = s.assets.GetInsights(ctx, userId, insightIDs)
		if err != nil {
			return FavouritePageDTO{}, err
		}
	}
	if len(audienceIDs) > 0 {
		audiences, err = s.assets.GetAudiences(ctx, userId, audienceIDs)
		if err != nil {
			return FavouritePageDTO{}, err
		}
	}
	if len(chartIDs) > 0 {
		charts, err = s.assets.GetCharts(ctx, userId, chartIDs)
		if err != nil {
			return FavouritePageDTO{}, err
		}
	}

	out := make([]FavouriteDTO, 0, len(favs))

	for _, fav := range favs {
		var asset AssetDTO
		var ok bool
		switch fav.Type {
		case domain.FavouriteInsight:
			asset, ok = insights[fav.AssetID]
		case domain.FavouriteAudience:
			asset, ok = audiences[fav.AssetID]
		case domain.FavouriteChart:
			asset, ok = charts[fav.AssetID]
		}

		if !ok {
			log.Printf("svc favourites: list user=%s favourite_id=%s asset=%s missing; skipping", userId, fav.ID, fav.AssetID)
			continue
		}

		out = append(out, FavouriteDTO{
			ID:          fav.ID,
			UserID:      fav.UserID,
			Type:        fav.Type,
			Description: fav.Description,
			Asset:       asset,
		})
	}

	// Recompute totals based on available assets only.
	total = len(out)
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return FavouritePageDTO{
		Items:      out,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *FavouriteService) AddFavourite(ctx context.Context, userId string, favType string, assetId string, description string) (string, error) {

	// var asset AssetDTO
	var err error

	log.Printf("svc favourites: add user=%s type=%s asset=%s", userId, favType, assetId)
	favouriteType := domain.FavouriteType(favType)
	switch favouriteType {
	case domain.FavouriteInsight:
		_, err = s.assets.GetInsight(ctx, userId, assetId)
	case domain.FavouriteAudience:
		_, err = s.assets.GetAudience(ctx, userId, assetId)
	case domain.FavouriteChart:
		_, err = s.assets.GetChart(ctx, userId, assetId)
	default:
		return "", fmt.Errorf("unknown favourite type: %s", favouriteType)
	}

	if err != nil {
		log.Printf("svc favourites: add user=%s type=%s asset=%s validation err=%v", userId, favType, assetId, err)
		return "", err
	}

	fav, err := domain.NewFavourite(userId, domain.FavouriteType(favType), assetId, description)
	if err != nil {
		log.Printf("svc favourites: add user=%s type=%s asset=%s build err=%v", userId, favType, assetId, err)
		return "", err
	}
	err = s.repo.Save(ctx, fav)
	if err != nil {
		log.Printf("svc favourites: add user=%s type=%s asset=%s save err=%v", userId, favType, assetId, err)
		return "", err
	}
	log.Printf("svc favourites: add user=%s favourite_id=%s created", userId, fav.ID)
	return fav.ID, nil
}

func (s *FavouriteService) DeleteFavourite(ctx context.Context, userId string, favouriteID string) error {
	log.Printf("svc favourites: delete user=%s favourite_id=%s", userId, favouriteID)
	fav, err := s.repo.FindByID(ctx, favouriteID)
	if err != nil {
		log.Printf("svc favourites: delete user=%s favourite_id=%s find err=%v", userId, favouriteID, err)
		return err
	}

	if fav.UserID != userId {
		log.Printf("svc favourites: delete user=%s favourite_id=%s forbidden owner=%s", userId, favouriteID, fav.UserID)
		return ErrFavouriteForbidden
	}

	return s.repo.Delete(ctx, favouriteID)
}

func (s *FavouriteService) UpdateFavouriteDescription(ctx context.Context, userId string, favouriteID string, description string) error {
	log.Printf("svc favourites: update description user=%s favourite_id=%s", userId, favouriteID)
	fav, err := s.repo.FindByID(ctx, favouriteID)
	if err != nil {
		log.Printf("svc favourites: update description user=%s favourite_id=%s find err=%v", userId, favouriteID, err)
		return err
	}

	if fav.UserID != userId {
		log.Printf("svc favourites: update description user=%s favourite_id=%s forbidden owner=%s", userId, favouriteID, fav.UserID)
		return ErrFavouriteForbidden
	}

	if err := s.repo.UpdateDescription(ctx, favouriteID, description); err != nil {
		log.Printf("svc favourites: update description user=%s favourite_id=%s save err=%v", userId, favouriteID, err)
		return err
	}

	log.Printf("svc favourites: update description user=%s favourite_id=%s updated", userId, favouriteID)
	return nil
}

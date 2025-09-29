package services

import (
	"time"

	"github.com/google/uuid"
	"favourite_assets/server/models"
	"favourite_assets/server/repositories"
	"favourite_assets/server/errors"

)

type AssetService struct {
	repo *repositories.AssetRepository
}

func NewAssetService(repo *repositories.AssetRepository) *AssetService {
	return &AssetService{
		repo: repo,
	}
}

func (s *AssetService) CreateAsset(asset models.Asset) (models.Asset, error) {
	if asset.GetID() == uuid.Nil {
		switch a := asset.(type) {
		case *models.Chart:
			a.ID = uuid.New()
			a.CreatedAt = time.Now()
			a.UpdatedAt = time.Now()
		case *models.Insight:
			a.ID = uuid.New()
			a.CreatedAt = time.Now()
			a.UpdatedAt = time.Now()
		case *models.Audience:
			a.ID = uuid.New()
			a.CreatedAt = time.Now()
			a.UpdatedAt = time.Now()
		default:
			return nil, errors.ErrBadRequest
		}
	}

	if err := s.repo.Create(asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (s *AssetService) GetAsset(id uuid.UUID) (models.Asset, error) {
	return s.repo.GetByID(id)
}

func (s *AssetService) UpdateAsset(assetID uuid.UUID, updatedData map[string]interface{}) (models.Asset, error) {
    existing, err := s.repo.GetByID(assetID)
    if err != nil {
        return nil, errors.ErrNotFound
    }

    existing.SetDescription(updatedData["description"].(string))

    switch a := existing.(type) {
    case *models.Chart:
        a.Title = updatedData["title"].(string)
        a.XAxis = updatedData["xAxis"].(string)
        a.YAxis = updatedData["yAxis"].(string)
		a.UpdatedAt = time.Now() 
    case *models.Insight:
        a.Text = updatedData["text"].(string)
		a.UpdatedAt = time.Now() 
    case *models.Audience:
        a.Gender = updatedData["gender"].(string)
        a.BirthCountry = updatedData["birthCountry"].(string)
        a.AgeGroup = updatedData["ageGroup"].(string)
        a.HoursOnSocial = int(updatedData["hoursSocialDaily"].(float64))
        a.PurchasesLastMonth = int(updatedData["purchasesLastMonth"].(float64))
		a.UpdatedAt = time.Now() 
    }

    if err := s.repo.Update(existing); err != nil {
        return nil, err
    }

    return existing, nil
}

func (s *AssetService) DeleteAsset(id uuid.UUID) error {
	if err := s.repo.Delete(id); err != nil {
		return errors.ErrNotFound
	}
	return nil
}

func (s *AssetService) ListAssets() []models.Asset {
	return s.repo.ListAll()
}

func (s *AssetService) ListAssetsByType(assetType models.AssetType) []models.Asset {
	var result []models.Asset
	for _, a := range s.repo.ListAll() {
		if a.GetType() == assetType { 
			result = append(result, a)
		}
	}
	return result
}


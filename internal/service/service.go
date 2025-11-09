package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/KostasDasios/platform-go-challenge/internal/models"
	"github.com/KostasDasios/platform-go-challenge/internal/repo"
)

type Service struct {
	repo repo.Repository
}

// NewService constructs a Service using the provided Repository.
func NewService(repo repo.Repository) *Service { return &Service{repo: repo} }

var userIDRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,64}$`)

func (s *Service) ValidateUserID(id string) bool { return userIDRe.MatchString(id) }

// ListFavourites returns all favourites for a user after validating the identifier.
func (s *Service) ListFavourites(userID string) ([]*models.Favourite, error) {
	if !s.ValidateUserID(userID) {
		return nil, fmt.Errorf("invalid user id")
	}
	return s.repo.List(userID)
}

// CreateFavourite validates the raw asset payload, normalises metadata and persists a new favourite.
func (s *Service) CreateFavourite(userID string, raw json.RawMessage) (*models.Favourite, error) {
	if !s.ValidateUserID(userID) {
		return nil, fmt.Errorf("invalid user id")
	}
	t, desc, err := validateAsset(raw)
	if err != nil {
		return nil, err
	}
	f := &models.Favourite{
		ID:          newID(),
		Type:        t,
		Description: desc,
		Asset:       raw,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.Create(userID, f); err != nil {
		return nil, err
	}
	return f, nil
}

// UpdateFavouriteDescription updates only the editable description field for a favourite.
func (s *Service) UpdateFavouriteDescription(userID, favID, desc string) (*models.Favourite, error) {
	if !s.ValidateUserID(userID) || strings.TrimSpace(favID) == "" {
		return nil, fmt.Errorf("invalid path")
	}
	return s.repo.UpdateDescription(userID, favID, desc)
}

// DeleteFavourite removes a favourite by id.
func (s *Service) DeleteFavourite(userID, favID string) error {
	if !s.ValidateUserID(userID) || strings.TrimSpace(favID) == "" {
		return fmt.Errorf("invalid path")
	}
	return s.repo.Delete(userID, favID)
}

// validateAsset performs a two-step decode: probe for type, then validate concrete schema.
// This keeps the service flexible for additional asset types without changing the transport contract.
func validateAsset(raw json.RawMessage) (models.AssetType, string, error) {
	var probe struct {
		Type        models.AssetType `json:"type"`
		Description string    `json:"description"`
	}
	if err := json.Unmarshal(raw, &probe); err != nil {
		return "", "", fmt.Errorf("invalid asset json: %w", err)
	}
	switch probe.Type {
	case models.AssetChart:
		var c models.Chart
		if err := json.Unmarshal(raw, &c); err != nil {
			return "", "", fmt.Errorf("invalid chart: %w", err)
		}
		if strings.TrimSpace(c.Title) == "" || len(c.Data) == 0 {
			return "", "", errors.New("chart needs title and non-empty data")
		}
		return models.AssetChart, c.Description, nil
	case models.AssetInsight:
		var in models.Insight
		if err := json.Unmarshal(raw, &in); err != nil {
			return "", "", fmt.Errorf("invalid insight: %w", err)
		}
		if strings.TrimSpace(in.Text) == "" {
			return "", "", errors.New("insight needs text")
		}
		return models.AssetInsight, in.Description, nil
	case models.AssetAudience:
		var a models.Audience
		if err := json.Unmarshal(raw, &a); err != nil {
			return "", "", fmt.Errorf("invalid audience: %w", err)
		}
		if a.Gender == "" || len(a.AgeGroups) == 0 {
			return "", "", errors.New("audience needs gender and age_groups")
		}
		return models.AssetAudience, a.Description, nil
	default:
		return "", "", errors.New("unknown asset type")
	}
}

// newID generates a short random identifier.
// In production this would be replaced with ULID/UUIDv7 for sortability and uniqueness guarantees.
func newID() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

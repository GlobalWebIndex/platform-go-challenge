package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go-favorites-app/internal/core/domain/favorites"

	"github.com/google/uuid"
)

// Pagination helper
type Pagination struct {
	Limit  int
	Offset int
}

func NewPagination(r *http.Request) Pagination {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}
	// Cap limit if necessary (e.g. 100)
	if limit > 1000 {
		limit = 1000
	}

	offset := (page - 1) * limit
	return Pagination{Limit: limit, Offset: offset}
}

// createAssetRequest is a helper struct to handle polymorphic unmarshal
type createAssetRequest struct {
	Type favorites.AssetType `json:"type"`
	Raw  json.RawMessage
}

func (r *createAssetRequest) UnmarshalJSON(data []byte) error {
	type typeHeader struct {
		Type favorites.AssetType `json:"type"`
	}
	var header typeHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return err
	}
	r.Type = header.Type
	// We want to store the full raw data to unmarshal into the concrete type later
	// However, json.Unmarshal logic for RawMessage would just copy data.
	// But since we are inside UnmarshalJSON, we need to be careful not to recurse if we reused the type.
	// Here typeHeader only captures 'type', we'll just copy data to Raw.
	r.Raw = make(json.RawMessage, len(data))
	copy(r.Raw, data)
	return nil
}

func parseAsset(data []byte, assetType favorites.AssetType, userID string) (favorites.Asset, error) {
	generateID := func(currentID string) string {
		if currentID == "" {
			return uuid.NewString()
		}
		return currentID
	}

	switch assetType {
	case favorites.AssetTypeChart:
		var c favorites.Chart
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		c.UserID = userID
		c.ID = generateID(c.ID)
		return c, nil
	case favorites.AssetTypeInsight:
		var i favorites.Insight
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}
		i.UserID = userID
		i.ID = generateID(i.ID)
		return i, nil
	case favorites.AssetTypeAudience:
		var a favorites.Audience
		if err := json.Unmarshal(data, &a); err != nil {
			return nil, err
		}
		a.UserID = userID
		a.ID = generateID(a.ID)
		return a, nil
	default:
		return nil, errors.New("unknown asset type")
	}
}

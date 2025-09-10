package entadapter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/SokratisChaimanas/platform-go-challenge/ent"
	"github.com/SokratisChaimanas/platform-go-challenge/ent/favourite"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/ports"
	"github.com/google/uuid"
)

// Compile time safety for ports.FavouriteRepository implementation.
var _ ports.FavouriteRepository = (*FavouriteRepo)(nil)

// FavouriteRepo implements ports.FavouriteRepository using Ent.
type FavouriteRepo struct {
	client *ent.Client
}

func NewFavouriteRepo(client *ent.Client) *FavouriteRepo {
	return &FavouriteRepo{client: client}
}

// Create inserts a new favourite. Duplicate entries map to ErrFavouriteAlreadyExists.
func (favouriteRepo *FavouriteRepo) Create(ctx context.Context, favouriteToCreate *domain.Favourite) error {
	_, err := favouriteRepo.client.Favourite.
		Create().
		SetUserID(favouriteToCreate.UserID).
		SetAssetID(favouriteToCreate.AssetID).
		SetCreatedAt(favouriteToCreate.CreatedAt).
		Save(ctx)

	return err
}

// Delete removes a favourite by (userID, assetID). Missing rows map to ErrFavouriteNotFound.
func (favouriteRepo *FavouriteRepo) Delete(ctx context.Context, userID, assetID uuid.UUID) error {
	f, err := favouriteRepo.client.Favourite.
		Delete().
		Where(
			favourite.UserID(userID),
			favourite.AssetID(assetID),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	if f == 0 {
		return domain.ErrFavouriteNotFound
	}

	return nil
}

// Exists checks if a favourite already exists for (userID, assetID).
func (favouriteRepo *FavouriteRepo) Exists(ctx context.Context, userID, assetID uuid.UUID) (bool, error) {
	return favouriteRepo.client.Favourite.
		Query().
		Where(
			favourite.UserID(userID),
			favourite.AssetID(assetID),
		).
		Exist(ctx)
}

// ListAssetsFavouritedByUserKeyset returns assets favourited by userID using keyset pagination.
func (favouriteRepo *FavouriteRepo) ListAssetsFavouritedByUserKeyset(
	ctx context.Context, userID uuid.UUID, limit int, after string,
) ([]domain.Asset, *string, error) {

	// Bound limits the same way your parsePagination does (default 20, cap 50).
	if limit <= 0 {
		limit = 20
	} else if limit > 50 {
		limit = 50
	}

	q := favouriteRepo.client.Favourite.
		Query().
		Where(favourite.UserID(userID)).
		// Deterministic total order: created_at ASC, id ASC
		Order(
			favourite.ByCreatedAt(sql.OrderAsc()),
			favourite.ByID(sql.OrderAsc()),
		).
		WithAsset() // eager load assets since we return assets, not favourites

	// Seek to > (created_at,id) if a cursor was provided.
	if after != "" {
		cur, err := decodeCursor(after)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", domain.ErrBadCursor, err)
		}
		q = q.Where(
			favourite.Or(
				favourite.CreatedAtGT(cur.T),
				favourite.And(
					favourite.CreatedAtEQ(cur.T),
					favourite.IDGT(cur.I),
				),
			),
		)
	}

	// Pull one extra to know if there's another page.
	rows, err := q.Limit(limit + 1).All(ctx)
	if err != nil {
		return nil, nil, err
	}

	var nextAfter *string
	if len(rows) > limit {
		last := rows[limit-1]
		cstr, err := encodeCursor(ksCursor{T: last.CreatedAt, I: last.ID})
		if err != nil {
			return nil, nil, err
		}
		nextAfter = &cstr
		rows = rows[:limit]
	}

	assets := make([]domain.Asset, 0, len(rows))
	for _, f := range rows {
		a := f.Edges.Asset
		if a == nil {
			// Shouldn't happen because WithAsset(), but guard anyway.
			continue
		}
		assets = append(assets, domain.Asset{
			ID:          a.ID,
			Type:        domain.AssetType(a.AssetType),
			Description: a.Description,
			Payload:     a.Payload,
		})
	}

	return assets, nextAfter, nil
}

// ksCursor encodes the "position" of the last favourite (created_at, id).
type ksCursor struct {
	T time.Time `json:"t"` // favourite.created_at
	I uuid.UUID `json:"i"` // favourite.id
}

// encodeCursor turns a cursor into a URL-safe base64 string.
func encodeCursor(c ksCursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// decodeCursor parses a URL-safe base64 cursor string.
func decodeCursor(s string) (ksCursor, error) {
	var c ksCursor
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, nil
}

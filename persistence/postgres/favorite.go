package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"pgc/domain"
)

var errDatabaseConnectionNil = errors.New("database connection cannot be nil")

const (
	// psqlColumnDescription is the name of the column used for storing the description of a favorite asset.
	psqlColumnDescription = "description"
)

// FavoriteRepository is a repository for managing user favorite assets in the database.
type FavoriteRepository struct {
	db *bun.DB // Database connection.
}

// Favorite represents a user's favorite asset in the database.
type Favorite struct {
	bun.BaseModel `bun:"table:favorites,alias:f"`

	ID          int    `bun:"id,autoincrement,pk"`  // Unique identifier for the favorite asset.
	UserID      int    `bun:"user_id,notnull"`      // ID of the user who marked the asset as favorite.
	AssetID     int    `bun:"asset_id,notnull"`     // ID of the asset marked as favorite.
	Description string `bun:"description,nullzero"` // Optional description for the favorite asset.

	User  *User  `bun:"rel:belongs-to,join:user_id=id"`  // Relation to the User entity.
	Asset *Asset `bun:"rel:belongs-to,join:asset_id=id"` // Relation to the Asset entity.
}

// Insert adds a new favorite asset into the database for a specific user.
func (r *FavoriteRepository) Insert(ctx context.Context, favorite domain.Favorite) (domain.Favorite, error) {
	switch {
	case favorite.UserID <= 0:
		return domain.Favorite{}, domain.ErrInvalidUserID
	case favorite.AssetID <= 0:
		return domain.Favorite{}, domain.ErrInvalidAssetID
	}

	psqlFavorite := Favorite{
		UserID:      favorite.UserID,
		AssetID:     favorite.AssetID,
		Description: favorite.Description,
	}

	_, err := r.db.NewInsert().
		Model(&psqlFavorite).
		Exec(ctx)

	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Field('C') == "23503" && pgErr.Field('n') == "favorites_user_id_fkey":
				// Handle foreign key violation for user ID.
				return domain.Favorite{}, domain.ErrInvalidUserID
			case pgErr.Field('C') == "23503" && pgErr.Field('n') == "favorites_asset_id_fkey":
				// Handle foreign key violation for asset ID.
				return domain.Favorite{}, domain.ErrInvalidAssetID
			case pgErr.Field('C') == "23505":
				// Handle unique constraint violation.
				return domain.Favorite{}, domain.ErrDuplicateFavorite
			default:
				slog.Error("PostgreSQL error creating favorite", "error", err, "user_id", favorite.UserID, "asset_id", favorite.AssetID)
				return domain.Favorite{}, domain.ErrInternalServerError
			}
		}

		slog.Error("Failed to insert favorite", "error", err, "user_id", favorite.UserID, "asset_id", favorite.AssetID)
		return domain.Favorite{}, domain.ErrInternalServerError
	}

	return domain.Favorite{
		ID:          psqlFavorite.ID,
		UserID:      psqlFavorite.UserID,
		AssetID:     psqlFavorite.AssetID,
		Description: psqlFavorite.Description,
	}, nil
}

// Remove deletes a favorite asset from the database for a specific user.
func (r *FavoriteRepository) Remove(ctx context.Context, favoriteID, userID int) error {
	switch {
	case favoriteID <= 0:
		return domain.ErrInvalidFavoriteID
	case userID <= 0:
		return domain.ErrInvalidUserID
	}

	res, err := r.db.NewDelete().
		Model((*Favorite)(nil)).
		Where("id = ? AND user_id = ?", favoriteID, userID).
		Exec(ctx)

	if err != nil {
		slog.Error("Failed to delete favorite", "error", err, "favoriteID", favoriteID, "user_id", userID)
		return domain.ErrInternalServerError
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected", "error", err)
		return domain.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return domain.ErrFavoriteNotFound
	}

	return nil
}

// SearchByUserID retrieves all favorite assets for a specific user from the database.
func (r *FavoriteRepository) SearchByUserID(ctx context.Context, userID int) ([]domain.Favorite, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	var psqlFavorites []Favorite
	err := r.db.NewSelect().
		Model(&psqlFavorites).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		slog.Error("Failed to retrieve favorites for user", "error", err, "user_id", userID)
		return nil, domain.ErrInternalServerError
	}

	// Map database favorites to domain favorites.
	favorites := make([]domain.Favorite, len(psqlFavorites))
	for k, favorite := range psqlFavorites {
		favorites[k] = domain.Favorite{
			ID:          favorite.ID,
			UserID:      favorite.UserID,
			AssetID:     favorite.AssetID,
			Description: favorite.Description,
		}
	}

	return favorites, nil
}

// UpdateDescription updates the description of a favorite asset in the database for a specific user.
func (r *FavoriteRepository) UpdateDescription(ctx context.Context, favoriteID, userID int, description string) (domain.Favorite, error) {
	switch {
	case favoriteID <= 0:
		return domain.Favorite{}, domain.ErrInvalidFavoriteID
	case userID <= 0:
		return domain.Favorite{}, domain.ErrInvalidUserID
	}

	psqlFavorite := Favorite{
		ID:          favoriteID,
		UserID:      userID,
		Description: description,
	}

	err := r.db.NewUpdate().
		Model(&psqlFavorite).
		Column(psqlColumnDescription).
		WherePK().
		Where("user_id = ?", psqlFavorite.UserID).
		Returning("*").
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Favorite{}, domain.ErrFavoriteNotFound
		}

		slog.Error("Failed to update favorite description", "error", err, "favorite_id", favoriteID)
		return domain.Favorite{}, domain.ErrInternalServerError
	}

	return domain.Favorite{
		ID:          psqlFavorite.ID,
		UserID:      psqlFavorite.UserID,
		AssetID:     psqlFavorite.AssetID,
		Description: psqlFavorite.Description,
	}, nil
}

// NewFavoriteRepository creates a new instance of FavoriteRepository with the provided database connection.
func NewFavoriteRepository(db *bun.DB) (*FavoriteRepository, error) {
	if db == nil {
		slog.Error("Failed to create FavoriteRepository: database connection is nil")
		return nil, errDatabaseConnectionNil
	}

	return &FavoriteRepository{db: db}, nil
}

// Package main contains the GWI Favorites Service implementation in Go
//

// Architecture:

// - models: data structures for assets, users, favorites

// - storage: database operations (PostgreSQL)

// - service: business logic layer

// - handlers: HTTP endpoint handlers

//

// Key Go decisions:

// - goroutines: handle multiple requests concurrently (built-in)

// - channels: coordinate between handlers and storage

// - interfaces: make components testable and swappable

// - connection pooling: database/sql handles this automatically

// - prepared statements: prevent SQL injection

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	DBConnString     = "user=user password=password host=postgres port=5432 dbname=gwi_challenge sslmode=disable"
	DefaultPageSize  = 20
	MaxPageSize      = 100
	CacheTTLSeconds  = 300 // 5 minutes
	MaxConnections   = 25  // database/sql pools automatically
	RequestTimeout   = 30 * time.Second
)

// ValidAssetTypes defines which asset types are allowed
var ValidAssetTypes = map[string]bool{
	"chart":    true,
	"insight":  true,
	"audience": true,
}

// ============================================================================
// DATA MODELS
// ============================================================================

// Asset represents any asset in the system (Chart, Insight, or Audience).
// Using a single struct with type discrimination keeps code simpler.
type Asset struct {
	ID   string          `json:"id"`
	Type string          `json:"type"` // "chart", "insight", "audience"
	Data json.RawMessage `json:"data"` // Type-specific data as JSON
}

// Favorite represents an asset favorited by a user.
// The description_override lets users customize how the asset appears in their list.
type Favorite struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	Asset               *Asset     `json:"asset"`
	DescriptionOverride *string    `json:"description_override"`
	AddedAt             time.Time  `json:"added_at"`
	IsDeleted           bool       `json:"is_deleted"`
}

// PaginatedResponse wraps a list of favorites with pagination metadata.
type PaginatedResponse struct {
	Favorites  []*Favorite    `json:"favorites"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo contains metadata about pagination.
type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ErrorResponse formats errors for HTTP responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ============================================================================
// DATABASE LAYER
// ============================================================================

// Storage handles all database operations. Keeping storage separate from
// business logic makes the code testable and follows single responsibility.
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new Storage instance with database connection.
// The connection pool is created automatically by database/sql.
func NewStorage() (*Storage, error) {
	// database/sql automatically manages a connection pool
	// Default max connections is 0 (unlimited), but we limit it
	db, err := sql.Open("postgres", DBConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	// Max idle: keep this many connections ready
	// Max open: allow at most this many connections
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(MaxConnections)
	db.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return &Storage{db: db}, nil
}

// Close closes the database connection pool.
func (s *Storage) Close() error {
	return s.db.Close()
}

// CreateUser creates a user (idempotent). Users are minimal - just ID.
func (s *Storage) CreateUser(userID string) error {
	query := `
		INSERT INTO users (id)
		VALUES ($1)
		ON CONFLICT (id) DO NOTHING
	`
	_, err := s.db.Exec(query, userID)
	return err
}

// UserExists checks if a user exists. Used for validation.
func (s *Storage) UserExists(userID string) (bool, error) {
	query := "SELECT id FROM users WHERE id = $1"
	var id string
	err := s.db.QueryRow(query, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListUsers fetches all users with pagination.
// Returns (users, totalCount, error)
func (s *Storage) ListUsers(limit int, offset int) ([]*struct {
	ID        string
	CreatedAt time.Time
}, int, error) {
	// Get total count
	countQuery := "SELECT COUNT(*) FROM users"
	var total int
	err := s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch page
	query := `
		SELECT id, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*struct {
		ID        string
		CreatedAt time.Time
	}
	for rows.Next() {
		var id string
		var createdAt time.Time
		if err := rows.Scan(&id, &createdAt); err != nil {
			return nil, 0, err
		}
		users = append(users, &struct {
			ID        string
			CreatedAt time.Time
		}{ID: id, CreatedAt: createdAt})
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ============================================================================
// USER MANAGEMENT - DELETE USER
// ============================================================================

// DeleteUser deletes a user and all their associated data.
// Cascades to remove all their favorites.
// Returns true if user found and deleted, false if not found.
func (s *Storage) DeleteUser(userID string) (bool, error) {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := s.db.Exec(query, userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// ============================================================================
// ASSET MANAGEMENT - CREATE, READ, DELETE ASSETS
// ============================================================================

// CreateAsset creates a new asset and returns its ID.
// Data is stored as JSONB for flexibility and queryability.
func (s *Storage) CreateAsset(assetType string, data json.RawMessage) (string, error) {
	assetID := uuid.New().String()
	query := `
		INSERT INTO assets (id, type, data)
		VALUES ($1, $2, $3)
	`
	_, err := s.db.Exec(query, assetID, assetType, string(data))
	if err != nil {
		return "", err
	}
	return assetID, nil
}

// GetAsset fetches a single asset by ID. Returns nil if not found.
func (s *Storage) GetAsset(assetID string) (*Asset, error) {
	query := "SELECT id, type, data FROM assets WHERE id = $1"
	var id, assetType string
	var dataStr string
	err := s.db.QueryRow(query, assetID).Scan(&id, &assetType, &dataStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &Asset{
		ID:   id,
		Type: assetType,
		Data: json.RawMessage(dataStr),
	}, nil
}

// ListAssets fetches all assets with pagination.
// Returns (assets, totalCount, error)
func (s *Storage) ListAssets(limit int, offset int, assetType *string) ([]*Asset, int, error) {
	// Get total count
	whereClause := ""
	queryArgs := []interface{}{}
	if assetType != nil && ValidAssetTypes[*assetType] {
		whereClause = " WHERE type = $1"
		queryArgs = append(queryArgs, *assetType)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM assets%s", whereClause)
	var total int
	err := s.db.QueryRow(countQuery, queryArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch page
	queryArgs = append(queryArgs, limit, offset)
	argCount := len(queryArgs) - 1
	query := fmt.Sprintf(`
		SELECT id, type, data
		FROM assets%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	rows, err := s.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []*Asset
	for rows.Next() {
		var id, assetType string
		var dataStr string
		if err := rows.Scan(&id, &assetType, &dataStr); err != nil {
			return nil, 0, err
		}
		assets = append(assets, &Asset{
			ID:   id,
			Type: assetType,
			Data: json.RawMessage(dataStr),
		})
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

// AssetExists checks if an asset exists. Used for validation.
func (s *Storage) AssetExists(assetID string) (bool, error) {
	query := "SELECT id FROM assets WHERE id = $1"
	var id string
	err := s.db.QueryRow(query, assetID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteAsset deletes an asset by ID.
// Returns true if found and deleted, false if not found.
func (s *Storage) DeleteAsset(assetID string) (bool, error) {
	query := `
		DELETE FROM assets
		WHERE id = $1
	`
	result, err := s.db.Exec(query, assetID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// ============================================================================
// FAVORITE MANAGEMENT - CREATE, READ, UPDATE, DELETE FAVORITES
// ============================================================================

// AddToFavorites adds an asset to a user's favorites.
// Returns the favorite ID or "" if already favorited.
// This uses a prepared statement automatically (sql.Exec handles this).
func (s *Storage) AddToFavorites(
	userID string,
	assetID string,
	descriptionOverride *string,
) (string, error) {
	favoriteID := uuid.New().String()
	query := `
		INSERT INTO favorites (id, user_id, asset_id, description_override)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, asset_id) WHERE deleted_at IS NULL
		DO NOTHING
	`
	result, err := s.db.Exec(query, favoriteID, userID, assetID, descriptionOverride)
	if err != nil {
		// Check if it's a foreign key constraint violation
		return "", fmt.Errorf("failed to add favorite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	// If no rows affected, the favorite already exists
	if rowsAffected == 0 {
		return "", nil
	}

	return favoriteID, nil
}

// GetFavorites fetches paginated favorites for a user.
// Returns (favorites, totalCount, error)
//
// This query uses indexes efficiently:
// - (user_id, deleted_at, added_at) index speeds up filtering and sorting
// - deleted_at IS NULL filter is part of the index predicate
// - JOIN to assets table is fast because asset_id is indexed
func (s *Storage) GetFavorites(
	userID string,
	limit int,
	offset int,
	assetType *string,
) ([]*Favorite, int, error) {
	// Build query dynamically based on filters
	whereClause := "WHERE deleted_at IS NULL AND user_id = $1"
	queryArgs := []interface{}{userID}
	argCount := 2

	if assetType != nil && ValidAssetTypes[*assetType] {
		whereClause += fmt.Sprintf(" AND a.type = $%d", argCount)
		queryArgs = append(queryArgs, *assetType)
		argCount++
	}

	// First, get the total count (needed for pagination metadata)
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM favorites f
		JOIN assets a ON f.asset_id = a.id
		%s
	`, whereClause)
	var total int
	err := s.db.QueryRow(countQuery, queryArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Now fetch the actual page
	// ORDER BY f.added_at DESC: newest favorites first
	// LIMIT $n OFFSET $n: pagination
	queryArgs = append(queryArgs, limit, offset)
	query := fmt.Sprintf(`
		SELECT
			f.id,
			f.user_id,
			f.description_override,
			f.added_at,
			a.id,
			a.type,
			a.data
		FROM favorites f
		JOIN assets a ON f.asset_id = a.id
		%s
		ORDER BY f.added_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	rows, err := s.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var favorites []*Favorite
	for rows.Next() {
		var (
			favID, userID, assetID, assetType string
			descOverride                       *string
			addedAt                            time.Time
			dataStr                            string
		)

		err := rows.Scan(
			&favID,
			&userID,
			&descOverride,
			&addedAt,
			&assetID,
			&assetType,
			&dataStr,
		)
		if err != nil {
			return nil, 0, err
		}

		fav := &Favorite{
			ID:                  favID,
			UserID:              userID,
			DescriptionOverride: descOverride,
			AddedAt:             addedAt,
			Asset: &Asset{
				ID:   assetID,
				Type: assetType,
				Data: json.RawMessage(dataStr),
			},
		}

		favorites = append(favorites, fav)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}

// UpdateFavoriteDescription updates the description for a favorited asset.
// Returns true if found and updated, false if not found.
func (s *Storage) UpdateFavoriteDescription(
	userID string,
	assetID string,
	description string,
) (bool, error) {
	query := `
		UPDATE favorites
		SET description_override = $1
		WHERE user_id = $2 AND asset_id = $3 AND deleted_at IS NULL
	`
	result, err := s.db.Exec(query, description, userID, assetID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// RemoveFromFavorites soft-deletes a favorite (marks as deleted, doesn't remove).
// This preserves data for auditing and recovery.
// Returns true if found and deleted, false if not found.
func (s *Storage) RemoveFromFavorites(userID string, assetID string) (bool, error) {
	query := `
		UPDATE favorites
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND asset_id = $2 AND deleted_at IS NULL
	`
	result, err := s.db.Exec(query, userID, assetID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// ============================================================================
// SERVICE LAYER - Business Logic
// ============================================================================

// Service orchestrates operations between HTTP handlers and storage.
// This layer contains business logic and validation.
type Service struct {
	storage *Storage
}

// NewService creates a new service.
func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

// CreateUser creates a new user and returns the created user object.
func (s *Service) CreateUser() (map[string]interface{}, error) {
	userID := uuid.New().String()
	err := s.storage.CreateUser(userID)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return map[string]interface{}{
		"id":         userID,
		"created_at": time.Now().UTC(),
	}, nil
}

// ListUsers retrieves paginated user list.
func (s *Service) ListUsers(page int, limit int) (map[string]interface{}, error) {
	// Validate and constrain pagination
	if limit < 1 {
		limit = 1
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	// Fetch from storage
	users, total, err := s.storage.ListUsers(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}

	// Transform to API format
	userList := []map[string]interface{}{}
	for _, u := range users {
		userList = append(userList, map[string]interface{}{
			"id":         u.ID,
			"created_at": u.CreatedAt,
		})
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return map[string]interface{}{
		"users": userList,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	}, nil
}

// ============================================================================
// USER MANAGEMENT - DELETE USER SERVICE METHOD
// ============================================================================

// DeleteUser removes a user from the system.
func (s *Service) DeleteUser(userID string) error {
	// Check if user exists
	exists, err := s.storage.UserExists(userID)
	if err != nil {
		return fmt.Errorf("error checking user: %w", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Delete user
	_, err = s.storage.DeleteUser(userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

// ============================================================================
// ASSET MANAGEMENT - CREATE, LIST, DELETE ASSET SERVICE METHODS
// ============================================================================

// CreateAsset creates a new asset in the system.
func (s *Service) CreateAsset(assetType string, data json.RawMessage) (map[string]interface{}, error) {
	// Validate asset type
	if !ValidAssetTypes[assetType] {
		return nil, fmt.Errorf("invalid asset type")
	}

	// Create asset
	assetID, err := s.storage.CreateAsset(assetType, data)
	if err != nil {
		return nil, fmt.Errorf("error creating asset: %w", err)
	}

	return map[string]interface{}{
		"id":   assetID,
		"type": assetType,
		"data": json.RawMessage(data),
	}, nil
}

// ListAssets retrieves paginated asset list.
func (s *Service) ListAssets(page int, limit int, assetType *string) (map[string]interface{}, error) {
	// Validate and constrain pagination
	if limit < 1 {
		limit = 1
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	if page < 1 {
		page = 1
	}

	// Validate asset type if provided
	if assetType != nil && *assetType != "" && !ValidAssetTypes[*assetType] {
		return nil, fmt.Errorf("invalid asset type")
	}

	offset := (page - 1) * limit

	// Fetch from storage
	assets, total, err := s.storage.ListAssets(limit, offset, assetType)
	if err != nil {
		return nil, fmt.Errorf("error fetching assets: %w", err)
	}

	// Transform to API format
	assetList := []map[string]interface{}{}
	for _, a := range assets {
		assetList = append(assetList, map[string]interface{}{
			"id":   a.ID,
			"type": a.Type,
			"data": a.Data,
		})
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return map[string]interface{}{
		"assets": assetList,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	}, nil
}

// DeleteAsset removes an asset from the system.
func (s *Service) DeleteAsset(assetID string) error {
	// Check if asset exists
	exists, err := s.storage.AssetExists(assetID)
	if err != nil {
		return fmt.Errorf("error checking asset: %w", err)
	}
	if !exists {
		return fmt.Errorf("asset not found")
	}

	// Delete asset
	_, err = s.storage.DeleteAsset(assetID)
	if err != nil {
		return fmt.Errorf("error deleting asset: %w", err)
	}

	return nil
}

// ============================================================================
// FAVORITES SERVICE METHODS
// ============================================================================

// AddFavorite adds an asset to user's favorites with validation.
func (s *Service) AddFavorite(
	userID string,
	assetID string,
	description *string,
) (*Favorite, error) {
	// Validate user exists
	exists, err := s.storage.UserExists(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Validate asset exists
	asset, err := s.storage.GetAsset(assetID)
	if err != nil {
		return nil, fmt.Errorf("error getting asset: %w", err)
	}
	if asset == nil {
		return nil, fmt.Errorf("asset not found")
	}

	// Try to add to favorites
	favoriteID, err := s.storage.AddToFavorites(userID, assetID, description)
	if err != nil {
		return nil, fmt.Errorf("error adding favorite: %w", err)
	}
	if favoriteID == "" {
		// Empty ID means already favorited
		return nil, fmt.Errorf("asset already in favorites")
	}

	return &Favorite{
		ID:                  favoriteID,
		UserID:              userID,
		Asset:               asset,
		DescriptionOverride: description,
		AddedAt:             time.Now(),
	}, nil
}

// GetFavorites retrieves user's favorites with pagination.
func (s *Service) GetFavorites(
	userID string,
	page int,
	limit int,
	assetType *string,
) (*PaginatedResponse, error) {
	// Validate user exists
	exists, err := s.storage.UserExists(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Validate and constrain pagination
	if limit < 1 {
		limit = 1
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	// Fetch from storage
	favorites, total, err := s.storage.GetFavorites(userID, limit, offset, assetType)
	if err != nil {
		return nil, fmt.Errorf("error fetching favorites: %w", err)
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit // Ceiling division
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedResponse{
		Favorites: favorites,
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// UpdateFavoriteDescription updates a favorite's description.
func (s *Service) UpdateFavoriteDescription(
	userID string,
	assetID string,
	description string,
) (*Favorite, error) {
	// Validate user exists
	exists, err := s.storage.UserExists(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Get current favorite to return full object
	favorites, _, err := s.storage.GetFavorites(userID, 1000, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching favorites: %w", err)
	}

	var favorite *Favorite
	for _, f := range favorites {
		if f.Asset.ID == assetID {
			favorite = f
			break
		}
	}

	if favorite == nil {
		return nil, fmt.Errorf("asset not in user's favorites")
	}

	// Update description
	success, err := s.storage.UpdateFavoriteDescription(userID, assetID, description)
	if err != nil {
		return nil, fmt.Errorf("error updating favorite: %w", err)
	}
	if !success {
		return nil, fmt.Errorf("failed to update description")
	}

	// Update and return
	favorite.DescriptionOverride = &description
	return favorite, nil
}

// RemoveFavorite removes an asset from user's favorites.
func (s *Service) RemoveFavorite(userID string, assetID string) error {
	// Validate user exists
	exists, err := s.storage.UserExists(userID)
	if err != nil {
		return fmt.Errorf("error checking user: %w", err)
	}
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Remove favorite
	success, err := s.storage.RemoveFromFavorites(userID, assetID)
	if err != nil {
		return fmt.Errorf("error removing favorite: %w", err)
	}
	if !success {
		return fmt.Errorf("asset not in user's favorites")
	}

	return nil
}

// ============================================================================
// HTTP HANDLERS
// ============================================================================

// RequestHandler holds dependencies for all HTTP handlers.
type RequestHandler struct {
	service *Service
}

// Helper to send error responses with proper status codes.
func (h *RequestHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// Helper to send JSON responses.
func (h *RequestHandler) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// ============================================================================
// USER HANDLERS
// ============================================================================

// CreateUser handles POST /api/v1/users
func (h *RequestHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create new user
	result, err := h.service.CreateUser()
	if err != nil {
		log.Printf("Error creating user: %v", err)
		h.sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.sendJSON(w, http.StatusCreated, result)
}

// ListUsers handles GET /api/v1/users
func (h *RequestHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = DefaultPageSize
	}

	// Fetch users
	result, err := h.service.ListUsers(page, limit)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		h.sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.sendJSON(w, http.StatusOK, result)
}

// ============================================================================
// USER DELETE HANDLER
// ============================================================================

// DeleteUser handles DELETE /api/v1/users/{userID}
func (h *RequestHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	// Delete user
	err := h.service.DeleteUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendError(w, http.StatusNotFound, err.Error())
		} else {
			log.Printf("Error deleting user: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// ASSET HANDLERS
// ============================================================================

// CreateAsset handles POST /api/v1/assets
func (h *RequestHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Type == "" {
		h.sendError(w, http.StatusBadRequest, "type is required")
		return
	}

	if req.Data == nil || len(req.Data) == 0 {
		h.sendError(w, http.StatusBadRequest, "data is required")
		return
	}

	// Create asset
	asset, err := h.service.CreateAsset(req.Type, req.Data)
	if err != nil {
		if err.Error() == "invalid asset type" {
			h.sendError(w, http.StatusBadRequest, err.Error())
		} else {
			log.Printf("Error creating asset: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.sendJSON(w, http.StatusCreated, asset)
}

// ListAssets handles GET /api/v1/assets
func (h *RequestHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = DefaultPageSize
	}

	assetType := r.URL.Query().Get("type")
	var assetTypePtr *string
	if assetType != "" {
		assetTypePtr = &assetType
	}

	// Fetch assets
	result, err := h.service.ListAssets(page, limit, assetTypePtr)
	if err != nil {
		if err.Error() == "invalid asset type" {
			h.sendError(w, http.StatusBadRequest, err.Error())
		} else {
			log.Printf("Error listing assets: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.sendJSON(w, http.StatusOK, result)
}

// DeleteAsset handles DELETE /api/v1/assets/{assetID}
func (h *RequestHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["assetID"]

	// Delete asset
	err := h.service.DeleteAsset(assetID)
	if err != nil {
		if err.Error() == "asset not found" {
			h.sendError(w, http.StatusNotFound, err.Error())
		} else {
			log.Printf("Error deleting asset: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// FAVORITE HANDLERS
// ============================================================================

// GetFavorites handles GET /api/v1/users/{userID}/favorites
func (h *RequestHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = DefaultPageSize
	}

	assetType := r.URL.Query().Get("type")
	if assetType == "" {
		assetType = ""
	} else if !ValidAssetTypes[assetType] {
		h.sendError(w, http.StatusBadRequest, "invalid asset type")
		return
	}

	// Fetch favorites
	result, err := h.service.GetFavorites(userID, page, limit, &assetType)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendError(w, http.StatusNotFound, err.Error())
		} else {
			log.Printf("Error fetching favorites: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.sendJSON(w, http.StatusOK, result)
}

// AddFavorite handles POST /api/v1/users/{userID}/favorites
func (h *RequestHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	// Parse request body
	var req struct {
		AssetID     string `json:"asset_id"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AssetID == "" {
		h.sendError(w, http.StatusBadRequest, "asset_id is required")
		return
	}

	// Add favorite
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	favorite, err := h.service.AddFavorite(userID, req.AssetID, description)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "asset not found" {
			h.sendError(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "asset already in favorites" {
			h.sendError(w, http.StatusConflict, err.Error())
		} else {
			log.Printf("Error adding favorite: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.sendJSON(w, http.StatusCreated, favorite)
}

// UpdateFavorite handles PUT /api/v1/users/{userID}/favorites/{assetID}
func (h *RequestHandler) UpdateFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	assetID := vars["assetID"]

	// Parse request body
	var req struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Description == "" {
		h.sendError(w, http.StatusBadRequest, "description is required")
		return
	}

	// Update favorite
	favorite, err := h.service.UpdateFavoriteDescription(userID, assetID, req.Description)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "asset not in user's favorites" {
			h.sendError(w, http.StatusNotFound, err.Error())
		} else {
			log.Printf("Error updating favorite: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	h.sendJSON(w, http.StatusOK, favorite)
}

// RemoveFavorite handles DELETE /api/v1/users/{userID}/favorites/{assetID}
func (h *RequestHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	assetID := vars["assetID"]

	// Remove favorite
	err := h.service.RemoveFavorite(userID, assetID)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "asset not in user's favorites" {
			h.sendError(w, http.StatusNotFound, err.Error())
		} else {
			log.Printf("Error removing favorite: %v", err)
			h.sendError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck handles GET /health
func (h *RequestHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.sendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	// Initialize database
	storage, err := NewStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	// Create service and handler
	service := NewService(storage)
	handler := &RequestHandler{service: service}

	// Setup routes using gorilla/mux for better routing
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", handler.ListUsers).Methods("GET")
	api.HandleFunc("/users", handler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{userID}", handler.DeleteUser).Methods("DELETE")

	// Asset routes
	api.HandleFunc("/assets", handler.ListAssets).Methods("GET")
	api.HandleFunc("/assets", handler.CreateAsset).Methods("POST")
	api.HandleFunc("/assets/{assetID}", handler.DeleteAsset).Methods("DELETE")

	// Favorite routes
	api.HandleFunc("/users/{userID}/favorites", handler.GetFavorites).Methods("GET")
	api.HandleFunc("/users/{userID}/favorites", handler.AddFavorite).Methods("POST")
	api.HandleFunc("/users/{userID}/favorites/{assetID}", handler.UpdateFavorite).Methods("PUT")
	api.HandleFunc("/users/{userID}/favorites/{assetID}", handler.RemoveFavorite).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Start server
	// Using gorilla/mux router which is more robust than default mux
	log.Println("Starting server on :8080")
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  RequestTimeout,
		WriteTimeout: RequestTimeout,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

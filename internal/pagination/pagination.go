package pagination

import (
	"strconv"
)

type PaginationParams struct {
	Cursor string // Cursor string for pagination format: timestamp|asset_id
	Limit  int    // Number of items per page
}

const DefaultLimit = 20
const MaxLimit = 100

func ParsePaginationParams(cursor string, limitStr string) PaginationParams {
	limit := DefaultLimit

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	// Enforce limits
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return PaginationParams{
		Cursor: cursor,
		Limit:  limit,
	}
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total,omitempty"`
}

func NewPaginatedResponse(data interface{}, nextCursor string, hasMore bool, limit int) *PaginatedResponse {
	return &PaginatedResponse{
		Data: data,
		Meta: Meta{
			NextCursor: nextCursor,
			HasMore:    hasMore,
			Limit:      limit,
		},
	}
}

func (r *PaginatedResponse) WithTotal(total int) *PaginatedResponse {
	r.Meta.Total = total
	return r
}

package utils

import (
	"platform-go-challenge/internal/domain/entities"
	"sort"
)

// PaginatedResult sorts a slice of assets and returns a paginated sub-slice.
func PaginatedResult(assets []entities.Asset, page, size int) []entities.Asset {
	// Default Sort: Sort the assets slice by a consistent field, such as ID.
	// This ensures that pagination always returns the same results for a given page.
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].GetID() < assets[j].GetID()
	})

	// Calculate the start and end indices for pagination.
	start := (page - 1) * size
	end := start + size

	// Handle out-of-bounds cases for a non-existent page.
	if start >= len(assets) {
		return []entities.Asset{}
	}

	// Adjust the end index if it exceeds the slice's length.
	if end > len(assets) {
		end = len(assets)
	}

	// Return the paginated and sorted slice.
	return assets[start:end]
}

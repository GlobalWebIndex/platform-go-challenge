package main

import (
	"fmt"
	"go-challenge/models"
	"time"
)

// func ProcessAsset(assetType models.AssetType) {
// 	if assetType == models.AssetTypeChart {
// 		fmt.Println("Correct type:", assetType)
// 	} else {
// 		fmt.Println("Wrong type:", assetType)
// 	}

// }

// func main() {
// 	ProcessAsset("chart")
// 	fmt.Println("----------------")
// 	ProcessAsset(models.AssetTypeChart)

// 	// validation check
// 	assetType := models.AssetType(models.AssetTypeChart)

// 	if !assetType.IsValid() {
// 		fmt.Println("invalid asset type")
// 	} else {
// 		fmt.Println("Valid asset type")
// 	}
// }

// Version 1: Returns a COPY
func CreateFavouriteCopy() models.Favourite {
	fav := models.Favourite{
		ID:        "123",
		UserID:    "user-001",
		AssetID:   "asset-001",
		CreatedAt: time.Now(),
	}

	fmt.Printf("Inside function - Address: %p\n", &fav)
	return fav // Returns a COPY
}

// Version 2: Returns a POINTER
func CreateFavouritePointer() *models.Favourite {
	fav := models.Favourite{
		ID:        "123",
		UserID:    "user-001",
		AssetID:   "asset-001",
		CreatedAt: time.Now(),
	}

	fmt.Printf("Inside function - Address: %p\n", &fav)
	return &fav // Returns a POINTER (address)
}

// func main() {
// 	fmt.Println("=== Testing COPY ===")
// 	favCopy := CreateFavouriteCopy()
// 	fmt.Printf("Outside function - Address: %p\n\n", &favCopy)
// 	// Different addresses! Two copies exist.

// 	fmt.Println("=== Testing POINTER ===")
// 	favPointer := CreateFavouritePointer()
// 	fmt.Printf("Outside function - Address: %p\n", favPointer)
// 	// Same address! Only one copy exists.
// }

package testdata

import (
	"database/sql"
	"fmt"
	"log"
	"platform-go-challenge/internal/service/types"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

var Users = []string{"user1", "user2", "user3"}

func SeedDb(db *sql.DB) error {
	generator := NewGenerator()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Setup worker pools
	numWorkers := 4
	assetChan := make(chan types.Asset, numWorkers*100)
	resultChan := make(chan types.Asset, numWorkers*100)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	// Asset generation workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10000/numWorkers; i++ {
				assetChan <- generator.GenerateChart()
				assetChan <- generator.GenerateInsight()
				assetChan <- generator.GenerateAudience()
			}
		}()
	}

	// Close asset channel when generation is done
	go func() {
		wg.Wait()
		close(assetChan)
	}()

	// Prepare statement once
	stmt, err := tx.Prepare(`
        INSERT INTO assets (id, type, description, data, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare asset statement: %w", err)
	}
	defer stmt.Close()

	// Insert workers
	var insertWg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		insertWg.Add(1)
		go func() {
			defer insertWg.Done()
			for asset := range assetChan {
				jsonData, err := generator.GenerateAssetJSON(asset)
				if err != nil {
					select {
					case errChan <- fmt.Errorf("failed to generate JSON: %w", err):
					default:
					}
					return
				}

				_, err = stmt.Exec(
					asset.GetID(),
					asset.GetType(),
					asset.GetDescription(),
					jsonData,
					time.Now(),
					time.Now(),
				)
				if err != nil {
					select {
					case errChan <- fmt.Errorf("failed to insert asset: %w", err):
					default:
					}
					return
				}
				resultChan <- asset
			}
		}()
	}

	// Close result channel when all inserts are done
	go func() {
		insertWg.Wait()
		close(resultChan)
	}()

	// Collect results
	var assets []types.Asset
	for asset := range resultChan {
		assets = append(assets, asset)
	}

	// Check for errors
	select {
	case err := <-errChan:
		return err
	default:
	}

	// Insert favorites concurrently
	favStmt, err := tx.Prepare(`
        INSERT INTO favorites (user_id, asset_id, created_at)
        VALUES (?, ?, ?)
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare favorites statement: %w", err)
	}
	defer favStmt.Close()

	var favWg sync.WaitGroup
	favErrChan := make(chan error, 1)

	for _, userID := range Users {
		favWg.Add(1)
		go func(uid string) {
			defer favWg.Done()

			indices := make([]int, len(assets))
			for i := range indices {
				indices[i] = i
			}

			seed := uint64(time.Now().UnixNano())

			localRand := rand.New(rand.NewSource(uint64(seed)))

			localRand.Shuffle(len(indices), func(i, j int) {
				indices[i], indices[j] = indices[j], indices[i]
			})

			for i := 0; i < 20; i++ {
				asset := assets[indices[i]]
				_, err = favStmt.Exec(uid, asset.GetID(), time.Now())
				if err != nil {
					select {
					case favErrChan <- fmt.Errorf("failed to insert favorite: %w", err):
					default:
					}
					return
				}
			}
		}(userID)
	}

	favWg.Wait()
	close(favErrChan)

	if err := <-favErrChan; err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully seeded database with %d assets and favorites for %d users",
		len(assets), len(Users))
	return nil
}

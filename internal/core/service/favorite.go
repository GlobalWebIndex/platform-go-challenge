package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"time"

	"go-favorites-app/internal/core/domain/favorites"
	"go-favorites-app/internal/core/ports"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("internal/core/service")

type Service struct {
	repo     ports.FavoriteRepository
	cache    ports.Cache
	enricher ports.Enricher
	logger   *slog.Logger
}

func NewService(repo ports.FavoriteRepository, cache ports.Cache, enricher ports.Enricher, logger *slog.Logger) *Service {
	s := &Service{
		repo:     repo,
		cache:    cache,
		enricher: enricher,
		logger:   logger,
	}

	return s
}

func (s *Service) Save(ctx context.Context, asset favorites.Asset) error {
	ctx, span := tracer.Start(ctx, "Service.Save", trace.WithAttributes(
		attribute.String("asset.id", asset.GetID()),
		attribute.String("asset.type", string(asset.GetType())),
	))
	defer span.End()

	s.logger.InfoContext(ctx, "saving asset", "id", asset.GetID())

	// 1. Validate
	if err := asset.Validate(); err != nil {
		span.RecordError(err)
		return fmt.Errorf("validation failed: %w", err)
	}

	// 2. Save DB
	if err := s.repo.Save(ctx, asset); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to save to db: %w", err)
	}

	// 3. Write-Through (Enrich + Cache)
	// We perform this synchronously to ensure the user gets "Hydrated" data initially
	// and to prevent "loading" states in the UI.
	if _, err := s.enrichAndSaveCache(ctx, asset); err != nil {
		// Log but don't fail the request if enrichment causes issues (though enrichAndSaveCache essentially swallows errors)
		s.logger.Warn("failed to enrich and cache on save", "id", asset.GetID(), "error", err)
	}

	return nil
}

// enrichAndSaveCacheEnriched enriches the asset and updates the cache, returning the enriched result.
func (s *Service) enrichAndSaveCache(ctx context.Context, asset favorites.Asset) (favorites.Asset, error) {
	// 1. Enrich (Optimistic)
	// We operate directly on the asset (assuming caller gave us a safe copy or we are fine enriching in place)
	if err := s.enricher.Enrich(ctx, asset); err != nil {
		s.logger.Warn("enrichment failed, caching unenriched asset", "id", asset.GetID(), "error", err)
		// Proceed to cache anyway (Best Effort)
	}

	// 2. Cache
	s.updateCache(ctx, asset)

	return asset, nil
}

func (s *Service) updateCache(ctx context.Context, asset favorites.Asset) {
	data, err := json.Marshal(asset)
	if err != nil {
		s.logger.Error("failed to marshal asset for cache", "error", err)
		return
	}

	// ZAdd with timestamp score
	score := float64(time.Now().Unix())
	if err := s.cache.AddToSet(ctx, asset.GetID(), score); err != nil {
		s.logger.Error("failed to update cache set", "error", err)
	}
	if err := s.cache.Set(ctx, asset.GetID(), data); err != nil {
		s.logger.Error("failed to set cache data", "error", err)
	}
}

func (s *Service) FindByID(ctx context.Context, id string) (favorites.Asset, error) {
	ctx, span := tracer.Start(ctx, "Service.FindByID", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	batch, err := s.cache.GetBatch(ctx, []string{id})
	if err == nil && len(batch) > 0 {
		return s.unmarshal(batch[id])
	}

	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Write-through (optional read-repair)
	// We need to enrich the asset before caching it, as the cache is expected to hold enriched data.
	return s.enrichAndSaveCache(ctx, asset)
}

func (s *Service) FindAll(ctx context.Context, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	ctx, span := tracer.Start(ctx, "Service.FindAll", trace.WithAttributes(
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	))
	span.End() // End setup span

	// 1. Check Redis Set for IDs
	start := int64(offset)
	stop := int64(offset + limit - 1)

	ids, err := s.cache.GetIdsFromSet(ctx, start, stop)
	if err == nil && len(ids) > 0 {
		s.logger.Info("cache hit for favorites list (chunked)")
		return s.chunkedCacheIterator(ctx, ids), nil
	}

	// 2. Stream from DB
	s.logger.Info("streaming favorites from db")
	repoIter, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// 3. Cache and Return (No blocking enrichment on read)
	return s.cacheIterator(ctx, repoIter), nil
}

func (s *Service) chunkedCacheIterator(ctx context.Context, ids []string) iter.Seq2[favorites.Asset, error] {
	return func(yield func(favorites.Asset, error) bool) {
		const batchSize = 100
		for i := 0; i < len(ids); i += batchSize {
			end := i + batchSize
			if end > len(ids) {
				end = len(ids)
			}
			chunkIds := ids[i:end]

			dataMap, err := s.cache.GetBatch(ctx, chunkIds)
			if err != nil {
				yield(nil, fmt.Errorf("failed to fetch cache batch: %w", err))
				return
			}

			for _, id := range chunkIds {
				data, found := dataMap[id]
				if !found {
					// Read-Repair: ID exists in Set but Data missing in Hash/Set
					s.logger.Warn("cache inconsistency detected (missing data), repairing from db", "id", id)
					asset, err := s.repo.FindByID(ctx, id)
					if err != nil {
						yield(nil, fmt.Errorf("failed to repair cache for id %s: %w", id, err))
						return
					}

					// Synchronous Enrichment on Read Repair
					// Ensures FindAll returns fully enriched data
					enrichedAsset, _ := s.enrichAndSaveCache(ctx, asset)

					if !yield(enrichedAsset, nil) {
						return
					}
					continue
				}

				asset, err := s.unmarshal(data)
				if err != nil {
					yield(nil, fmt.Errorf("failed to unmarshal cached asset %s: %w", id, err))
					return
				}

				if !yield(asset, nil) {
					return
				}
			}
		}
	}
}

func (s *Service) FindAllByUser(ctx context.Context, userID string, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	ctx, span := tracer.Start(ctx, "Service.FindAllByUser", trace.WithAttributes(
		attribute.String("user.id", userID),
	))
	defer span.End()

	// Direct DB call for now (can add caching later)
	return s.repo.FindByUser(ctx, userID, limit, offset)
}

func (s *Service) Delete(ctx context.Context, id, userID string) error {
	ctx, span := tracer.Start(ctx, "Service.Delete")
	defer span.End()

	// Verify ownership
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if asset.GetUserID() != userID {
		return errors.New("forbidden: you do not own this asset")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return s.cache.Remove(ctx, id)
}

func (s *Service) UpdateDescription(ctx context.Context, id, description, userID string) (favorites.Asset, error) {
	ctx, span := tracer.Start(ctx, "Service.UpdateDescription")
	defer span.End()

	// Verify ownership
	asset, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset.GetUserID() != userID {
		return nil, errors.New("forbidden: you do not own this asset")
	}

	updatedAsset, err := s.repo.UpdateDescription(ctx, id, description)
	if err != nil {
		return nil, err
	}

	// Invalidate or update cache
	if err := s.cache.Remove(ctx, id); err != nil {
		// Log error but don't fail the operation as DB is already updated ??
		// Ideally we should have a way to retry or ensure consistency.
		// For now, logging.
		s.logger.Error("failed to invalidate cache after update", "id", id, "error", err)
	}

	return updatedAsset, nil
}

// Helpers

func (s *Service) unmarshal(data []byte) (favorites.Asset, error) {
	var base favorites.BaseAsset
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case favorites.AssetTypeChart:
		var c favorites.Chart
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, err
		}
		return c, nil
	case favorites.AssetTypeInsight:
		var i favorites.Insight
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}
		return i, nil
	case favorites.AssetTypeAudience:
		var a favorites.Audience
		if err := json.Unmarshal(data, &a); err != nil {
			return nil, err
		}
		return a, nil
	}
	return nil, fmt.Errorf("unknown type: %s", base.Type)
}

func (s *Service) cacheIterator(ctx context.Context, input iter.Seq2[favorites.Asset, error]) iter.Seq2[favorites.Asset, error] {
	return func(yield func(favorites.Asset, error) bool) {
		for asset, err := range input {
			if err != nil {
				yield(nil, err)
				return
			}

			// Synchronous Write-Through for List Streaming
			// We optimize for consistency over raw speed here to prevent "flickering" data
			if _, err := s.enrichAndSaveCache(ctx, asset); err != nil {
				s.logger.Warn("failed to enrich asset during stream", "id", asset.GetID())
			}

			// We return the enriched asset implicitly because enrichAndSaveCache updates cache
			// But here we are iterating over DB results (input).
			// If we want to return enriched, we should use result of enrichAndSaveCache
			// However, since we are streaming from DB, we might already have raw data.
			// Ideally we yield the result of enrichAndSaveCache.
			// For simplicity and matching current flow, let's just use the updated asset if possible.
			// But enrichAndSaveCache helper returns (Asset, error). Let's use it.

			// Actually, let's re-enrich.
			enriched, _ := s.enrichAndSaveCache(ctx, asset)

			if !yield(enriched, nil) {
				return
			}
		}
	}
}

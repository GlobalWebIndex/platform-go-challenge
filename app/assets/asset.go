package assets

import (
    "context"
    "gwi-platform/database"
    "gwi-platform/models"
)

type Asset interface {
    Add(ctx context.Context, db database.DB) error
    Delete(ctx context.Context, db database.DB, id models.ID) error
    Modify(ctx context.Context, db database.DB) error
    Get(ctx context.Context, db database.DB, id models.ID) (interface{}, error)
}

type AssetFactory func() Asset

var assetFactories = map[string]AssetFactory{
    "INSIGHT":  func() Asset { return &Insight{} },
    "CHART":    func() Asset { return &Chart{} },
    "AUDIENCE": func() Asset { return &Audience{} },
}

func GetAsset(assetType string) (Asset, bool) {
    factory, exists := assetFactories[assetType]
    if !exists {
        return nil, false
    }
    return factory(), true
}
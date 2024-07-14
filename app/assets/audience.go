package assets

import (
    "context"
    "database/sql"
    "fmt"
    "gwi-platform/models"
    "gwi-platform/database"
)

type Audience struct {
    models.BaseAsset
    Gender                  string  `json:"gender"`
    BirthCountry            string  `json:"birth_country"`
    AgeGroup                string  `json:"age_group"`
    DailyHoursOnSocialMedia float64 `json:"daily_hours_on_social_media"`
    PurchasesLastMonth      int     `json:"purchases_last_month"`
}

func (a *Audience) Add(ctx context.Context, db database.DB) error {
    stmt, err := db.PrepareContext(ctx, sqlInsertAudience)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()

    result, err := stmt.ExecContext(ctx, a.Gender, a.BirthCountry, a.AgeGroup, a.DailyHoursOnSocialMedia, a.PurchasesLastMonth)
    if err != nil {
        return fmt.Errorf("failed to insert audience: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert ID: %w", err)
    }

    a.AssetID = models.ID(id)
    a.AssetType = models.AudienceAsset
    return nil
}

func (a *Audience) Delete(ctx context.Context, db database.DB, id models.ID) error {
    _, err := db.ExecContext(ctx, sqlDeleteAudience, id)
    if err != nil {
        return fmt.Errorf("failed to delete audience: %w", err)
    }
    return nil
}

func (a *Audience) Modify(ctx context.Context, db database.DB) error {
    _, err := db.ExecContext(ctx, sqlUpdateAudience,
        a.Gender, a.BirthCountry, a.AgeGroup,
        a.DailyHoursOnSocialMedia, a.PurchasesLastMonth, a.AssetID)
    if err != nil {
        return fmt.Errorf("failed to update audience: %w", err)
    }
    return nil
}

func (a *Audience) Get(ctx context.Context, db database.DB, id models.ID) (interface{}, error) {
    var audience models.Audience
    err := db.QueryRowContext(ctx, sqlGetAudience, id).Scan(
        &audience.AssetID,
        &audience.Gender,
        &audience.BirthCountry,
        &audience.AgeGroup,
        &audience.DailyHoursOnSocialMedia,
        &audience.PurchasesLastMonth,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("audience with id %d not found", id)
        }
        return nil, fmt.Errorf("failed to get audience: %w", err)
    }
    audience.AssetType = models.AudienceAsset
    return audience, nil
}
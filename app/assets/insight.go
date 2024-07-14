package assets

import (
    "context"
    "database/sql"
    "fmt"
    "gwi-platform/models"
    "gwi-platform/database"
)

type Insight struct {
    models.BaseAsset
    Text string `json:"text"`
}

func (i *Insight) Add(ctx context.Context, db database.DB) error {
    stmt, err := db.PrepareContext(ctx, sqlInsertInsight)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()

    result, err := stmt.ExecContext(ctx, i.Text)
    if err != nil {
        return fmt.Errorf("failed to insert insight: %w", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert ID: %w", err)
    }

    i.AssetID = models.ID(id)
    i.AssetType = models.InsightAsset
    return nil
}

func (i *Insight) Delete(ctx context.Context, db database.DB, id models.ID) error {
    _, err := db.ExecContext(ctx, sqlDeleteInsight, id)
    if err != nil {
        return fmt.Errorf("failed to delete insight: %w", err)
    }
    return nil
}

func (i *Insight) Modify(ctx context.Context, db database.DB) error {
    _, err := db.ExecContext(ctx, sqlUpdateInsight, i.Text, i.AssetID)
    if err != nil {
        return fmt.Errorf("failed to update insight: %w", err)
    }
    return nil
}

func (i *Insight) Get(ctx context.Context, db database.DB, id models.ID) (interface{}, error) {
    var insight models.Insight
    err := db.QueryRowContext(ctx, sqlGetInsight, id).Scan(&insight.AssetID, &insight.Text)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("insight with id %d not found", id)
        }
        return nil, fmt.Errorf("failed to get insight: %w", err)
    }
    insight.AssetType = models.InsightAsset
    return insight, nil
}
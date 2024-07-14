package assets

import (
    "context"
    "database/sql"
    "fmt"
    "gwi-platform/models"
    "gwi-platform/database"
)

type Chart struct {
    models.BaseAsset
    Title      string           `json:"title"`
    XAxisTitle string           `json:"x_axis_title"`
    YAxisTitle string           `json:"y_axis_title"`
    DataPoints *[]models.DataPoint `json:"data_points"`
}

func (c *Chart) Add(ctx context.Context, db database.DB) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    result, err := tx.ExecContext(ctx, sqlInsertChart, c.Title, c.XAxisTitle, c.YAxisTitle)
    if err != nil {
        return fmt.Errorf("failed to insert chart: %w", err)
    }

    chartID, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert ID: %w", err)
    }

    stmt, err := tx.PrepareContext(ctx, sqlInsertDataPoint)
    if err != nil {
        return fmt.Errorf("failed to prepare data point statement: %w", err)
    }
    defer stmt.Close()

    for _, dp := range *c.DataPoints {
        _, err = stmt.ExecContext(ctx, chartID, dp.X, dp.Y)
        if err != nil {
            return fmt.Errorf("failed to insert data point: %w", err)
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    c.AssetID = models.ID(chartID)
    c.AssetType = models.ChartAsset
    return nil
}

func (c *Chart) Delete(ctx context.Context, db database.DB, id models.ID) error {
    _, err := db.ExecContext(ctx, sqlDeleteChart, id)
    if err != nil {
        return fmt.Errorf("failed to delete chart: %w", err)
    }
    return nil
}

func (c *Chart) Modify(ctx context.Context, db database.DB) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    _, err = tx.ExecContext(ctx, sqlUpdateChart, c.Title, c.XAxisTitle, c.YAxisTitle, c.AssetID)
    if err != nil {
        return fmt.Errorf("failed to update chart: %w", err)
    }

    _, err = tx.ExecContext(ctx, sqlDeleteChartDataPoints, c.AssetID)
    if err != nil {
        return fmt.Errorf("failed to delete existing data points: %w", err)
    }

    stmt, err := tx.PrepareContext(ctx, sqlInsertDataPoint)
    if err != nil {
        return fmt.Errorf("failed to prepare data point statement: %w", err)
    }
    defer stmt.Close()

    for _, dp := range *c.DataPoints {
        _, err = stmt.ExecContext(ctx, c.AssetID, dp.X, dp.Y)
        if err != nil {
            return fmt.Errorf("failed to insert data point: %w", err)
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

func (c *Chart) Get(ctx context.Context, db database.DB, id models.ID) (interface{}, error) {
    var chart models.Chart
    err := db.QueryRowContext(ctx, sqlGetChart, id).Scan(
        &chart.AssetID,
        &chart.Title,
        &chart.XAxisTitle,
        &chart.YAxisTitle,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("chart with id %d not found", id)
        }
        return nil, fmt.Errorf("failed to get chart: %w", err)
    }
    
    chart.AssetType = models.ChartAsset
    
    // Fetch data points
    dataPoints, err := getChartDataPoints(ctx, db, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get chart data points: %w", err)
    }
    chart.DataPoints = dataPoints

    return chart, nil
}

func getChartDataPoints(ctx context.Context, db database.DB, id models.ID) (*[]models.DataPoint, error) {
    rows, err := db.QueryContext(ctx, sqlGetChartDataPoints, id)
    if err != nil {
        return nil, fmt.Errorf("failed to query chart data points: %w", err)
    }
    defer rows.Close()

    var dataPoints []models.DataPoint
    for rows.Next() {
        var dp models.DataPoint
        if err := rows.Scan(&dp.X, &dp.Y); err != nil {
            return nil, fmt.Errorf("failed to scan data point: %w", err)
        }
        dataPoints = append(dataPoints, dp)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error after scanning data points: %w", err)
    }

    return &dataPoints, nil
}
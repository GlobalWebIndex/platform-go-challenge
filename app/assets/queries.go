package assets

const (
    // Chart queries
    sqlInsertChart = `
        INSERT INTO tcharts (title, x_axis_title, y_axis_title) 
        VALUES (?, ?, ?)
    `
    sqlInsertDataPoint = `
        INSERT INTO tchart_data_points (chart_id, x_value, y_value) 
        VALUES (?, ?, ?)
    `
    sqlUpdateChart = `
        UPDATE tcharts 
        SET title = ?, x_axis_title = ?, y_axis_title = ? 
        WHERE chart_id = ?
    `
    sqlDeleteChart = `
        DELETE FROM tcharts 
        WHERE chart_id = ?
    `
    sqlGetChart = `
        SELECT chart_id, title, x_axis_title, y_axis_title 
        FROM tcharts 
        WHERE chart_id = ?
    `
    sqlGetChartDataPoints = `
        SELECT x_value, y_value 
        FROM tchart_data_points 
        WHERE chart_id = ?
    `

    sqlDeleteChartDataPoints = `
        DELETE FROM tchart_data_points 
        WHERE chart_id = ?
    `

    // Insight queries
    sqlInsertInsight = `
        INSERT INTO tinsights (text) 
        VALUES (?)
    `
    sqlUpdateInsight = `
        UPDATE tinsights 
        SET text = ? 
        WHERE insight_id = ?
    `
    sqlDeleteInsight = `
        DELETE FROM tinsights 
        WHERE insight_id = ?
    `
    sqlGetInsight = `
        SELECT insight_id, text 
        FROM tinsights 
        WHERE insight_id = ?
    `

    // Audience queries
    sqlInsertAudience = `
        INSERT INTO taudiences 
        (gender, birth_country, age_group, daily_hours_on_social_media, purchases_last_month) 
        VALUES (?, ?, ?, ?, ?)
    `
    sqlUpdateAudience = `
        UPDATE taudiences 
        SET gender = ?, birth_country = ?, age_group = ?, 
            daily_hours_on_social_media = ?, purchases_last_month = ? 
        WHERE audience_id = ?
    `
    sqlDeleteAudience = `
        DELETE FROM taudiences 
        WHERE audience_id = ?
    `
    sqlGetAudience = `
        SELECT audience_id, gender, birth_country, age_group, 
               daily_hours_on_social_media, purchases_last_month 
        FROM taudiences 
        WHERE audience_id = ?
    `
)
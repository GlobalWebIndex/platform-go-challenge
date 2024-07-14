package handlers

const (
	sqlGetFavorites = `
    SELECT
        fa.fav_id,
        fa.asset_type,
        fa.asset_id,
        c.title,
        c.x_axis_title,
        c.y_axis_title,
        i.text,
        a.gender,
        a.birth_country,
        a.age_group,
        a.daily_hours_on_social_media,
        a.purchases_last_month
    FROM
        tfavorite_assets fa
    LEFT JOIN tcharts    c ON fa.asset_type = 'CHART'    AND fa.asset_id = c.chart_id
    LEFT JOIN tinsights  i ON fa.asset_type = 'INSIGHT'  AND fa.asset_id = i.insight_id
    LEFT JOIN taudiences a ON fa.asset_type = 'AUDIENCE' AND fa.asset_id = a.audience_id
    WHERE
        fa.shard = ?;`

	sqlGetDataPoints = `
    SELECT
        x_value,
        y_value
    FROM
        tchart_data_points
    WHERE
    	chart_id = ?;`

	sqlGetShards = `
	SELECT shard FROM tshards WHERE user_id = ?`

	sqlGetUser = `
	SELECT
        user_id,
        status
    FROM
        tusers
    WHERE
        user_id = ?`

	sqlAddUser = `
	INSERT INTO tusers (username, email) VALUES (?, ?)`

	sqlAddUserShards = `
	CALL CreateUserShards(?, ?)`
)

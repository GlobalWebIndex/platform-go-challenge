-- Users tables
CREATE TABLE tusers (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    status CHAR(1) DEFAULT 'A' NOT NULL,
    INDEX iusers_1 (user_id)
);

CREATE TABLE tusers_aud (
    aud_id INT AUTO_INCREMENT PRIMARY KEY,
    aud_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    action CHAR(1) NOT NULL,
    user_id INT NOT NULL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    status CHAR(1) NOT NULL
);

-- Charts tables
CREATE TABLE tcharts (
    chart_id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    x_axis_title VARCHAR(50) NOT NULL,
    y_axis_title VARCHAR(50) NOT NULL,
    INDEX icharts_1 (chart_id)
);

CREATE TABLE tcharts_aud (
    aud_id INT AUTO_INCREMENT PRIMARY KEY,
    aud_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    chart_id INT NOT NULL,
    title VARCHAR(100) NOT NULL,
    x_axis_title VARCHAR(50) NOT NULL,
    y_axis_title VARCHAR(50) NOT NULL,
    action CHAR(1) NOT NULL
);

CREATE TABLE tchart_data_points (
    data_points_id INT AUTO_INCREMENT PRIMARY KEY,
    chart_id INT NOT NULL,
    x_value FLOAT NOT NULL DEFAULT 0,
    y_value FLOAT NOT NULL DEFAULT 0,
    FOREIGN KEY (chart_id) REFERENCES tcharts(chart_id) ON DELETE CASCADE,
    INDEX ichart_data_points_1 (data_points_id),
    INDEX ichart_data_points_2 (chart_id,data_points_id)
);

CREATE TABLE tchart_data_points_aud (
    aud_id INT AUTO_INCREMENT PRIMARY KEY,
    aud_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_points_id INT NOT NULL,
    chart_id INT NOT NULL,
    x_value FLOAT NOT NULL,
    y_value FLOAT NOT NULL,
    action CHAR(1) NOT NULL
);

-- Insights tables
CREATE TABLE tinsights (
    insight_id INT AUTO_INCREMENT PRIMARY KEY,
    text TEXT NOT NULL,
    INDEX iinsights_1 (insight_id)
);

CREATE TABLE tinsights_aud (
    aud_id INT AUTO_INCREMENT PRIMARY KEY,
    aud_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    insight_id INT NOT NULL,
    text TEXT NOT NULL,
    action CHAR(1) NOT NULL
);

-- Audiences tables
CREATE TABLE taudiences (
    audience_id INT AUTO_INCREMENT PRIMARY KEY,
    gender ENUM('Male', 'Female', 'Other') NOT NULL,
    birth_country VARCHAR(50) NOT NULL,
    age_group VARCHAR(20) NOT NULL,
    daily_hours_on_social_media FLOAT NOT NULL DEFAULT 0,
    purchases_last_month INT NOT NULL DEFAULT 0,
    INDEX iaudiences_i1 (audience_id)
);

CREATE TABLE taudiences_aud (
    aud_id INT AUTO_INCREMENT PRIMARY KEY,
    aud_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    audience_id INT NOT NULL,
    gender ENUM('Male', 'Female', 'Other') NOT NULL,
    birth_country VARCHAR(50) NOT NULL,
    age_group VARCHAR(20) NOT NULL,
    daily_hours_on_social_media FLOAT NOT NULL,
    purchases_last_month INT NOT NULL,
    action CHAR(1) NOT NULL
);

CREATE TABLE tshards(
    user_id INT NOT NULL,
    shard VARCHAR(50) NOT NULL,
    count INT NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES tusers(user_id) ON DELETE CASCADE,
    INDEX ishards_1 (user_id),
    INDEX ishards_2 (shard),
    INDEX ishards_3 (user_id,shard)
);

-- Favorite Assets tables
CREATE TABLE tfavorite_assets (
    fav_id INT AUTO_INCREMENT PRIMARY KEY,
    asset_type ENUM('CHART', 'INSIGHT', 'AUDIENCE') NOT NULL,
    asset_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    shard VARCHAR(50) NOT NULL,
    FOREIGN KEY (shard) REFERENCES tshards(shard) ON DELETE CASCADE,
    UNIQUE KEY unique_favorite (shard, asset_id, asset_type),
    INDEX ifavorite_assets_1 (fav_id),
    INDEX ifavorite_assets_2 (shard, asset_id, asset_type)
);

CREATE TABLE tfavorite_assets_aud (
    aud_id INT AUTO_INCREMENT PRIMARY KEY,
    aud_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fav_id INT NOT NULL,
    asset_type ENUM('CHART', 'INSIGHT', 'AUDIENCE') NOT NULL,
    asset_id INT NOT NULL,
    created_at TIMESTAMP,
    shard VARCHAR(50) NOT NULL,
    action CHAR(1) NOT NULL
);


CREATE TABLE tsettings (
    setting_id INT AUTO_INCREMENT PRIMARY KEY,
    setting_name VARCHAR(100) NOT NULL,
    setting_value TEXT,
    value_type ENUM('VARCHAR', 'INT', 'FLOAT', 'BOOLEAN') NOT NULL,
    UNIQUE KEY unique_setting (setting_name)
);

-- init values for tsettings
INSERT INTO tsettings (setting_name,setting_value,value_type) VALUES ('PAGE_SIZE',5,'INT');


DELIMITER //


CREATE PROCEDURE CreateUserShards(
    IN p_user_id INT,
    IN p_num_shards INT
)
BEGIN
    DECLARE v_shard INT DEFAULT 1;
    
    -- Check if p_num_shards is positive
    IF p_num_shards <= 0 THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Number of shards must be positive';
    END IF;
    
    -- Loop to create the specified number of shards
    WHILE v_shard < p_num_shards DO
        -- Insert a new shard for the user
        INSERT INTO tshards (user_id, shard, count) 
        VALUES (p_user_id, CONCAT(p_user_id, '_', v_shard), 0);
        
        -- Increment the shard counter
        SET v_shard = v_shard + 1;
    END WHILE;
END //

CREATE PROCEDURE sp_audit_user(
    IN p_action CHAR(1),
    IN p_user_id INT,
    IN p_username VARCHAR(50),
    IN p_email VARCHAR(100),
    IN p_status CHAR(1)
)
BEGIN
    INSERT INTO tusers_aud (aud_time, user_id, username, email, status, action)
    VALUES (CURRENT_TIMESTAMP, p_user_id, p_username, p_email, p_status, p_action);
END //

CREATE PROCEDURE sp_audit_chart(
    IN p_action CHAR(1),
    IN p_chart_id INT,
    IN p_title VARCHAR(100),
    IN p_x_axis_title VARCHAR(50),
    IN p_y_axis_title VARCHAR(50)
)
BEGIN
    INSERT INTO tcharts_aud (
        aud_time, chart_id, title, x_axis_title, y_axis_title, action
    )
    VALUES (
        CURRENT_TIMESTAMP, p_chart_id, p_title, p_x_axis_title, p_y_axis_title, p_action
    );
END //

CREATE PROCEDURE sp_audit_chart_data_point(
    IN p_action CHAR(1),
    IN p_data_points_id INT,
    IN p_chart_id INT,
    IN p_x_value FLOAT,
    IN p_y_value FLOAT
)
BEGIN
    INSERT INTO tchart_data_points_aud (
        aud_time, data_points_id, chart_id, x_value, y_value, action
    )
    VALUES (
        CURRENT_TIMESTAMP, p_data_points_id, p_chart_id, p_x_value, p_y_value, p_action
    );
END //

CREATE PROCEDURE sp_audit_insight(
    IN p_action CHAR(1),
    IN p_insight_id INT,
    IN p_text TEXT
)
BEGIN
    INSERT INTO tinsights_aud (
        aud_time, insight_id, text, action
    )
    VALUES (
        CURRENT_TIMESTAMP, p_insight_id, p_text, p_action
    );
END //

CREATE PROCEDURE sp_audit_audience(
    IN p_action CHAR(1),
    IN p_audience_id INT,
    IN p_gender ENUM('Male', 'Female', 'Other'),
    IN p_birth_country VARCHAR(50),
    IN p_age_group VARCHAR(20),
    IN p_daily_hours_on_social_media FLOAT,
    IN p_purchases_last_month INT
)
BEGIN
    INSERT INTO taudiences_aud (
        aud_time, audience_id, gender, birth_country, age_group, 
        daily_hours_on_social_media, purchases_last_month, action
    )
    VALUES (
        CURRENT_TIMESTAMP, p_audience_id, p_gender, p_birth_country, p_age_group, 
        p_daily_hours_on_social_media, p_purchases_last_month, p_action
    );
END //

CREATE PROCEDURE sp_audit_favorite_asset(
    IN p_action CHAR(1),
    IN p_fav_id INT,
    IN p_asset_type ENUM('CHART', 'INSIGHT', 'AUDIENCE'),
    IN p_asset_id INT,
    IN p_created_at TIMESTAMP,
    IN p_shard VARCHAR(50)
)
BEGIN
    INSERT INTO tfavorite_assets_aud (
        aud_time, fav_id, asset_type, asset_id, created_at, shard, action
    )
    VALUES (
        CURRENT_TIMESTAMP, p_fav_id, p_asset_type, p_asset_id, p_created_at, p_shard, p_action
    );
END //

DELIMITER //

CREATE PROCEDURE InsertAndUpdateShard(
    IN p_user_id INT,
    IN p_asset_id INT,
    IN p_asset_type VARCHAR(50)
)
BEGIN
    DECLARE v_shard VARCHAR(50);
    DECLARE v_page_size INT;
    DECLARE v_next_page INT;

    -- Get the PAGE_SIZE setting
    SELECT CAST(setting_value AS SIGNED) INTO v_page_size
    FROM tsettings
    WHERE setting_name = 'PAGE_SIZE';

    -- Try to find an existing shard with available space
    SELECT
        shard INTO v_shard
    FROM
        tshards
    WHERE
        user_id = p_user_id
    AND count < v_page_size
    ORDER BY count ASC
    LIMIT 1;

    -- If no shard found, create a new one
    IF v_shard IS NULL THEN
        -- Find the next page number
        SELECT
            IFNULL(MAX(SUBSTRING_INDEX(shard, '_', -1)), 0) + 1 INTO v_next_page
        FROM
            tshards
        WHERE
            user_id = p_user_id;

        -- Create the new shard name
        SET v_shard = CONCAT(p_user_id, '_', v_next_page);

        -- Insert the new shard
        INSERT INTO tshards (user_id, shard, count)
        VALUES (p_user_id, v_shard, 0);
    END IF;

    -- Insert the favorite asset
    INSERT INTO tfavorite_assets (asset_id, asset_type, shard)
    VALUES (p_asset_id, p_asset_type, v_shard);

    -- Update the shard count
    UPDATE tshards
    SET count = count + 1
    WHERE
        user_id = p_user_id
    AND shard   = v_shard;
END //


-- Users triggers
CREATE TRIGGER tr_users_insert 
AFTER INSERT ON tusers
FOR EACH ROW
BEGIN
    CALL sp_audit_user('I', NEW.user_id, NEW.username, NEW.email, NEW.status);
END //

CREATE TRIGGER tr_users_update 
AFTER UPDATE ON tusers
FOR EACH ROW
BEGIN
    CALL sp_audit_user('U', OLD.user_id, OLD.username, OLD.email, OLD.status);
END //

CREATE TRIGGER tr_users_delete 
BEFORE DELETE ON tusers
FOR EACH ROW
BEGIN
    CALL sp_audit_user('D', OLD.user_id, OLD.username, OLD.email, OLD.status);
END //

-- Charts triggers
CREATE TRIGGER tr_charts_insert 
AFTER INSERT ON tcharts
FOR EACH ROW
BEGIN
    CALL sp_audit_chart(
        'I', NEW.chart_id, NEW.title, NEW.x_axis_title, NEW.y_axis_title
    );
END //

CREATE TRIGGER tr_charts_update 
AFTER UPDATE ON tcharts
FOR EACH ROW
BEGIN
    CALL sp_audit_chart(
        'U', OLD.chart_id, OLD.title, OLD.x_axis_title, OLD.y_axis_title
    );
END //

CREATE TRIGGER tr_charts_delete 
BEFORE DELETE ON tcharts
FOR EACH ROW
BEGIN
    CALL sp_audit_chart(
        'D', OLD.chart_id, OLD.title, OLD.x_axis_title, OLD.y_axis_title
    );
END //

-- Chart Data Points triggers
CREATE TRIGGER tr_chart_data_points_insert 
AFTER INSERT ON tchart_data_points
FOR EACH ROW
BEGIN
    CALL sp_audit_chart_data_point(
        'I', NEW.data_points_id, NEW.chart_id, NEW.x_value, NEW.y_value
    );
END //

CREATE TRIGGER tr_chart_data_points_update 
AFTER UPDATE ON tchart_data_points
FOR EACH ROW
BEGIN
    CALL sp_audit_chart_data_point(
        'U', OLD.data_points_id, OLD.chart_id, OLD.x_value, OLD.y_value
    );
END //

CREATE TRIGGER tr_chart_data_points_delete 
BEFORE DELETE ON tchart_data_points
FOR EACH ROW
BEGIN
    CALL sp_audit_chart_data_point(
        'D', OLD.data_points_id, OLD.chart_id, OLD.x_value, OLD.y_value
    );
END //

-- Insights triggers
CREATE TRIGGER tr_insights_insert 
AFTER INSERT ON tinsights
FOR EACH ROW
BEGIN
    CALL sp_audit_insight(
        'I', NEW.insight_id, NEW.text
    );
END //

CREATE TRIGGER tr_insights_update 
AFTER UPDATE ON tinsights
FOR EACH ROW
BEGIN
    CALL sp_audit_insight(
        'U', OLD.insight_id, OLD.text
    );
END //

CREATE TRIGGER tr_insights_delete 
BEFORE DELETE ON tinsights
FOR EACH ROW
BEGIN
    CALL sp_audit_insight(
        'D', OLD.insight_id, OLD.text
    );
END //

-- Audiences triggers
CREATE TRIGGER tr_audiences_insert 
AFTER INSERT ON taudiences
FOR EACH ROW
BEGIN
    CALL sp_audit_audience(
        'I', NEW.audience_id, NEW.gender, NEW.birth_country, NEW.age_group, 
        NEW.daily_hours_on_social_media, NEW.purchases_last_month
    );
END //

CREATE TRIGGER tr_audiences_update 
AFTER UPDATE ON taudiences
FOR EACH ROW
BEGIN
    CALL sp_audit_audience(
        'U', OLD.audience_id, OLD.gender, OLD.birth_country, OLD.age_group, 
        OLD.daily_hours_on_social_media, OLD.purchases_last_month
    );
END //

CREATE TRIGGER tr_audiences_delete 
BEFORE DELETE ON taudiences
FOR EACH ROW
BEGIN
    CALL sp_audit_audience(
        'D', OLD.audience_id, OLD.gender, OLD.birth_country, OLD.age_group, 
        OLD.daily_hours_on_social_media, OLD.purchases_last_month
    );
END //

-- Favorite Assets triggers
CREATE TRIGGER tr_favorite_assets_insert 
AFTER INSERT ON tfavorite_assets
FOR EACH ROW
BEGIN
    CALL sp_audit_favorite_asset(
        'I', NEW.fav_id, NEW.asset_type, NEW.asset_id, NEW.created_at, NEW.shard
    );
END //

CREATE TRIGGER tr_favorite_assets_update 
AFTER UPDATE ON tfavorite_assets
FOR EACH ROW
BEGIN
    CALL sp_audit_favorite_asset(
        'U', OLD.fav_id, OLD.asset_type, OLD.asset_id, OLD.created_at, OLD.shard
    );
END //

CREATE TRIGGER tr_favorite_assets_delete 
BEFORE DELETE ON tfavorite_assets
FOR EACH ROW
BEGIN
    CALL sp_audit_favorite_asset(
        'D', OLD.fav_id, OLD.asset_type, OLD.asset_id, OLD.created_at, OLD.shard
    );
END //


-- Workaround to ensure data integrity by verifying that asset_id matches the correct table 
-- based on the asset_type in tfavorite_assets. We use triggers to enforce this rule 
-- by checking the existence of asset_id in the appropriate table (tcharts, tinsights, or taudiences) 
-- before insert and update operations.

CREATE TRIGGER tfavorite_assets_insert BEFORE INSERT ON tfavorite_assets
FOR EACH ROW
BEGIN
    DECLARE asset_exists INT;

    IF NEW.asset_type = 'CHART' THEN
        SELECT COUNT(*) INTO asset_exists FROM tcharts WHERE chart_id = NEW.asset_id;
        IF asset_exists = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid asset_id for asset_type CHART';
        END IF;
    ELSEIF NEW.asset_type = 'INSIGHT' THEN
        SELECT COUNT(*) INTO asset_exists FROM tinsights WHERE insight_id = NEW.asset_id;
        IF asset_exists = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid asset_id for asset_type INSIGHT';
        END IF;
    ELSEIF NEW.asset_type = 'AUDIENCE' THEN
        SELECT COUNT(*) INTO asset_exists FROM taudiences WHERE audience_id = NEW.asset_id;
        IF asset_exists = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid asset_id for asset_type AUDIENCE';
        END IF;
    END IF;
END //

CREATE TRIGGER tfavorite_assets_update BEFORE UPDATE ON tfavorite_assets
FOR EACH ROW
BEGIN
    DECLARE asset_exists INT;

    IF NEW.asset_type = 'CHART' THEN
        SELECT COUNT(*) INTO asset_exists FROM tcharts WHERE chart_id = NEW.asset_id;
        IF asset_exists = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid asset_id for asset_type CHART';
        END IF;
    ELSEIF NEW.asset_type = 'INSIGHT' THEN
        SELECT COUNT(*) INTO asset_exists FROM tinsights WHERE insight_id = NEW.asset_id;
        IF asset_exists = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid asset_id for asset_type INSIGHT';
        END IF;
    ELSEIF NEW.asset_type = 'AUDIENCE' THEN
        SELECT COUNT(*) INTO asset_exists FROM taudiences WHERE audience_id = NEW.asset_id;
        IF asset_exists = 0 THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid asset_id for asset_type AUDIENCE';
        END IF;
    END IF;
END //

DELIMITER ;
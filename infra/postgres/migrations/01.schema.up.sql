-- Users table
CREATE TABLE users
(
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(255) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL
);

-- Assets table
CREATE TABLE assets
(
    id   SERIAL PRIMARY KEY,
    type VARCHAR(36) NOT NULL CHECK (type IN ('chart', 'insight', 'audience'))
);

-- Charts table (specialized asset)
CREATE TABLE charts
(
    asset_id     INT PRIMARY KEY REFERENCES assets (id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    title_x_axis VARCHAR(255),
    title_y_axis VARCHAR(255),
    data         JSONB        NOT NULL -- Using JSONB for flexible data structure
);

-- Insights table (specialized asset)
CREATE TABLE insights
(
    asset_id INT PRIMARY KEY REFERENCES assets (id) ON DELETE CASCADE,
    content  TEXT NOT NULL
);

-- Audiences table (specialized asset)
CREATE TABLE audiences
(
    asset_id                          INT PRIMARY KEY REFERENCES assets (id) ON DELETE CASCADE,
    gender                            VARCHAR(50),
    birth_country                     VARCHAR(100),
    age_group                         VARCHAR(50),
    hours_spent_daily_on_social_media INTEGER,
    last_month_purchases              INTEGER
);

-- Favorites table
CREATE TABLE favorites
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    asset_id    INTEGER NOT NULL REFERENCES assets (id) ON DELETE CASCADE,
    description TEXT,
    UNIQUE (user_id, asset_id)
);

-- Indexes for performance optimization
CREATE INDEX idx_favorites_user_id ON favorites (user_id);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255)
);

CREATE TABLE assets (
    asset_id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    asset_type VARCHAR(50),
    user_id UUID REFERENCES users(id)
);

CREATE TABLE favourites (
    user_id UUID REFERENCES users(id),
    asset_id VARCHAR(255) REFERENCES assets(asset_id),
    asset_type VARCHAR(50),
    PRIMARY KEY (user_id, asset_id)
);

CREATE TABLE charts (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    x_axis_title VARCHAR(255),
    y_axis_title VARCHAR(255)
);

CREATE TABLE chart_data (
    chart_id VARCHAR(255),
    datapoint_code VARCHAR(255),
    value NUMERIC
);

CREATE TABLE insights (
    id VARCHAR(255) PRIMARY KEY,
    description TEXT
);

CREATE TABLE audiences (
    id VARCHAR(255) PRIMARY KEY,
    gender VARCHAR(50),
    country VARCHAR(255),
    age_group VARCHAR(50),
    social_hours INT,
    purchases INT,
    description TEXT
);
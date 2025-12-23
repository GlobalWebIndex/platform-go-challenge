-- GWI Favorites Service - Database Schema
-- PostgreSQL 12+
--
-- Design principles:
-- 1. Single assets table with JSONB data: avoids three separate tables
-- 2. Soft deletes: preserves audit trail
-- 3. Proper indexes: query optimization
-- 4. Foreign keys: referential integrity
-- 5. UUIDs: distributed system friendly

-- ============================================================================
-- TABLES
-- ============================================================================

-- Users table (minimal - just identity)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Assets table (all types: chart, insight, audience)
-- Using JSONB for type-specific data keeps schema simple while remaining queryable
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('chart', 'insight', 'audience')),
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Favorites junction table
-- Links users to assets with optional description override
-- Soft deletes preserve data for auditing
CREATE TABLE IF NOT EXISTS favorites (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    description_override TEXT,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Why these indexes?
-- Most queries filter by user_id and need fast sorting.
-- deleted_at IS NULL is checked on every query.
-- Main query pattern: "get all favorites for user, sorted by date, excluding deleted"
-- This index makes that very fast by clustering data the way we query it

-- Filtered unique index: only enforce uniqueness for active (non-deleted) favorites
-- This allows soft deletes to work correctly - users can re-favorite deleted items
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_active_favorite_unique ON favorites (user_id, asset_id)
WHERE deleted_at IS NULL;

-- Sorted favorites for user (newest first)
CREATE INDEX IF NOT EXISTS idx_user_active_favorites ON favorites (user_id, added_at DESC)
WHERE deleted_at IS NULL;

-- Alternative sort: "get favorites by type"
CREATE INDEX IF NOT EXISTS idx_user_type_favorites ON favorites (user_id, asset_id)
WHERE deleted_at IS NULL;

-- Foreign key lookups
CREATE INDEX IF NOT EXISTS idx_favorite_asset ON favorites (asset_id);

-- Search by asset type (useful for filtering without joining)
CREATE INDEX IF NOT EXISTS idx_asset_type ON assets (type);

-- ============================================================================
-- HELPER FUNCTIONS
-- ============================================================================

-- Automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- TRIGGERS (PostgreSQL 12 compatible - drop then create)
-- ============================================================================

-- Drop triggers if they exist (PostgreSQL 12 doesn't support CREATE TRIGGER IF NOT EXISTS)
DROP TRIGGER IF EXISTS users_update_timestamp ON users;
DROP TRIGGER IF EXISTS assets_update_timestamp ON assets;

-- Apply the trigger to users table
CREATE TRIGGER users_update_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Apply the trigger to assets table
CREATE TRIGGER assets_update_timestamp
BEFORE UPDATE ON assets
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ============================================================================
-- SAMPLE DATA (for testing)
-- ============================================================================

-- Create test user (if not exists)
INSERT INTO users (id) 
VALUES ('550e8400-e29b-41d4-a716-446655440001')
ON CONFLICT (id) DO NOTHING;

-- Create sample chart asset (if not exists)
INSERT INTO assets (id, type, data) 
VALUES (
    '550e8400-e29b-41d4-a716-446655440101',
    'chart',
    '{"title": "Daily Social Media Usage", "x_axis": "Age Group", "y_axis": "Hours Per Day", "data": [3.5, 5.2, 4.8, 3.1, 2.0]}'::jsonb
)
ON CONFLICT (id) DO NOTHING;

-- Create sample insight asset (if not exists)
INSERT INTO assets (id, type, data) 
VALUES (
    '550e8400-e29b-41d4-a716-446655440102',
    'insight',
    '{"text": "40% of millennials spend more than 3 hours on social media daily", "topic": "millennial demographics"}'::jsonb
)
ON CONFLICT (id) DO NOTHING;

-- Create sample audience asset (if not exists)
INSERT INTO assets (id, type, data) 
VALUES (
    '550e8400-e29b-41d4-a716-446655440103',
    'audience',
    '{"gender": "Female", "birth_country": "US", "age_groups": ["25-34", "35-44"], "social_media_hours_daily": "3-5", "purchases_last_month": 5}'::jsonb
)
ON CONFLICT (id) DO NOTHING;

-- Add one to favorites (if not exists)
INSERT INTO favorites (id, user_id, asset_id, description_override) 
VALUES (
    '550e8400-e29b-41d4-a716-446655440201',
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440101',
    'Interesting trend to monitor'
)
ON CONFLICT (id) DO NOTHING;

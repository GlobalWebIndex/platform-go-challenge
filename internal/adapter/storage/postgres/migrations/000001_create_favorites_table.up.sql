CREATE TABLE IF NOT EXISTS favorites (
    id UUID PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    asset_data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_favorites_asset_data ON favorites USING GIN (asset_data);
CREATE INDEX IF NOT EXISTS idx_favorites_type ON favorites (type);

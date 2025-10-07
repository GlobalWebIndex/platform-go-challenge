CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     email TEXT UNIQUE
);

DO $$ BEGIN
CREATE TYPE asset_type AS ENUM ('chart', 'insight', 'audience');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS assets (
                                      id BIGSERIAL PRIMARY KEY,
                                      type asset_type NOT NULL,
                                      description TEXT,
                                      data JSONB NOT NULL,
                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS favourites (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_id BIGINT NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    custom_description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, asset_id)
    );

CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type);
CREATE INDEX IF NOT EXISTS idx_favourites_user ON favourites(user_id);

INSERT INTO users (email) VALUES ('gwi@demo.local') ON CONFLICT DO NOTHING;

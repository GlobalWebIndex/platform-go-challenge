BEGIN;

-- ===== users (new) =====
CREATE TABLE IF NOT EXISTS users (
  id          text PRIMARY KEY,
  email       text NOT NULL UNIQUE,
  full_name   text NOT NULL,
  is_active   boolean NOT NULL DEFAULT TRUE,
  created_at  timestamptz NOT NULL,
  updated_at  timestamptz NOT NULL
);

-- helpful index if you later search by email prefix
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- ===== assets (unchanged structure; snake_case columns already) =====
CREATE TABLE IF NOT EXISTS assets (
  id          text PRIMARY KEY,
  type        text NOT NULL CHECK (type IN ('chart','insight','audience')),
  title       text NOT NULL,
  description text NOT NULL,
  payload     jsonb NOT NULL,
  created_at  timestamptz NOT NULL,
  updated_at  timestamptz NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type);

-- ===== favorites (now references users + assets) =====
CREATE TABLE IF NOT EXISTS favorites (
  user_id    text NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
  asset_id   text NOT NULL REFERENCES assets(id)  ON DELETE CASCADE,
  created_at timestamptz NOT NULL,
  PRIMARY KEY (user_id, asset_id)
);

-- keyset pagination support
CREATE INDEX IF NOT EXISTS idx_fav_user_created
  ON favorites (user_id, created_at DESC, asset_id);

-- --- Safety nets if favorites existed before without FKs ---
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'fk_favorites_user_id'
  ) THEN
    ALTER TABLE favorites
      ADD CONSTRAINT fk_favorites_user_id
      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conname = 'fk_favorites_asset_id'
  ) THEN
    ALTER TABLE favorites
      ADD CONSTRAINT fk_favorites_asset_id
      FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE;
  END IF;
END $$;

COMMIT;

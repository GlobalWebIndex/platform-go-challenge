BEGIN;

-- ===== users =====
CREATE TABLE IF NOT EXISTS users (
  id          text PRIMARY KEY,
  email       text NOT NULL UNIQUE,
  full_name   text NOT NULL,
  is_active   boolean NOT NULL DEFAULT TRUE,
  created_at  timestamptz NOT NULL,
  updated_at  timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- ===== assets =====
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

-- ===== favorites (user-specific relation with overrides) =====
CREATE TABLE IF NOT EXISTS favorites (
  user_id              text NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
  asset_id             text NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
  created_at           timestamptz NOT NULL,
  override             jsonb NOT NULL DEFAULT '{}'::jsonb,   -- per-user fields (e.g., title/description)
  override_updated_at  timestamptz,
  PRIMARY KEY (user_id, asset_id)
);

-- keyset pagination support: (created_at DESC, asset_id ASC) per user
CREATE INDEX IF NOT EXISTS idx_fav_user_created
  ON favorites (user_id, created_at DESC, asset_id);

COMMIT;
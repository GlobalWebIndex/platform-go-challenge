-- CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT now()
);

-- ASSETS
CREATE TYPE asset_type AS ENUM ('chart', 'insight', 'audience');

CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type asset_type NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    data JSONB,
    created_at TIMESTAMP DEFAULT now()
);

-- FAVOURITES
CREATE TABLE favourites (
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    asset_id UUID REFERENCES assets(id) ON DELETE CASCADE,
  description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id, asset_id)
);

-- SEED DATA
INSERT INTO users (id, name, password_hash) VALUES
('u1', 'Alice', '$2a$10$q8PPH5ykvZ24Sq9Gu0QC6OfYVWYw5cQnczvDcUC3HWjDDixSf.I3.'),
('u2', 'Bob', '$2a$10$HXoxscLQW5pjFX3CxTyka./faujDQzJzlLgFPcE3zZ0cnd.RNRMHe');

INSERT INTO assets (type, title, description, data) VALUES
(
  'insight',
  'Social Media Usage',
  '40% of millennials spend more than 3 hours on social media daily',
  '{"source":"survey"}'
),
(
  'chart',
  'Purchases per Age Group',
  'Monthly purchases',
  '{"x":["18-24","25-34"],"y":[5,9]}'
);

-- Seed a couple of favourites linked to the inserted assets
INSERT INTO favourites (user_id, asset_id, description)
VALUES
  (
    'u1',
    (SELECT id FROM assets WHERE title = 'Social Media Usage' LIMIT 1),
    'Interesting social insight'
  ),
  (
    'u2',
    (SELECT id FROM assets WHERE title = 'Purchases per Age Group' LIMIT 1),
    'Chart I like to watch'
  )
ON CONFLICT DO NOTHING;

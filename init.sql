CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS assets (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) REFERENCES users(user_id),
    asset_id VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    data JSONB NOT NULL
);

-- Insert initial users
INSERT INTO users (user_id) VALUES
('user1'),
('user2'),
('user3');

-- Insert initial assets for user1
INSERT INTO assets (user_id, asset_id, type, description, data) VALUES
('user1', 'chart1', 'Chart', 'A test chart', '{"title": "Test Chart", "axisTitle": "Test Axis", "data": [0, 1, 2, 3, 4, 5]}'),
('user1', 'insight1', 'Insight', 'A test text', '{"text": "40% of millennials spend more than 3 hours on social media daily."}'),
('user1', 'audience1', 'Audience', 'A test Audience characteristics', '{"gender": "Male", "birthCountry": "Greece", "ageGroup": "24-35", "socialMediaHours": 3, "purchasesLastMonth": 6}');

-- Insert initial assets for user2
INSERT INTO assets (user_id, asset_id, type, description, data) VALUES
('user2', 'chart2', 'Chart', 'A test chart', '{"title": "Another Sample Chart", "axisTitle": "Test Axis", "data": [10, 20, 30, 40, 50]}'),
('user2', 'insight2', 'Insight', 'A test text', '{"text": "60% of Gen Z spends more than 5 hours on social media daily."}'),
('user2', 'audience2', 'Audience', 'A test Audience characteristics', '{"gender": "Female", "birthCountry": "Canada", "ageGroup": "18-24", "socialMediaHours": 5, "purchasesLastMonth": 10}');

-- Insert initial assets for user3
INSERT INTO assets (user_id, asset_id, type, description, data) VALUES
('user3', 'chart3', 'Chart', 'A test chart', '{"title": "Another Sample Chart", "axisTitle": "Test Axis", "data": [500, 0]}'),
('user3', 'insight3', 'Insight', 'A test text', '{"text": "90% of the third age spend zero time on social media."}'),
('user3', 'audience3', 'Audience', 'A test Audience characteristics', '{"gender": "Male", "birthCountry": "Greece", "ageGroup": "70-80", "socialMediaHours": 0, "purchasesLastMonth": 0}');
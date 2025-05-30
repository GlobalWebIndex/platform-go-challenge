BEGIN;

-- Insert users
INSERT INTO users (username, hashed_password)
VALUES ('alice', '$2b$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW'),    -- 1
       ('bob', '$2b$12$GdF9y8Yf1x3eRkBn5Z5rUu5vO6dJcQl7f0tYdXv7wXK1Ld0zV1XvG2'),     -- 2
       ('charlie', '$2b$12$4ZvXbY2X7bO5d5Z5q5j5v.8V0X7d0X7bO5d5Z5q5j5v.8V0X7d0X7b'), -- 3
       ('dave', '$2b$12$5Z5q5j5v.8V0X7d0X7bO5d5Z5q5j5v.8V0X7d0X7bO5d5Z5q5j5v.8V0'); -- 4 rows

INSERT INTO assets (type)
VALUES ('chart'),    -- 1
       ('audience'), -- 2
       ('chart'),    -- 3
       ('insight'),  -- 4
       ('audience'), -- 5
       ('insight'),  -- 6
       ('chart'),    -- 7
       ('insight'),  -- 8
       ('audience'), -- 9
       ('insight'),  -- 10
       ('chart'),    -- 11
       ('audience'); -- 12 rows


INSERT INTO charts (asset_id, title, title_x_axis, title_y_axis, data)
VALUES (1, 'Quarterly Sales', 'Quarters', 'Revenue (USD)', '[
  {
    "Q1": 12000
  },
  {
    "Q2": 18000
  },
  {
    "Q3": 15000
  },
  {
    "Q4": 22000
  }
]'::jsonb),
       (3, 'User Engagement', 'Days of Week', 'Active Users', '[
         {
           "Mon": 1200
         },
         {
           "Tue": 1400
         },
         {
           "Wed": 1600
         },
         {
           "Thu": 1800
         },
         {
           "Fri": 2100
         }
       ]'::jsonb),
       (7, 'Sales by Region', 'Regions', 'Sales (USD)', '[
         {
           "North": 30000
         },
         {
           "South": 25000
         },
         {
           "East": 40000
         },
         {
           "West": 35000
         }
       ]'::jsonb),
       (11, 'Product Popularity', 'Products', 'Sales Count', '[
         {
           "Product A": 1500
         },
         {
           "Product B": 2000
         },
         {
           "Product C": 1800
         }
       ]'::jsonb);


INSERT INTO audiences (asset_id, gender, birth_country, age_group, hours_spent_daily_on_social_media,
                       last_month_purchases)
VALUES (2, 'f', 'USA', '18-24', 3, 5),
       (5, 'm', 'Canada', '25-34', 2, 8),
       (9, 'f', 'UK', '35-44', 1, 10),
       (12, 'm', 'Germany', '45-54', 4, 2);

INSERT INTO insights (asset_id, content)
VALUES (4, 'Sales increased by 20% in Q4 compared to Q3 due to holiday season promotions.'),
       (6, 'User engagement peaks on Fridays with 75% higher activity compared to Mondays.'),
       (8, 'Monthly trends show a steady increase in visitors over the last three months.'),
       (10, 'Insights reveal that product A is the most popular among users aged 18-24.');

INSERT INTO favorites (user_id, asset_id, description)
VALUES ((SELECT id FROM users WHERE username = 'alice'),
        (SELECT id FROM assets WHERE type = 'chart' AND id = 1),
        'Primary sales chart'),
       ((SELECT id FROM users WHERE username = 'alice'),
        (SELECT id FROM assets WHERE type = 'insight' AND id = 4),
        'Key Q4 insight'),
       ((SELECT id FROM users WHERE username = 'bob'),
        (SELECT id FROM assets WHERE type = 'chart' AND id = 3),
        'Engagement tracking'),
       ((SELECT id FROM users WHERE username = 'bob'),
        (SELECT id FROM assets WHERE type = 'audience' AND id = 2),
        'US'),
       ((SELECT id FROM users WHERE username = 'charlie'),
        (SELECT id FROM assets WHERE type = 'insight' AND id = 6),
        'Friday peak analysis');

COMMIT;
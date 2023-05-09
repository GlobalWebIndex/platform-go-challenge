-- Connect to the database
-- \c host=172.18.0.1 port=5432 user=username dbname=mydb sslmode=disable password=password TimeZone=Europe/Athens

-- Insert 10 users
DO $$
    DECLARE
        i INTEGER := 1;
    BEGIN
        WHILE i <= 10 LOOP
                INSERT INTO "users" ("created_at", "updated_at", "deleted_at", "username", "password")
                VALUES (NOW(), NOW(), NULL, 'user' || i, 'password' || i);

                -- Insert 3 assets for each user
                INSERT INTO "insights" ("created_at", "updated_at", "deleted_at", "text", "description")
                VALUES (NOW(), NOW(), NULL, 'Insight ' || i || '1', 'Description for Insight ' || i || '1');
                INSERT INTO "audiences" ("created_at", "updated_at", "deleted_at", "gender", "birth_country", "age_group", "hours_spent_on_social", "purchases_last_month", "description")
                VALUES (NOW(), NOW(), NULL, 'Gender ' || i, 'Country ' || i, 'Age ' || i, i*10, i*5, 'Description for Audience ' || i);
                INSERT INTO "charts" ("created_at", "updated_at", "deleted_at", "title", "x_title", "y_title", "data", "description")
                VALUES (NOW(), NOW(), NULL, 'Chart ' || i, 'X Title ' || i, 'Y Title ' || i, 'Data for Chart ' || i, 'Description for Chart ' || i);

                -- Associate the assets with the user
                INSERT INTO "users_charts" ("user_id", "chart_id")
                VALUES (i, i);
                INSERT INTO "users_audiences" ("user_id", "audience_id")
                VALUES (i, i);
                INSERT INTO "users_insights" ("user_id", "insight_id")
                VALUES (i, i);

                i := i + 1;
            END LOOP;
    END $$;
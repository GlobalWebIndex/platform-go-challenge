-- Insert 10 users
DO
$$
    DECLARE
        i INTEGER := 1;
    BEGIN
        WHILE i <= 10
            LOOP
                INSERT INTO "users" ("created_at", "updated_at", "deleted_at", "username", "password")
                VALUES (NOW(), NOW(), NULL, 'user' || i, 'password' || i);

                -- Insert 1000 assets for each user
                DECLARE
                    j INTEGER := 1;
                BEGIN
                    WHILE j <= 10000
                        LOOP
                            INSERT INTO "insights" ("created_at", "updated_at", "deleted_at", "text", "description")
                            VALUES (NOW(), NOW(), NULL, 'Insight ' || i || '-' || j,
                                    'Description for Insight ' || i || '-' || j);
                            INSERT INTO "audiences" ("created_at", "updated_at", "deleted_at", "gender",
                                                     "birth_country", "age_group", "hours_spent_on_social",
                                                     "purchases_last_month", "description")
                            VALUES (NOW(), NOW(), NULL, 'Gender ' || i || '-' || j, 'Country ' || i || '-' || j,
                                    'Age ' || i || '-' || j, j * 10, j * 5, 'Description for Audience ' || i || '-' ||
                                                                            j);
                            INSERT INTO "charts" ("created_at", "updated_at", "deleted_at", "title", "x_title",
                                                  "y_title", "data", "description")
                            VALUES (NOW(), NOW(), NULL, 'Chart ' || i || '-' || j, 'X Title ' || i || '-' || j,
                                    'Y Title ' || i || '-' || j, 'Data for Chart ' || i || '-' || j,
                                    'Description for Chart ' || i || '-' || j);

                            -- Associate the assets with the user
                            INSERT INTO "users_charts" ("user_id", "chart_id")
                            VALUES (i, i * j);
                            INSERT INTO "users_audiences" ("user_id", "audience_id")
                            VALUES (i, i * j);
                            INSERT INTO "users_insights" ("user_id", "insight_id")
                            VALUES (i, i * j);

                            j := j + 1;
                        END LOOP;
                END;

                i := i + 1;
            END LOOP;
    END
$$;
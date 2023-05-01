DROP TABLE IF EXISTS tUser;
CREATE TABLE tUser (
    user_id VARCHAR(256) NOT NULL,
    username VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL,
    CONSTRAINT cUser_u1_user_id_username
        UNIQUE (user_id)
);

DROP TABLE IF EXISTS tAsset;
CREATE TABLE tAsset (
    asset_id SERIAL PRIMARY KEY,
    description VARCHAR(256),
    type VARCHAR(32) NOT NULL,
    data JSON NOT NULL,
    CONSTRAINT cAsset_c1
        CHECK (type in ('CHART', 'INSIGHT', 'AUDIENCE'))
);

DROP TABLE IF EXISTS tFavorite;
CREATE TABLE tFavorite (
    user_id VARCHAR(256) NOT NULL,
    asset_id BIGINT NOT NULL,
    CONSTRAINT cFavorite_f1 FOREIGN KEY (user_id)
        REFERENCES tUser (user_id),
    CONSTRAINT cFavorite_f2 FOREIGN KEY (asset_id)
        REFERENCES tAsset (asset_id)
);

INSERT INTO tUser (user_id, username, email) VALUES ('afefr2345sdfs', 'Bill', 'bill@gmail.com');
INSERT INTO tUser (user_id, username, email) VALUES ('gdfg90908', 'John', 'john@gmail.com');
INSERT INTO tUser (user_id, username, email) VALUES ('842wf887ghsfsdf', 'Nick', 'nick@gmail.com');


INSERT INTO tAsset (description, type, data) VALUES ('This is asset of type Chart', 'CHART', '{"X_axis_title":"X title","Y_axis_title":"Y title","data":"some data"}');
INSERT INTO tAsset (description, type, data) VALUES ('This is asset of type Insight', 'INSIGHT', '{"text":"40% of millenials spend more than 3hours on social media daily"}');
INSERT INTO tAsset (description, type, data) VALUES ('This is asset of type Audience', 'AUDIENCE', '{"gender":"Male","birth_country":"Greece","age_groups":"30-40","daily_hours_on_social":1,"purchases_last_month":3}');

INSERT INTO tFavorite (user_id, asset_id) VALUES ('afefr2345sdfs', 1);
INSERT INTO tFavorite (user_id, asset_id) VALUES ('afefr2345sdfs', 2);
INSERT INTO tFavorite (user_id, asset_id) VALUES ('afefr2345sdfs', 3);

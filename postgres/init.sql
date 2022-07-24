\c hypertube;

-- Users

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR (30) NOT NULL,
    firstname VARCHAR (30) NOT NULL,
    lastname VARCHAR (30) NOT NULL,
    email VARCHAR (90) UNIQUE NOT NULL,
    id_42 INTEGER UNIQUE,
    id_google VARCHAR (30) UNIQUE,
    password VARCHAR (65)
);

-- (used for tests)
-- unencrypted from server side password: 1e8392fcefc860ef9714dcf4ad2249a995118c7a3bdbf4a96e8ffd7fe354c2e6
-- user test 0 (for api-tester)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (0, 'username_test_0', 'firstname_test_0', 'lastname_test_0', 'email.test_0@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 1
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (101, 'username_test_1', 'firstname_test_1', 'lastname_test_1', 'email.test_1@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 2
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (102, 'username_test_2', 'firstname_test_2', 'lastname_test_2', 'email.test_2@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 3
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (103, 'username_test_3', 'firstname_test_3', 'lastname_test_3', 'email.test_3@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 4
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (104, 'username_test_4', 'firstname_test_4', 'lastname_test_4', 'email.test_4@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');
-- user test 5 (reserved for recover password apply)
INSERT INTO users (id, username, firstname, lastname, email, password)
VALUES (105, 'username_test_5', 'firstname_test_5', 'lastname_test_5', 'email.test_5@test.com', '1eba53d83fcffd42a3e3113fe52e68b8e9bbf478e29a12eb840557942386b482');

-- Person Informations

CREATE TABLE names (
    id BIGSERIAL PRIMARY KEY,
    tmdb_id INTEGER NOT NULL,
    name VARCHAR (250) NOT NULL,
    thumbnail VARCHAR (250) NULL,
    birth_year INTEGER NULL,
    death_year INTEGER NULL
);
CREATE UNIQUE INDEX unique_name_tmdb_id ON names(tmdb_id);

-- Genres

CREATE TABLE genres (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR (250) NOT NULL
);
CREATE UNIQUE INDEX unique_genre_name ON genres(name);

-- Media Informations

CREATE TABLE medias (
    id BIGSERIAL PRIMARY KEY,
    imdb_id VARCHAR (10) NULL,
    tmdb_id INTEGER NOT NULL,
    description TEXT NULL,
    duration INTEGER NULL,
    thumbnail TEXT NULL,
    background VARCHAR (250) NULL,
    year INTEGER NULL,
    rating REAL NULL
);
CREATE UNIQUE INDEX unique_media_imdb_id ON medias(imdb_id);

CREATE TABLE media_names (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    name VARCHAR (250) NOT NULL,
    lang VARCHAR (5) NOT NULL
);
ALTER TABLE ONLY media_names
ADD CONSTRAINT media_names_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;
CREATE UNIQUE INDEX unique_media_name_by_lang ON media_names(media_id, name, lang);

CREATE TABLE media_genres (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL
);
ALTER TABLE ONLY media_genres
ADD CONSTRAINT media_genres_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY media_genres
ADD CONSTRAINT media_genres_genre_id_foreign FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE NOT DEFERRABLE;
CREATE UNIQUE INDEX unique_media_genre_by_genre ON media_genres(media_id, genre_id);

CREATE TABLE media_staffs (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    name_id INTEGER NOT NULL,
    role VARCHAR (250) NULL
);
ALTER TABLE ONLY media_staffs
ADD CONSTRAINT media_staffs_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;
CREATE UNIQUE INDEX unique_media_staff_role ON media_staffs(media_id, name_id, role);

CREATE TABLE media_actors (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    name_id INTEGER NOT NULL,
    character VARCHAR (250) NULL
);
ALTER TABLE ONLY media_actors
ADD CONSTRAINT media_actors_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;
CREATE UNIQUE INDEX unique_media_actor_character ON media_actors(media_id, name_id, character);

-- Scrapper

CREATE TABLE torrents (
    id BIGSERIAL PRIMARY KEY,
    full_url TEXT NOT NULL,
    media_id INTEGER NULL,
    name VARCHAR (500) NOT NULL,
    type VARCHAR (50) NOT NULL,
    seed INTEGER NULL,
    leech INTEGER NULL,
    size VARCHAR (250) NULL,
    upload_time TIMESTAMP WITH TIME ZONE NULL,
    description_html TEXT NULL,
    torrent_url TEXT NULL,
    magnet TEXT NULL,
    last_update TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    file_path VARCHAR (500) NULL,
    downloaded BOOLEAN DEFAULT false,
    last_access TIMESTAMPTZ NULL
);
CREATE UNIQUE INDEX unique_torrent_url ON torrents(full_url);
ALTER TABLE ONLY torrents
ADD CONSTRAINT torrent_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE SET NULL;

CREATE TABLE torrent_files (
    id BIGSERIAL PRIMARY KEY,
    torrent_id INTEGER NOT NULL,
    path VARCHAR (250) NULL,
    name VARCHAR (250) NOT NULL,
    size VARCHAR (250) NULL
);
ALTER TABLE ONLY torrent_files
ADD CONSTRAINT torrent_files_torrent_id_foreign FOREIGN KEY (torrent_id) REFERENCES torrents(id) ON DELETE CASCADE NOT DEFERRABLE;

-- Comments

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	media_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
ALTER TABLE ONLY comments
ADD CONSTRAINT comments_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE ONLY comments
ADD CONSTRAINT comments_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE SET NULL;

-- Positions

CREATE TABLE positions (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    torrent_id INTEGER NOT NULL,
    position INTEGER
);

ALTER TABLE ONLY positions
ADD CONSTRAINT positions_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY positions
ADD CONSTRAINT positions_torrent_id_foreign FOREIGN KEY (torrent_id) REFERENCES torrents(id) ON DELETE CASCADE NOT DEFERRABLE;
CREATE UNIQUE INDEX unique_position ON positions(user_id, torrent_id);
ALTER TABLE ONLY positions
ADD CONSTRAINT unique_user_torrent_relation UNIQUE USING INDEX unique_position;

-- Subtitles

CREATE TABLE torrent_subtitles (
    id BIGSERIAL PRIMARY KEY,
	torrent_id INTEGER NOT NULL,
	lang VARCHAR (250) NOT NULL,
	path VARCHAR (250) NOT NULL
);
ALTER TABLE ONLY torrent_subtitles
ADD CONSTRAINT torrent_subtitles_foreign FOREIGN KEY (torrent_id) REFERENCES torrents(id) ON DELETE CASCADE;

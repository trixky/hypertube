-- Scrapper

CREATE TABLE torrents (
    id BIGSERIAL PRIMARY KEY,
    full_url TEXT NOT NULL,
    media_id INTEGER NULL,
    name VARCHAR (250) NOT NULL,
    type VARCHAR (50) NOT NULL,
    seed INTEGER NULL,
    leech INTEGER NULL,
    size VARCHAR (250) NULL,
    upload_time TIMESTAMP WITH TIME ZONE NULL,
    description_html TEXT NULL,
    torrent_url TEXT NULL,
    magnet TEXT NULL
);
CREATE UNIQUE INDEX unique_torrent_url ON torrents(full_url);

CREATE TABLE torrent_files (
    id BIGSERIAL PRIMARY KEY,
    torrent_id INTEGER NOT NULL,
    path VARCHAR (250) NULL,
    name VARCHAR (250) NOT NULL,
    size VARCHAR (250) NULL
);
ALTER TABLE ONLY torrent_files
ADD CONSTRAINT torrent_files_torrent_id_foreign FOREIGN KEY (torrent_id) REFERENCES torrents(id) ON DELETE CASCADE NOT DEFERRABLE;

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

-- Media Informations

CREATE TABLE medias (
    id BIGSERIAL PRIMARY KEY,
    imdb_id VARCHAR (10) NOT NULL,
    tmdb_id INTEGER NOT NULL,
    description TEXT NULL,
    duration INTEGER NULL,
    thumbnail TEXT NULL,
    background VARCHAR (250) NULL,
    year INTEGER NULL,
    genres VARCHAR (250) NOT NULL,
    rating REAL NULL
);
CREATE UNIQUE INDEX unique_media_imdb_id ON medias(imdb_id);
ALTER TABLE ONLY torrents
ADD CONSTRAINT torrent_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE SET NULL;

CREATE TABLE media_names (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    name VARCHAR (250) NOT NULL,
    lang VARCHAR (5) NOT NULL
);
ALTER TABLE ONLY media_names
ADD CONSTRAINT media_names_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;
CREATE UNIQUE INDEX unique_media_name_by_lang ON media_names(media_id, name, lang);

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

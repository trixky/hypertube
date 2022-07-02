-- Scrapper

CREATE TABLE torrents (
    id BIGSERIAL PRIMARY KEY,
    full_url TEXT NOT NULL,
    imdb_title_id INTEGER NULL,
    name VARCHAR (250) NOT NULL,
    type VARCHAR (50) NOT NULL,
    seed INTEGER NULL,
    leech INTEGER NULL,
    size VARCHAR (250) NULL,
    upload_time TIMESTAMP WITH TIME ZONE NULL,
    description_html TEXT NULL,
    torrent_url TEXT NULL,
    magnet TEXT NULL,
    imdb_id VARCHAR(10) NULL
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

-- IMDB Informations

CREATE TABLE medias (
    id BIGSERIAL PRIMARY KEY,
    imdb_id VARCHAR (10) NOT NULL,
    description TEXT NULL,
    duration INTEGER NOT NULL,
    thumbnail TEXT NULL,
    background VARCHAR (250) NULL,
    year INTEGER NOT NULL,
    genres VARCHAR (250) NOT NULL,
    rating INTEGER NULL
);
CREATE UNIQUE INDEX unique_media_imdb_id ON medias(imdb_id);

CREATE TABLE media_names (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    lang VARCHAR (5) NOT NULL,
    name VARCHAR (250) NOT NULL
);
ALTER TABLE ONLY media_names
ADD CONSTRAINT media_names_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;

CREATE TABLE media_staffs (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    name_id INTEGER NOT NULL,
    role VARCHAR (250) NOT NULL
);
ALTER TABLE ONLY media_staffs
ADD CONSTRAINT media_staffs_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;

CREATE TABLE media_relations (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    relation_id INTEGER NOT NULL
);
ALTER TABLE ONLY media_relations
ADD CONSTRAINT media_relations_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;

CREATE TABLE names (
    id BIGSERIAL PRIMARY KEY,
    imdb_id VARCHAR (10) NOT NULL,
    name VARCHAR (250) NOT NULL,
    birth_year INTEGER NULL,
    death_year INTEGER NULL
);
CREATE UNIQUE INDEX unique_name_imdb_id ON names(imdb_id);

CREATE TABLE name_relations (
    id BIGSERIAL PRIMARY KEY,
    name_id INTEGER NOT NULL,
    media_id INTEGER NOT NULL
);
ALTER TABLE ONLY name_relations
ADD CONSTRAINT name_relations_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;

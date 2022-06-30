-- Scrapper

CREATE TABLE medias (
    id BIGSERIAL PRIMARY KEY,
    full_url TEXT NOT NULL,
    imdb_title_id INTEGER NULL,
    name VARCHAR (250) NULL,
    type VARCHAR (50) NULL,
    seed INTEGER NULL,
    leech INTEGER NULL,
    size VARCHAR (250) NULL,
    upload_time TIMESTAMP WITH TIME ZONE NULL,
    description_html TEXT NULL,
    torrent_url TEXT NOT NULL,
    magnet TEXT NOT NULL,
    imdb_id VARCHAR(10) NOT NULL
);

CREATE TABLE media_files (
    id BIGSERIAL PRIMARY KEY,
    media_id INTEGER NOT NULL,
    path VARCHAR (250) NULL,
    name VARCHAR (250) NOT NULL,
    size VARCHAR (250) NULL
);
ALTER TABLE ONLY media_files
ADD CONSTRAINT media_files_media_id_foreign FOREIGN KEY (media_id) REFERENCES medias(id) ON DELETE CASCADE NOT DEFERRABLE;

-- IMDB Informations

CREATE TABLE imdb_titles (
    id BIGSERIAL PRIMARY KEY,
    imdb_id VARCHAR (10) NOT NULL,
    description TEXT NULL,
    duration VARCHAR (50) NULL,
    thumbnail TEXT NULL,
    background VARCHAR (250) NULL,
    year INTEGER NOT NULL,
    genres VARCHAR (250) NOT NULL,
    rating INTEGER NULL
);

CREATE TABLE imdb_title_names (
    id BIGSERIAL PRIMARY KEY,
    imdb_title_id INTEGER NOT NULL,
    lang VARCHAR (5) NOT NULL,
    name VARCHAR (250) NOT NULL
);
ALTER TABLE ONLY imdb_title_names
ADD CONSTRAINT imdb_title_names_imdb_title_id_foreign FOREIGN KEY (imdb_title_id) REFERENCES imdb_titles(id) ON DELETE CASCADE NOT DEFERRABLE;

CREATE TABLE imdb_title_staffs (
    id BIGSERIAL PRIMARY KEY,
    imdb_title_id INTEGER NOT NULL,
    name VARCHAR (250) NOT NULL,
    thumbnail TEXT NULL,
    url TEXT NOT NULL,
    role VARCHAR (250) NOT NULL
);
ALTER TABLE ONLY imdb_title_staffs
ADD CONSTRAINT imdb_title_staffs_imdb_title_id_foreign FOREIGN KEY (imdb_title_id) REFERENCES imdb_titles(id) ON DELETE CASCADE NOT DEFERRABLE;

CREATE TABLE imdb_title_relations (
    id BIGSERIAL PRIMARY KEY,
    imdb_title_id INTEGER NOT NULL,
    relation_imdb_id VARCHAR (10) NOT NULL,
    name VARCHAR (250) NOT NULL,
    thumbnail TEXT NULL
);
ALTER TABLE ONLY imdb_title_relations
ADD CONSTRAINT imdb_title_relations_imdb_title_id_foreign FOREIGN KEY (imdb_title_id) REFERENCES imdb_titles(id) ON DELETE CASCADE NOT DEFERRABLE;
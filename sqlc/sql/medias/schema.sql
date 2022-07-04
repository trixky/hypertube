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

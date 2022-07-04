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

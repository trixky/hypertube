-- Subtitles

CREATE TABLE torrent_subtitles (
    id BIGSERIAL PRIMARY KEY,
	torrent_id INTEGER NOT NULL,
	lang VARCHAR (250) NOT NULL,
	path VARCHAR (250) NOT NULL
);
ALTER TABLE ONLY torrent_subtitles
ADD CONSTRAINT torrent_subtitles_foreign FOREIGN KEY (torrent_id) REFERENCES torrents(id) ON DELETE CASCADE;

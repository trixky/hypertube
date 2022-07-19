-- Positions Informations

CREATE TABLE positions (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    torrent_id INTEGER NOT NULL,
    position INTEGER
);

CREATE UNIQUE INDEX unique_position ON positions(user_id, torrent_id);
ALTER TABLE ONLY positions
ADD CONSTRAINT unique_user_torrent_relation UNIQUE USING INDEX unique_position;
ALTER TABLE ONLY positions
ADD CONSTRAINT positions_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY positions
ADD CONSTRAINT positions_torrent_id_foreign FOREIGN KEY (torrent_id) REFERENCES torrents(id) ON DELETE CASCADE NOT DEFERRABLE;

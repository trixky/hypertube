-- name: SavePosition :one
INSERT INTO positions (user_id, torrent_id, position)
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT unique_user_torrent_relation DO UPDATE
SET position = EXCLUDED.position;

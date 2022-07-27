-- name: GetMediaByID :one
SELECT medias.*
FROM medias
WHERE medias.id = $1
LIMIT 1;

-- name: GetMediaNames :many
SELECT media_names.*
FROM media_names
WHERE media_names.media_id = $1;

-- name: GetMediaTorrents :many
SELECT torrents.*
FROM torrents
RIGHT JOIN medias ON torrents.media_id = medias.id
WHERE torrents.media_id = $1
ORDER BY seed DESC, leech ASC;

-- name: GetMediaTorrentsForUser :many
SELECT torrents.*, positions.position
FROM torrents
RIGHT JOIN medias ON torrents.media_id = medias.id
LEFT JOIN positions ON torrents.id = positions.torrent_id AND positions.user_id = $2
WHERE torrents.media_id = $1
ORDER BY seed DESC, leech ASC;

-- name: WatchedMedia :one
SELECT * FROM positions
RIGHT JOIN torrents ON torrents.id = positions.torrent_id
WHERE positions.user_id = $1 AND torrents.media_id = $2
LIMIT 1;

-- name: GetMediaGenres :many
SELECT genres.*
FROM media_genres
RIGHT JOIN genres ON media_genres.genre_id = genres.id
WHERE media_genres.media_id = $1;

-- name: GetMediaActors :many
SELECT names.id, media_actors.id as actor_id, names.name, names.thumbnail, media_actors.character
FROM media_actors
RIGHT JOIN names ON media_actors.name_id = names.id
WHERE media_actors.media_id = $1
	AND names.thumbnail IS NOT NULL
	AND media_actors.character IS NOT NULL
ORDER BY cast_order
LIMIT 15;

-- name: GetMediaStaffs :many
SELECT names.id, names.name, names.thumbnail, media_staffs.role
FROM media_staffs
RIGHT JOIN names ON media_staffs.name_id = names.id
WHERE media_staffs.media_id = $1
	AND names.thumbnail IS NOT NULL
	AND media_staffs.role IS NOT NULL;

-- name: GetGenres :many
SELECT *
FROM genres;

-- name: GetTorrent :one
SELECT *
FROM torrents
WHERE id = $1
LIMIT 1;

-- name: GetTorrentByURL :one
SELECT *
FROM torrents
WHERE full_url = $1
LIMIT 1;

-- name: ListTorrents :many
SELECT *
FROM torrents
ORDER BY id DESC;

-- name: CreateTorrent :one
INSERT INTO torrents
	(
		full_url,
		media_id,
		name,
		type,
		seed,
		leech,
		size,
		upload_time,
		description_html,
		torrent_url,
		magnet
	)
VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11
	)
RETURNING *;

-- name: SetTorrentPeers :exec
UPDATE torrents
SET seed = $2,
	leech = $3
WHERE id = $1;

-- name: SetTorrentInformations :exec
UPDATE torrents
SET torrent_url = $2,
	magnet = $3,
	description_html = $4,
	size = $5
WHERE id = $1;

-- name: AddTorrentMediaId :exec
UPDATE torrents
SET media_id = $2
WHERE id = $1;

-- name: AddTorrentFile :one
INSERT INTO torrent_files
	(torrent_id, path, name, size)
VALUES
	($1, $2, $3, $4)
RETURNING *;

-- name: DeleteTorrent :exec
DELETE FROM torrents
WHERE id = $1;

-- name: CreateMedia :one
INSERT INTO medias
	(
		imdb_id,
		tmdb_id,
		description,
		duration,
		thumbnail,
		background,
		year,
		genres,
		rating
	)
VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9
	)
RETURNING *;

-- name: CreateMedias :copyfrom
INSERT INTO medias
	(
		imdb_id,
		description,
		duration,
		thumbnail,
		background,
		year,
		genres,
		rating
	)
VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	);

-- name: GetMediaByIMDB :one
SELECT *
FROM medias
WHERE imdb_id = $1
LIMIT 1;

-- name: CheckMediaExistByIMDB :one
SELECT count(id)
FROM medias
WHERE imdb_id = $1
LIMIT 1;

-- name: CreateMediaName :exec
INSERT INTO media_names
	(media_id, name, lang)
VALUES
	($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: CheckMediaNameExist :one
SELECT count(id)
FROM media_names
WHERE media_id = $1 AND name = $2 AND lang = $3
LIMIT 1;

-- name: DeleteMediaName :exec
DELETE FROM media_names
WHERE id = $1;

-- name: CreateName :one
INSERT INTO names
	(
		tmdb_id,
		name,
		thumbnail,
		birth_year,
		death_year
	)
VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5
	)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetNameByTMDB :one
SELECT *
FROM names
WHERE tmdb_id = $1
LIMIT 1;

-- name: CheckNameExistByTMDB :one
SELECT count(id)
FROM names
WHERE tmdb_id = $1
LIMIT 1;

-- name: CreateMediaStaff :exec
INSERT INTO media_staffs
	(media_id, name_id, role)
VALUES
	($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: DeleteMediaStaff :exec
DELETE FROM media_staffs
WHERE id = $1;

-- name: CreateMediaActor :exec
INSERT INTO media_actors
	(media_id, name_id, character)
VALUES
	($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: DeleteMediaActor :exec
DELETE FROM media_actors
WHERE id = $1;

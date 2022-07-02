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
		imdb_title_id,
		name,
		type,
		seed,
		leech,
		size,
		upload_time,
		description_html,
		torrent_url,
		magnet,
		imdb_id
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
		$11,
		$12
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
	size = $5,
	imdb_id = $6
WHERE id = $1;

-- name: AddTorrentIMDBId :exec
UPDATE torrents
SET imdb_title_id = $2
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

-- name: CreateName :exec
INSERT INTO names
	(
		imdb_id,
		name,
		birth_year,
		death_year
	)
VALUES
	(
		$1,
		$2,
		$3,
		$4
	)
ON CONFLICT DO NOTHING;

-- name: GetNameByIMDB :one
SELECT *
FROM names
WHERE imdb_id = $1
LIMIT 1;

-- name: CheckNameExistByIMDB :one
SELECT count(id)
FROM names
WHERE imdb_id = $1
LIMIT 1;

-- name: CreateNameRelation :exec
INSERT INTO name_relations
	(name_id, media_id)
VALUES
	($1, $2);

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

-- name: CreateMediaRelation :one
INSERT INTO media_relations
	(media_id, relation_id)
VALUES
	($1, $2)
RETURNING *;

-- name: DeleteMediaRelation :exec
DELETE FROM media_relations
WHERE id = $1;

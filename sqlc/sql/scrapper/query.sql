-- name: GetTorrent :one
SELECT *
FROM torrents
WHERE id = $1
LIMIT 1;
-- name: ListTorrents :many
SELECT *
FROM torrents
ORDER BY id DESC;
-- name: CreateTorrent :one
INSERT INTO torrents (
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
VALUES (
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
-- name: AddTorrentFile :exec
INSERT INTO torrent_files (torrent_id, path, name, size)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: DeleteTorrent :exec
DELETE FROM torrents
WHERE id = $1;
-- name: CreateMedia :one
INSERT INTO medias (
		imdb_id,
		description,
		duration,
		thumbnail,
		background,
		year,
		genres,
		rating
	)
VALUES (
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
-- name: AddMediaName :one
INSERT INTO media_names (media_id, lang, name)
VALUES ($1, $2, $3)
RETURNING *;
-- name: DeleteMediaName :exec
DELETE FROM media_names
WHERE id = $1;
-- name: AddMediaStaff :one
INSERT INTO media_staffs (media_id, name, thumbnail, url, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: DeleteMediaStaff :exec
DELETE FROM media_staffs
WHERE id = $1;
-- name: AddMediaRelation :one
INSERT INTO media_relations (
		media_id,
		relation_imdb_id,
		name,
		thumbnail
	)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: DeleteMediaRelation :exec
DELETE FROM media_relations
WHERE id = $1;
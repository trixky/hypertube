-- name: GetMedia :one
SELECT *
FROM medias
WHERE id = $1
LIMIT 1;
-- name: ListMedias :many
SELECT *
FROM medias
ORDER BY id DESC;
-- name: CreateMedia :one
INSERT INTO medias (
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
-- name: SetMediaPeers :exec
UPDATE medias
SET seed = $2,
	leech = $3
WHERE id = $1;
-- name: AddMediaIMDBId :exec
UPDATE medias
SET imdb_title_id = $2
WHERE id = $1;
-- name: AddMediaFile :exec
INSERT INTO media_files (media_id, path, name, size)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: DeleteMedia :exec
DELETE FROM medias
WHERE id = $1;
-- name: CreateIMDBTitle :one
INSERT INTO imdb_titles (
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
-- name: AddIMDBTitleName :one
INSERT INTO imdb_title_names (imdb_title_id, lang, name)
VALUES ($1, $2, $3)
RETURNING *;
-- name: DeleteIMDBTitleName :exec
DELETE FROM imdb_title_names
WHERE id = $1;
-- name: AddIMDBTitleStaff :one
INSERT INTO imdb_title_staffs (imdb_title_id, name, thumbnail, url, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: DeleteIMDBTitleStaff :exec
DELETE FROM imdb_title_staffs
WHERE id = $1;
-- name: AddIMDBTitleRelation :one
INSERT INTO imdb_title_relations (
		imdb_title_id,
		relation_imdb_id,
		name,
		thumbnail
	)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: DeleteIMDBTitleRelation :exec
DELETE FROM imdb_title_relations
WHERE id = $1;
-- name: GetMediaByID :one
SELECT *
FROM medias
WHERE id = $1
LIMIT 1;

-- name: GetGenres :many
SELECT *
FROM genres;

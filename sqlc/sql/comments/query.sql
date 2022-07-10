-- name: GetCommentById :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: GetMediaComments :many
SELECT *
FROM comments
WHERE media_id = $1;

-- name: GetUserComments :many
SELECT *
FROM comments
WHERE user_id = $1;

-- name: GetUserCommentsForMedia :many
SELECT *
FROM comments
WHERE user_id = $1 AND media_id = $2;

-- name: CreateComment :one
INSERT INTO comments
	(user_id, media_id, content)
VALUES
	($1, $2, $3)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

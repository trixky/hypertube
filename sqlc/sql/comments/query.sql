-- name: GetCommentById :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: GetMediaComments :many
SELECT users.username, comments.id, comments.user_id, comments.content, comments.created_at
FROM comments
RIGHT JOIN users ON comments.user_id = users.id
WHERE media_id = $1
ORDER BY id DESC;

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

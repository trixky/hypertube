-- name: GetCommentById :one
SELECT  users.username, comments.*
FROM comments
RIGHT JOIN users ON comments.user_id = users.id
WHERE comments.id = $1 LIMIT 1;

-- name: GetLastComments :many
SELECT users.username, comments.*
FROM comments
RIGHT JOIN users ON comments.user_id = users.id
ORDER BY id DESC LIMIT 15;

-- name: GetMediaComments :many
SELECT users.username, comments.*
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

-- name: UpdateComment :exec
UPDATE comments SET content = $2 WHERE id = $1;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

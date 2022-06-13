-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (
    username, firstname, lastname, email, password
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateUserUsername :exec
UPDATE users
set username = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
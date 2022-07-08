-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetInternalUserByCredentials :one
SELECT * FROM users
WHERE email = $1 
AND password = $2
AND id_42 IS NULL
LIMIT 1;

-- name: GetUserBy42Id :one
SELECT * FROM users
WHERE id_42 = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username;

-- name: CreateInternalUser :one
INSERT INTO users (
    username, firstname, lastname, email, password
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: Create42ExternalUser :one
INSERT INTO users (
    username, firstname, lastname, email, password, id_42
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
set username = $2,
firstname = $3,
lastname = $4,
email = $5,
password = $6
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
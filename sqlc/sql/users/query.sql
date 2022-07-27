-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetInternalUserByCredentials :one
SELECT * FROM users
WHERE email = $1 
AND password = $2
AND id_42 IS NULL
AND id_google IS NULL
LIMIT 1;

-- name: GetInternalUserByEmail :one
SELECT * FROM users
WHERE email = $1
AND id_42 IS NULL
AND id_google IS NULL
LIMIT 1;

-- name: GetUserBy42Id :one
SELECT * FROM users
WHERE id_42 = $1 LIMIT 1;

-- name: GetUserByGoogleId :one
SELECT * FROM users
WHERE id_google = $1 LIMIT 1;

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

-- name: CreateGoogleExternalUser :one
INSERT INTO users (
    username, firstname, lastname, email, password, id_google
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

-- name: UpdateUserPassword :exec
UPDATE users
set password = $2
WHERE id = $1
AND id_42 IS NULL
AND id_google IS NULL;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUserMedias :one
SELECT COUNT(distinct medias.id) as total
FROM torrents
RIGHT JOIN medias ON torrents.media_id = medias.id
RIGHT JOIN positions ON torrents.id = positions.torrent_id AND positions.user_id = $1;

-- name: GetUserMedias :many
SELECT DISTINCT ON (medias.id) medias.*, positions.id
FROM torrents
RIGHT JOIN medias ON torrents.media_id = medias.id
RIGHT JOIN positions ON torrents.id = positions.torrent_id AND positions.user_id = $2
ORDER BY medias.id DESC, positions.id DESC
LIMIT 50 OFFSET $1;

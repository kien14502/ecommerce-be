-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    password_hash,
    username,
    full_name
) VALUES (?, ?, ?, ?, ?);

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;
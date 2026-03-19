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
WHERE email = ?;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * From users
WHERE username = ?;

-- name: GetEmailVerifiedStatus :one
SELECT is_email_verified
FROM users
WHERE email = ?;

-- name: MarkEmailVerified :exec
UPDATE users
SET is_email_verified = TRUE
WHERE email = ?;

-- name: CreatePost :exec
INSERT INTO posts (
    id,
    user_id,
    content
) VALUES (?, ?, ?);

-- name: GetPost :one
SELECT * FROM posts
WHERE id = ?;

-- name: GetUserPosts :many
SELECT * FROM posts
WHERE user_id = ?
ORDER BY created_at DESC;
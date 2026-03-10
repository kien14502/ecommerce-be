-- name: CreatePost :exec
INSERT INTO posts (
    id,
    user_id,
    content
) VALUES (?, ?, ?);

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: GetUserPosts :many
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at DESC;
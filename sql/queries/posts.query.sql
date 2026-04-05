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
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdatePost :exec
UPDATE posts
SET
    content    = ?,
    updated_at = NOW()
WHERE id      = ?
  AND user_id = ?;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id      = ?
  AND user_id = ?;

-- name: GetUserPostsWithCount :many
SELECT 
    id,
    user_id,
    content,
    visibility,
    created_at,
    CAST(COUNT(*) OVER() AS SIGNED) AS total_count
FROM posts
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
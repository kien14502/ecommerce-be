-- name: CreateComment :exec
INSERT INTO comments (
    id,
    post_id,
    user_id,
    parent_id,
    content
) VALUES (?, ?, ?, ?, ?);

-- name: GetPostComments :many
SELECT * FROM comments
WHERE post_id = ?
ORDER BY created_at;
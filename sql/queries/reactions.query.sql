-- name: CreateReaction :exec
INSERT INTO reactions (
    id,
    user_id,
    post_id,
    reaction_type
) VALUES (?, ?, ?, ?);

-- name: DeleteReaction :exec
DELETE FROM reactions
WHERE user_id = $1
AND post_id = $2;
-- name: FollowUser :exec
INSERT INTO follows (
    follower_id, following_id
) VALUES (
    ?,?
);

-- name: UnfollowUser :exec
DELETE FROM follows
WHERE follower_id = ?
AND following_id = ?;
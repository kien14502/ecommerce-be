-- name: FollowUser :exec
INSERT INTO follows (
    follower_id, following_id
) VALUES (
    $1, $2
);

-- name: UnfollowUser :exec
DELETE FROM follows
WHERE follower_id = $1
AND following_id = $2;
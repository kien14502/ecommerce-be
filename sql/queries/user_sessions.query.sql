-- name: CreateSession :exec
INSERT INTO user_sessions (
    id,
    user_id,
    device_id,
    refresh_token_hash,
    expires_at
) VALUES (?, ?, ?, ?, ?);

-- name: GetSessionByToken :one
SELECT * FROM user_sessions
WHERE refresh_token_hash = ?
AND expires_at > NOW();

-- name: DeleteSession :exec
DELETE FROM user_sessions
WHERE id = ?;

-- name: DeleteAllSessions :exec
DELETE FROM user_sessions
WHERE user_id = ?;
-- name: CreateDevice :exec
INSERT INTO user_devices (
    id,
    user_id,
    device_name,
    device_type,
    user_agent,
    ip_address
) VALUES (?, ?, ?, ?, ?, ?);

-- name: ListUserDevices :many
SELECT d.*
FROM user_devices d
JOIN user_sessions s
ON s.device_id = d.id
WHERE s.user_id = $1;
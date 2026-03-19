-- name: CreateOTP :exec
INSERT INTO otp_verifications (
    id,
    email,
    otp_hash,
    purpose,
    expires_at
) VALUES (?, ?, ?, ?, ?);

-- name: GetOTP :one
SELECT * FROM otp_verifications
WHERE email = ?
AND purpose = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteOTP :exec
-- name: DeleteOTP :exec
DELETE FROM otp_verifications
WHERE id = ?
AND purpose = ?;
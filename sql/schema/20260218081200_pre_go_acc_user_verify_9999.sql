-- name: SetVerifyOTP :execresult
INSERT INTO pre_go_acc_user_verify_9999 (
    verify_otp,
    verify_key,
    verify_key_hash,
    verify_type
) VALUES (?, ?, ?, ?);

-- name: GetVerifyOTP :one
SELECT verify_id, verify_otp, verify_key, verify_key_hash, is_verified, verify_created_at
FROM pre_go_acc_user_verify_9999
WHERE verify_key = ? AND is_deleted = 0
LIMIT 1;

-- name: UpdateVerifyStatus :exec
UPDATE pre_go_acc_user_verify_9999
SET 
    is_verified = 1,
    verify_updated_at = CURRENT_TIMESTAMP()
WHERE verify_key = ?;

-- name: SoftDeleteVerifyKey :exec
UPDATE pre_go_acc_user_verify_9999
SET 
    is_deleted = 1,
    verify_updated_at = CURRENT_TIMESTAMP()
WHERE verify_key = ?;

-- name: GetValidOTP :one
SELECT verify_otp, verify_key_hash
FROM pre_go_acc_user_verify_9999
WHERE verify_key = ? AND is_verified = 0 AND is_deleted = 0
LIMIT 1;
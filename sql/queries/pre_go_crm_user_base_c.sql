-- name: GetUserBaseByAccount :one
SELECT user_id, user_account, user_password, user_salt
FROM pre_go_acc_user_base_9999
WHERE user_account = ?
LIMIT 1;

-- name: GetUserBaseById :one
SELECT user_id, user_account, user_password, user_salt, user_login_ip, user_created_at, user_updated_at
FROM pre_go_acc_user_base_9999
WHERE user_id = ?
LIMIT 1;

-- name: CreateUserBase :execresult
INSERT INTO pre_go_acc_user_base_9999 (
    user_account, 
    user_password, 
    user_salt, 
    user_login_ip
) VALUES (?, ?, ?, ?);

-- name: UpdateUserBaseLoginInfo :exec
UPDATE pre_go_acc_user_base_9999
SET 
    user_login_time = CURRENT_TIMESTAMP(),
    user_login_ip = ?,
    user_updated_at = CURRENT_TIMESTAMP()
WHERE user_id = ?;

-- name: UpdateUserBasePassword :exec
UPDATE pre_go_acc_user_base_9999
SET 
    user_password = ?,
    user_salt = ?,
    user_updated_at = CURRENT_TIMESTAMP()
WHERE user_id = ?;

-- name: CheckUserBaseExists :one
SELECT COUNT(*) 
FROM pre_go_acc_user_base_9999 
WHERE user_account = ?;

-- name: ListUsersBase :many
SELECT user_id, user_account, user_created_at
FROM pre_go_acc_user_base_9999
ORDER BY user_created_at DESC
LIMIT ? OFFSET ?;

-- name: DeleteUserBase :exec
DELETE FROM pre_go_acc_user_base_9999 
WHERE user_id = ?;
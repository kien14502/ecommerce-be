-- name: GetUserByEmail :one
SELECT user_email, user_id 
FROM pre_go_acc_user_9999 
WHERE user_email = ? 
LIMIT 1;

-- name: UpdateUserStatusByUserId :exec
UPDATE pre_go_acc_user_9999 
SET user_state = ?
WHERE user_id = ?;

-- name: CreateUser :execresult
INSERT INTO pre_go_acc_user_9999 (
    user_email, 
    user_mobile, 
    user_account, 
    created_at,
    user_state
) VALUES (?, ?, ?, ?, ?);

-- name: GetUserById :one
SELECT * 
FROM pre_go_acc_user_9999 
WHERE user_id = ?;

-- name: DeleteUser :exec
UPDATE pre_go_acc_user_9999 
SET user_state = -1 
WHERE user_id = ?;

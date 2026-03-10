-- name: CreateOAuthAccount :exec
INSERT INTO oauth_accounts (
    id,
    user_id,
    provider,
    provider_user_id
) VALUES (?, ?, ?, ?);

-- name: GetOAuthAccount :one
SELECT * FROM oauth_accounts
WHERE provider = $1
AND provider_user_id = $2;
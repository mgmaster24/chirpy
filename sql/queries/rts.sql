-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: ResetRefreshTokens :execresult
DELETE FROM refresh_tokens;

-- name: RevokeToken :execresult
UPDATE refresh_tokens
SET updated_at=NOW(), revoked_at=NOW()
WHERE refresh_tokens.token = $1;


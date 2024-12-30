-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  uuid_generate_v4(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserEmailById :one
SELECT * FROM users WHERE id = $1;

-- name: ResetUsers :execresult
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: UpdateUserEmailPass :one
UPDATE users
SET email=$2, hashed_password=$3, updated_at=NOW()
WHERE id=$1
RETURNING *;

-- name: MakeChirpyRed :one
UPDATE users
SET is_chirpy_red=$2
WHERE id=$1
RETURNING *;

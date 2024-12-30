-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
  uuid_generate_v4(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetGhirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetChirpsForUserById :many
SELECT * FROM chirps WHERE user_id = $1;

-- name: GetChirpsForUserByEmail :many
SELECT chirps.* FROM chirps
INNER JOIN users ON chirps.user_id = users.id
WHERE users.email = $1;

-- name: ResetChirps :execresult
DELETE FROM chirps;

-- name: DeleteChirp :execresult
DELETE FROM chirps
WHERE user_id = $1 AND id = $2;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllUsers :many
SELECT email FROM users;

-- name: Gethash :one
SELECT hashed_password FROM users
WHERE email = $1;

-- name: GetId :one
SELECT id from users
WHERE email = $1;

-- name: UpgradeToRedByID :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
    AND is_chirpy_red = FALSE;

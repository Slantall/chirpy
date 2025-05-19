-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;


-- name: DeleteAllUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;


-- name: UpdateEmailAndPass :one
UPDATE users
SET updated_at = NOW(), email = $1, hashed_password = $2
WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red;


-- name: UpdateChirpyRed :one
UPDATE users
SET updated_at = NOW(), is_chirpy_red = $1
WHERE id = $2
RETURNING id, created_at, updated_at, email, is_chirpy_red;
-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetUsers :many
SELECT 
    user_id, 
    email, 
    first_name, 
    last_name, 
    created_at 
FROM users 
ORDER BY created_at DESC 
LIMIT $1 
OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

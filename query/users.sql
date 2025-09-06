-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE user_id = $1;

-- name: CreateUser :one
INSERT INTO users (
  username,
  name,
  password
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;
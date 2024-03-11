-- name: CreateProduct :one
INSERT INTO products (
  name,
  price,
  image_url,
  stock,
  condition,
  tags,
  is_purchaseable
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

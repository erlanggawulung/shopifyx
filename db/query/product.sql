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

-- name: UpdateProduct :one
UPDATE products
SET
  name = $2,
  price = $3,
  image_url = $4,
  stock = $5,
  condition = $6,
  tags = $7,
  is_purchaseable = $8
WHERE
  id = $1
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE
  id = $1
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

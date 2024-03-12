-- name: CreateProduct :one
INSERT INTO products (
  name,
  price,
  image_url,
  stock,
  condition,
  tags,
  is_purchaseable,
  created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
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

-- name: UpdateProductStock :one
UPDATE products
SET
  stock = $2
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

-- name: ListProducts :one
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;
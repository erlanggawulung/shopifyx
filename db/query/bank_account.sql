-- name: CreateBankAccount :one
INSERT INTO bank_accounts (
  user_id,
  bank_name,
  bank_account_name,
  bank_account_number
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetBankAccountsByUserId :many
SELECT id, bank_name, bank_account_name, bank_account_number FROM bank_accounts
WHERE user_id = $1;

-- name: DeleteBankAccount :one
DELETE FROM bank_accounts
WHERE
  id = $1
  AND 
  user_id = $2
RETURNING *;

-- name: UpdateBankAccount :one
UPDATE bank_accounts
SET
  bank_name = $3,
  bank_account_name = $4,
  bank_account_number = $5
WHERE
  id = $1
  AND 
  user_id = $2
RETURNING *;
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
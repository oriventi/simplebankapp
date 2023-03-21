-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    $1, $2 ,$3
) RETURNING *;

-- name: LockAccountEntry :one
SELECT * FROM accounts 
WHERE id = $1 LIMIT 1 
FOR NO KEY UPDATE;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id ASC
LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $1
WHERE id = $2
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :execresult
DELETE FROM accounts WHERE id = $1;
-- name: GetTransfer :one
SELECT * FROM transfers 
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount,
    created_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: DeleteTransfer :execresult
DELETE FROM transfers
WHERE id = $1;
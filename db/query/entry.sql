-- name: GetEntry :one
SELECT * FROM entries 
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id 
LIMIT $1 OFFSET $2;

-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount,
    created_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: DeleteEntry :execresult
DELETE FROM entries
WHERE id = $1;


-- name: CreateRefund :one
WITH inserted_refund AS (
    INSERT INTO payment.refund (
        payment_id,
        method,
        status,
        reason,
        address
    )
    VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING *
),
inserted_resources AS (
    INSERT INTO product.resource (owner_id, s3_id)
    SELECT id, unnest(sqlc.arg('resources')::text[]) FROM inserted_refund
    RETURNING s3_id
)
SELECT r.id, COALESCE(array_agg(res.s3_id), '{}')::text[] as resources
FROM inserted_refund r
LEFT JOIN inserted_resources res ON true
GROUP BY r.id;

-- name: UpdateRefund :exec
UPDATE payment.refund
SET 
    method = COALESCE(sqlc.narg('method'), method),
    status = COALESCE(sqlc.narg('status'), status),
    reason = COALESCE(sqlc.narg('reason'), reason),
    address = CASE 
                 WHEN sqlc.arg('null_address')::bool THEN NULL 
                 ELSE COALESCE(sqlc.narg('address'), address) 
              END
WHERE id = $1;

-- name: DeleteRefund :exec
DELETE FROM payment.refund WHERE id = $1;
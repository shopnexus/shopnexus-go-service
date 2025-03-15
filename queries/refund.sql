-- name: ExistsRefund :one
SELECT EXISTS (
  SELECT 1
  FROM payment.refund r
  INNER JOIN payment.base p ON r.payment_id = p.id
  WHERE (
    r.id = $1 AND 
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
  )
) AS exists;

-- name: GetRefund :one
SELECT 
  r.*,
  COALESCE(array_agg(res.s3_id), '{}')::text[] AS resources
FROM payment.refund r
LEFT JOIN product.resource res ON r.id = res.owner_id
WHERE (
  r.id = $1 AND (
    sqlc.narg('user_id') IS NULL OR r.user_id = sqlc.narg('user_id')
  )
)
GROUP BY r.id;

-- name: CountRefunds :one
SELECT COUNT(id)
FROM payment.refund r
INNER JOIN payment.base p ON r.payment_id = p.id
WHERE (
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) AND
    (r.payment_id = sqlc.narg('payment_id') OR sqlc.narg('payment_id') IS NULL) AND
    (r.method = sqlc.narg('method') OR sqlc.narg('method') IS NULL) AND
    (r.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
    (r.reason ILIKE '%' || sqlc.narg('reason') || '%' OR sqlc.narg('reason') IS NULL) AND
    (r.address ILIKE '%' || sqlc.narg('address') || '%' OR sqlc.narg('address') IS NULL) AND
    (r.date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
    (r.date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
);

-- name: ListRefunds :many
SELECT 
    r.*,
    COALESCE(array_agg(res.s3_id), '{}')::text[] as resources
FROM payment.refund r
LEFT JOIN product.resource res ON res.owner_id = r.id
INNER JOIN payment.base p ON r.payment_id = p.id
WHERE (
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) AND
    (r.payment_id = sqlc.narg('payment_id') OR sqlc.narg('payment_id') IS NULL) AND
    (r.method = sqlc.narg('method') OR sqlc.narg('method') IS NULL) AND
    (r.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
    (r.reason ILIKE '%' || sqlc.narg('reason') || '%' OR sqlc.narg('reason') IS NULL) AND
    (r.address ILIKE '%' || sqlc.narg('address') || '%' OR sqlc.narg('address') IS NULL) AND
    (r.date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
    (r.date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
)
ORDER BY r.date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

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
    address = COALESCE(sqlc.narg('address'), address)
WHERE id = $1;

-- name: DeleteRefund :exec


DELETE FROM payment.refund WHERE id = $1;
-- name: ExistsPayment :one
SELECT EXISTS (
  SELECT 1
  FROM payment.base p
  WHERE (
    p.id = $1 AND 
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
  )
) AS exists;

-- name: GetPayment :one
SELECT p.*
FROM payment.base p
WHERE (
  p.id = $1 AND 
  (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
);

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

-- name: GetPaymentProducts :many
SELECT pop.*
FROM payment.product_on_payment pop
WHERE pop.payment_id = $1;

-- name: CreatePayment :one
INSERT INTO payment.base (
    user_id,
    method,
    status,
    address,
    total
)
VALUES (
    $1, $2, $3, $4, $5
) 
RETURNING *;

-- name: CreatePaymentProducts :copyfrom
INSERT INTO "payment".product_on_payment (
    payment_id,
    product_serial_id,
    quantity,
    price,
    total_price
)
VALUES (
    $1, $2, $3, $4, $5
);

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
-- name: ExistsPayment :one
SELECT EXISTS (
  SELECT 1
  FROM payment.base p
  WHERE (
    p.id = $1 AND 
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) AND 
    (p.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
  )
) AS exists;

-- name: GetPayment :one
SELECT p.*
FROM payment.base p
WHERE (
  p.id = $1 AND 
  (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
);

-- name: CountPayments :one
SELECT COUNT(p.id)
FROM payment.base p
WHERE (
  (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) AND
  (p.method = sqlc.narg('method') OR sqlc.narg('method') IS NULL) AND
  (p.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
  (p.address ILIKE '%' || sqlc.narg('address') || '%' OR sqlc.narg('address') IS NULL) AND
  (p.total >= sqlc.narg('total_from') OR sqlc.narg('total_from') IS NULL) AND
  (p.total <= sqlc.narg('total_to') OR sqlc.narg('total_to') IS NULL) AND
  (p.date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
  (p.date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
);

-- name: ListPayments :many
SELECT p.*
FROM payment.base p
WHERE (
  (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) AND
  (p.method = sqlc.narg('method') OR sqlc.narg('method') IS NULL) AND
  (p.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
  (p.address ILIKE '%' || sqlc.narg('address') || '%' OR sqlc.narg('address') IS NULL) AND
  (p.total >= sqlc.narg('total_from') OR sqlc.narg('total_from') IS NULL) AND
  (p.total <= sqlc.narg('total_to') OR sqlc.narg('total_to') IS NULL) AND
  (p.date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
  (p.date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
)
ORDER BY p.date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: GetPaymentProducts :many
SELECT pop.*, pm.id as product_model_id
FROM payment.product_on_payment pop
INNER JOIN product.base p ON pop.product_serial_id = p.serial_id
INNER JOIN product.model pm ON p.product_model_id = pm.id
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

-- name: UpdatePayment :exec
UPDATE payment.base
SET 
    method = COALESCE(sqlc.narg('method'), method),
    status = COALESCE(sqlc.narg('status'), status),
    address = COALESCE(sqlc.narg('address'), address),
    total = COALESCE(sqlc.narg('total'), total)
WHERE id = $1;

-- name: DeletePayment :exec
DELETE FROM payment.base WHERE id = $1;

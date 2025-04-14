-- name: ExistsRefund :one
SELECT EXISTS (
  SELECT 1
  FROM payment.refund r
  INNER JOIN payment.product_on_payment pop ON r.product_on_payment_id = pop.id
  INNER JOIN payment.base p ON pop.payment_id = p.id
  WHERE (
    (r.product_on_payment_id = sqlc.arg('product_on_payment_id')) AND
    (r.status = 'PENDING' OR r.status = 'SUCCESS') AND
    (p.user_id = sqlc.arg('user_id'))
  )
) AS exists;

-- name: GetRefund :one
SELECT 
  r.*,
  COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM payment.refund r
LEFT JOIN product.resource res ON r.id = res.owner_id
INNER JOIN payment.product_on_payment pop ON r.product_on_payment_id = pop.id
INNER JOIN payment.base p ON pop.payment_id = p.id
WHERE (
  r.id = $1 AND 
  (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
)
GROUP BY r.id;

-- name: CountRefunds :one
SELECT COUNT(r.id)
FROM payment.refund r
INNER JOIN payment.product_on_payment pop ON r.product_on_payment_id = pop.id
INNER JOIN payment.base p ON pop.payment_id = p.id
WHERE (
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) AND
    (r.product_on_payment_id = sqlc.narg('product_on_payment_id') OR sqlc.narg('product_on_payment_id') IS NULL) AND
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
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM payment.refund r
LEFT JOIN product.resource res ON res.owner_id = r.id
INNER JOIN payment.product_on_payment pop ON r.product_on_payment_id = pop.id
INNER JOIN payment.base p ON pop.payment_id = p.id
WHERE (
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) AND
    (r.product_on_payment_id = sqlc.narg('product_on_payment_id') OR sqlc.narg('product_on_payment_id') IS NULL) AND
    (r.method = sqlc.narg('method') OR sqlc.narg('method') IS NULL) AND
    (r.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
    (r.reason ILIKE '%' || sqlc.narg('reason') || '%' OR sqlc.narg('reason') IS NULL) AND
    (r.address ILIKE '%' || sqlc.narg('address') || '%' OR sqlc.narg('address') IS NULL) AND
    (r.date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
    (r.date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
)
GROUP BY r.id
ORDER BY r.date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateRefund :one
WITH inserted_refund AS (
    INSERT INTO payment.refund (
        product_on_payment_id,
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
    INSERT INTO product.resource (owner_id, url)
    SELECT id, unnest(sqlc.arg('resources')::text[]) FROM inserted_refund
    RETURNING url
)
SELECT r.id, COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM inserted_refund r
LEFT JOIN inserted_resources res ON true
GROUP BY r.id;

-- name: UpdateRefund :exec
UPDATE payment.refund r
SET 
    method = COALESCE(sqlc.narg('method'), method),
    status = COALESCE(sqlc.narg('status'), status),
    reason = COALESCE(sqlc.narg('reason'), reason),
    address = COALESCE(sqlc.narg('address'), address)
FROM payment.refund
INNER JOIN payment.product_on_payment pop ON r.product_on_payment_id = pop.id
INNER JOIN payment.base p ON pop.payment_id = p.id
WHERE (
  r.id = $1 AND
  (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
);

-- name: DeleteRefund :exec
DELETE FROM payment.refund r
WHERE r.id = $1
  AND EXISTS (
    SELECT 1
    FROM payment.product_on_payment pop
    JOIN payment.base p ON pop.payment_id = p.id
    WHERE r.product_on_payment_id = pop.id
      AND (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL)
  );


-- name: CanRefund :one
SELECT EXISTS (
  SELECT 1
  FROM payment.product_on_payment pop
  INNER JOIN payment.base p ON pop.payment_id = p.id
  LEFT JOIN payment.refund r ON pop.id = r.product_on_payment_id
  WHERE (
    pop.id = $1 AND
    p.status = 'SUCCESS' AND -- Refund only available for successful payment
    (r.id IS NULL OR r.status = 'FAILED' OR r.status = 'CANCELED') AND -- Refund must not exist or is failed/canceled (not pending/success)
    (p.user_id = sqlc.narg('user_id') OR sqlc.narg('user_id') IS NULL) -- Refund must belong to the user
  )
) AS can_refund;


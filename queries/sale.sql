-- name: GetSale :one
SELECT * FROM product.sale WHERE id = $1;

-- name: CountSales :one
SELECT COUNT(*) FROM product.sale
WHERE
    (sqlc.narg('tag')::text IS NULL OR tag = sqlc.narg('tag')) AND
    (sqlc.narg('product_model_id')::bigint IS NULL OR product_model_id = sqlc.narg('product_model_id')) AND
    (sqlc.narg('brand_id')::bigint IS NULL OR brand_id = sqlc.narg('brand_id')) AND
    (sqlc.narg('date_started_from')::timestamptz IS NULL OR date_started >= sqlc.narg('date_started_from')) AND
    (sqlc.narg('date_started_to')::timestamptz IS NULL OR date_started <= sqlc.narg('date_started_to')) AND
    (sqlc.narg('date_ended_from')::timestamptz IS NULL OR date_ended >= sqlc.narg('date_ended_from')) AND
    (sqlc.narg('date_ended_to')::timestamptz IS NULL OR date_ended <= sqlc.narg('date_ended_to')) AND
    (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'));

-- name: ListSales :many
SELECT * FROM product.sale
WHERE
    (sqlc.narg('tag')::text IS NULL OR tag = sqlc.narg('tag')) AND
    (sqlc.narg('product_model_id')::bigint IS NULL OR product_model_id = sqlc.narg('product_model_id')) AND
    (sqlc.narg('brand_id')::bigint IS NULL OR brand_id = sqlc.narg('brand_id')) AND
    (sqlc.narg('date_started_from')::timestamptz IS NULL OR date_started >= sqlc.narg('date_started_from')) AND
    (sqlc.narg('date_started_to')::timestamptz IS NULL OR date_started <= sqlc.narg('date_started_to')) AND
    (sqlc.narg('date_ended_from')::timestamptz IS NULL OR date_ended >= sqlc.narg('date_ended_from')) AND
    (sqlc.narg('date_ended_to')::timestamptz IS NULL OR date_ended <= sqlc.narg('date_ended_to')) AND
    (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
ORDER BY id
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateSale :one
INSERT INTO product.sale (
    tag,
    product_model_id,
    brand_id,
    date_started,
    date_ended,
    quantity,
    used,
    is_active,
    discount_percent,
    discount_price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateSale :exec
UPDATE product.sale
SET
    tag = COALESCE(sqlc.narg('tag'), tag),
    product_model_id = COALESCE(sqlc.narg('product_model_id'), product_model_id),
    brand_id = COALESCE(sqlc.narg('brand_id'), brand_id),
    date_started = COALESCE(sqlc.narg('date_started'), date_started),
    date_ended = COALESCE(sqlc.narg('date_ended'), date_ended),
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    used = COALESCE(sqlc.narg('used'), used),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    discount_percent = COALESCE(sqlc.narg('discount_percent'), discount_percent),
    discount_price = COALESCE(sqlc.narg('discount_price'), discount_price)
WHERE id = $1;

-- name: DeleteSale :exec
DELETE FROM product.sale WHERE id = $1;

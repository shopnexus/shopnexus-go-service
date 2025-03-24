-- name: GetAvailableProducts :many
SELECT *
FROM product.base
WHERE (
    product_model_id = $1 AND
    sold + sqlc.arg('amount') <= quantity
)
LIMIT sqlc.arg('amount');

-- name: CountProducts :one
SELECT COUNT(id)
FROM product.base
WHERE (
    (product_model_id = sqlc.narg('product_model_id') OR sqlc.narg('product_model_id') IS NULL) AND
    (quantity >= sqlc.narg('quantity_from') OR sqlc.narg('quantity_from') IS NULL) AND
    (quantity <= sqlc.narg('quantity_to') OR sqlc.narg('quantity_to') IS NULL) AND
    (sold >= sqlc.narg('sold_from') OR sqlc.narg('sold_from') IS NULL) AND
    (sold <= sqlc.narg('sold_to') OR sqlc.narg('sold_to') IS NULL) AND
    (add_price >= sqlc.narg('add_price_from') OR sqlc.narg('add_price_from') IS NULL) AND
    (add_price <= sqlc.narg('add_price_to') OR sqlc.narg('add_price_to') IS NULL) AND
    (is_active = sqlc.narg('is_active') OR sqlc.narg('is_active') IS NULL) AND
    (metadata @> sqlc.narg('metadata') OR sqlc.narg('metadata') IS NULL) AND
    (date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
    (date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
);

-- name: GetProduct :one
SELECT *
FROM product.base
WHERE (
    id = sqlc.narg('id') OR 
    serial_id = sqlc.narg('serial_id')
);

-- name: ListProducts :many
SELECT *
FROM product.base
WHERE (
    (product_model_id = sqlc.narg('product_model_id') OR sqlc.narg('product_model_id') IS NULL) AND
    (quantity >= sqlc.narg('quantity_from') OR sqlc.narg('quantity_from') IS NULL) AND
    (quantity <= sqlc.narg('quantity_to') OR sqlc.narg('quantity_to') IS NULL) AND
    (sold >= sqlc.narg('sold_from') OR sqlc.narg('sold_from') IS NULL) AND
    (sold <= sqlc.narg('sold_to') OR sqlc.narg('sold_to') IS NULL) AND
    (add_price >= sqlc.narg('add_price_from') OR sqlc.narg('add_price_from') IS NULL) AND
    (add_price <= sqlc.narg('add_price_to') OR sqlc.narg('add_price_to') IS NULL) AND
    (is_active = sqlc.narg('is_active') OR sqlc.narg('is_active') IS NULL) AND
    (metadata @> sqlc.narg('metadata') OR sqlc.narg('metadata') IS NULL) AND
    (date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
    (date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
)
ORDER BY date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateProduct :one
INSERT INTO product.base (
    serial_id,
    product_model_id,
    quantity,
    sold,
    add_price,
    is_active,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE product.base
SET
    serial_id = COALESCE(sqlc.narg('serial_id'), serial_id),
    product_model_id = COALESCE(sqlc.narg('product_model_id'), product_model_id),
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    sold = COALESCE(sqlc.narg('sold'), sold),
    add_price = COALESCE(sqlc.narg('add_price'), add_price),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM product.base WHERE (
    id = sqlc.narg('id') OR 
    serial_id = sqlc.narg('serial_id')
);
-- name: GetAvailableProducts :many
SELECT *
FROM product.base
WHERE (
    product_model_id = $1 AND
    sold = false
)
LIMIT sqlc.arg('amount');

-- name: CountProducts :one
SELECT COUNT(id)
FROM product.base
WHERE (
    (product_model_id = sqlc.narg('product_model_id') OR sqlc.narg('product_model_id') IS NULL) AND
    (sold = sqlc.narg('sold') OR sqlc.narg('sold') IS NULL) AND
    (date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
    (date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
);

-- name: ListProducts :many
SELECT *
FROM product.base
WHERE (
    (product_model_id = sqlc.narg('product_model_id') OR sqlc.narg('product_model_id') IS NULL) AND
    (sold = sqlc.narg('sold') OR sqlc.narg('sold') IS NULL) AND
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
    sold
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE product.base
SET
    serial_id = COALESCE(sqlc.narg('serial_id'), serial_id),
    product_model_id = COALESCE(sqlc.narg('product_model_id'), product_model_id),
    sold = COALESCE(sqlc.narg('sold'), sold),
    date_updated = NOW()
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM product.base WHERE (
    id = sqlc.narg('id') OR 
    serial_id = sqlc.narg('serial_id')
);
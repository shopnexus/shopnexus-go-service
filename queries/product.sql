-- name: GetProduct :one
WITH filtered_product AS (
    SELECT p.*
    FROM product.base p
    WHERE p.id = sqlc.arg('id')
),
filtered_resource AS (
    SELECT 
        res.owner_id,
        array_agg(res.url ORDER BY res.order ASC) AS resources
    FROM product.resource res
    WHERE res.url IS NOT NULL AND res.owner_id = sqlc.arg('id')
    GROUP BY res.owner_id
)
SELECT 
    p.*,
    COALESCE(r.resources, '{}') AS resources
FROM filtered_product p
LEFT JOIN filtered_resource r ON r.owner_id = p.id;

-- name: CountProducts :one
SELECT COUNT(id)
FROM product.base
WHERE (
    (id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
    (product_model_id = sqlc.narg('product_model_id') OR sqlc.narg('product_model_id') IS NULL) AND
    (quantity >= sqlc.narg('quantity_from') OR sqlc.narg('quantity_from') IS NULL) AND
    (quantity <= sqlc.narg('quantity_to') OR sqlc.narg('quantity_to') IS NULL) AND
    (sold >= sqlc.narg('sold_from') OR sqlc.narg('sold_from') IS NULL) AND
    (sold <= sqlc.narg('sold_to') OR sqlc.narg('sold_to') IS NULL) AND
    (add_price >= sqlc.narg('add_price_from') OR sqlc.narg('add_price_from') IS NULL) AND
    (add_price <= sqlc.narg('add_price_to') OR sqlc.narg('add_price_to') IS NULL) AND
    (is_active = sqlc.narg('is_active') OR sqlc.narg('is_active') IS NULL) AND
    (can_combine = sqlc.narg('can_combine') OR sqlc.narg('can_combine') IS NULL) AND
    (metadata @> sqlc.narg('metadata') OR sqlc.narg('metadata') IS NULL) AND
    (date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
    (date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
);

-- name: ListProducts :many
WITH filtered_product AS (
    SELECT p.*
    FROM product.base p
    WHERE (
        (p.id = sqlc.narg('id') OR sqlc.narg('id') IS NULL) AND
        (p.product_model_id = sqlc.narg('product_model_id') OR sqlc.narg('product_model_id') IS NULL) AND
        (p.quantity >= sqlc.narg('quantity_from') OR sqlc.narg('quantity_from') IS NULL) AND
        (p.quantity <= sqlc.narg('quantity_to') OR sqlc.narg('quantity_to') IS NULL) AND
        (p.sold >= sqlc.narg('sold_from') OR sqlc.narg('sold_from') IS NULL) AND
        (p.sold <= sqlc.narg('sold_to') OR sqlc.narg('sold_to') IS NULL) AND
        (p.add_price >= sqlc.narg('add_price_from') OR sqlc.narg('add_price_from') IS NULL) AND
        (p.add_price <= sqlc.narg('add_price_to') OR sqlc.narg('add_price_to') IS NULL) AND
        (p.is_active = sqlc.narg('is_active') OR sqlc.narg('is_active') IS NULL) AND
        (p.can_combine = sqlc.narg('can_combine') OR sqlc.narg('can_combine') IS NULL) AND
        (p.metadata @> sqlc.narg('metadata') OR sqlc.narg('metadata') IS NULL) AND
        (p.date_created >= sqlc.narg('date_created_from') OR sqlc.narg('date_created_from') IS NULL) AND
        (p.date_created <= sqlc.narg('date_created_to') OR sqlc.narg('date_created_to') IS NULL)
    )
),
filtered_resource AS (
    SELECT 
        res.owner_id,
        array_agg(res.url ORDER BY res.order ASC) AS resources
    FROM product.resource res
    WHERE res.url IS NOT NULL AND res.owner_id IN (SELECT id FROM filtered_product)
    GROUP BY res.owner_id
)
SELECT
    p.*,
    COALESCE(r.resources, '{}') AS resources
FROM filtered_product p
LEFT JOIN filtered_resource r ON r.owner_id = p.id
ORDER BY date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateProduct :one
 INSERT INTO product.base (
    product_model_id,
    quantity,
    sold,
    add_price,
    is_active,  
    can_combine,
    metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE product.base
SET
    product_model_id = COALESCE(sqlc.narg('product_model_id'), product_model_id),
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    sold = COALESCE(sqlc.narg('sold'), sold),
    add_price = COALESCE(sqlc.narg('add_price'), add_price),
    can_combine = COALESCE(sqlc.narg('can_combine'), can_combine),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = $1;

-- name: UpdateProductSold :exec
UPDATE product.base
SET
    sold = sold + sqlc.arg('amount')
WHERE
    (id = ANY(sqlc.arg('ids')::bigint[])) AND 
    (sold + sqlc.arg('amount') <= quantity);

-- name: DeleteProduct :exec
DELETE FROM product.base WHERE id = $1;
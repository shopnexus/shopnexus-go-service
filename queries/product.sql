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
SELECT 
    p.*,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM product.base p
LEFT JOIN product.resource res ON res.owner_id = p.id
WHERE (
    id = sqlc.narg('id') OR 
    serial_id = sqlc.narg('serial_id')
)
GROUP BY p.id;

-- name: ListProducts :many
SELECT
    p.*,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM product.base p
LEFT JOIN product.resource res ON res.owner_id = p.id
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
GROUP BY p.id
ORDER BY date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateProduct :one
WITH inserted_product AS (
    INSERT INTO product.base (
        serial_id,
        product_model_id,
        quantity,
        sold,
        add_price,
        is_active,
        metadata
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING *
),
inserted_resources AS (
    INSERT INTO product.resource (owner_id, url)
    SELECT id, unnest(sqlc.arg('resources')::text[]) FROM inserted_product
    RETURNING url
)
SELECT
    p.id,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM inserted_product p
LEFT JOIN inserted_resources res ON true
GROUP BY p.id;

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

-- name: GetResources :many
SELECT url
FROM product.resource
WHERE owner_id = $1;

-- name: AddResources :exec
INSERT INTO product.resource (owner_id, url)
SELECT $1, unnest(sqlc.arg('resources')::text[])
ON CONFLICT (owner_id, url) DO NOTHING;

-- name: RemoveResources :exec
DELETE FROM product.resource
WHERE owner_id = $1 AND url = ANY(sqlc.arg('resources')::text[]);
-- name: GetProductModel :one
SELECT 
    pm.*,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources,
    COALESCE(array_agg(DISTINCT t.tag) FILTER (WHERE t.tag IS NOT NULL), '{}')::text[] as tags
FROM product.model pm
LEFT JOIN product.resource res ON res.owner_id = pm.id
LEFT JOIN product.tag_on_product_model t ON t.product_model_id = pm.id
WHERE pm.id = $1
GROUP BY pm.id;

-- name: GetProductSerialIDs :many
SELECT serial_id
FROM product.serial
WHERE product_id = $1;

-- name: CountProductTypes :one
SELECT COUNT(id)
FROM product.type
WHERE (
    (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL)
);

-- name: ListProductTypes :many
SELECT t.*
FROM product.type t
WHERE (
    (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL)
)
ORDER BY t.id DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: CountProductModels :one
SELECT COUNT(id)
FROM product.model pm
WHERE (
    (pm.type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
    (pm.brand_id = sqlc.narg('brand_id') OR sqlc.narg('brand_id') IS NULL) AND
    (pm.name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
    (pm.description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL) AND
    (pm.list_price >= sqlc.narg('list_price_from') OR sqlc.narg('list_price_from') IS NULL) AND
    (pm.list_price <= sqlc.narg('list_price_to') OR sqlc.narg('list_price_to') IS NULL) AND
    (pm.date_manufactured >= sqlc.narg('date_manufactured_from') OR sqlc.narg('date_manufactured_from') IS NULL) AND
    (pm.date_manufactured <= sqlc.narg('date_manufactured_to') OR sqlc.narg('date_manufactured_to') IS NULL)
);

-- name: ListProductModels :many
SELECT 
    pm.*,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources,
    COALESCE(array_agg(DISTINCT t.tag) FILTER (WHERE t.tag IS NOT NULL), '{}')::text[] as tags
FROM product.model pm
LEFT JOIN product.resource res ON res.owner_id = pm.id
LEFT JOIN product.tag_on_product_model t ON t.product_model_id = pm.id
WHERE (
    (pm.type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
    (pm.brand_id = sqlc.narg('brand_id') OR sqlc.narg('brand_id') IS NULL) AND
    (pm.name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
    (pm.description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL) AND
    (pm.list_price >= sqlc.narg('list_price_from') OR sqlc.narg('list_price_from') IS NULL) AND
    (pm.list_price <= sqlc.narg('list_price_to') OR sqlc.narg('list_price_to') IS NULL) AND
    (pm.date_manufactured >= sqlc.narg('date_manufactured_from') OR sqlc.narg('date_manufactured_from') IS NULL) AND
    (pm.date_manufactured <= sqlc.narg('date_manufactured_to') OR sqlc.narg('date_manufactured_to') IS NULL)
)
GROUP BY pm.id
ORDER BY pm.id DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateProductModel :one
WITH inserted_model AS (
    INSERT INTO product.model (
        type, brand_id, name, description, list_price, date_manufactured
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    ) RETURNING *
),
inserted_resources AS (
    INSERT INTO product.resource (owner_id, url)
    SELECT id, unnest(sqlc.arg('resources')::text[]) FROM inserted_model
    RETURNING url
),
inserted_tags AS (
    INSERT INTO product.tag_on_product_model (product_model_id, tag)
    SELECT id, unnest(sqlc.arg('tags')::text[]) FROM inserted_model
    RETURNING tag
)
SELECT 
    m.id,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources,
    COALESCE(array_agg(DISTINCT t.tag) FILTER (WHERE t.tag IS NOT NULL), '{}')::text[] as tags
FROM inserted_model m
LEFT JOIN inserted_resources res ON true
LEFT JOIN inserted_tags t ON true
GROUP BY m.id;

-- name: UpdateProductModel :exec
UPDATE product.model
SET 
    type = COALESCE(sqlc.narg('type'), type),
    brand_id = COALESCE(sqlc.narg('brand_id'), brand_id),
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    list_price = COALESCE(sqlc.narg('list_price'), list_price),
    date_manufactured = COALESCE(sqlc.narg('date_manufactured'), date_manufactured)
WHERE id = $1;

-- name: DeleteProductModel :exec
DELETE FROM product.model WHERE id = $1;


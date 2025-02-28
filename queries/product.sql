-- name: GetBrand :one
SELECT 
    b.*,
    COALESCE(array_agg(i.s3_id) FILTER (WHERE i.s3_id IS NOT NULL), '{}')::TEXT[] as resources
FROM product.brand b
LEFT JOIN product.resource i ON i.owner_id = b.id
WHERE b.id = $1
GROUP BY b.id;

-- name: CountBrands :one
WITH filtered_brands AS (
  SELECT b.id
  FROM product.brand b
  WHERE (
    (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
    (description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL)
  )
)
SELECT COUNT(id)
FROM filtered_brands;

-- name: ListBrands :many
WITH filtered_brands AS (
  SELECT
    b.*, 
    COALESCE(array_agg(i.s3_id) FILTER (WHERE i.s3_id IS NOT NULL), '{}')::TEXT[] as resources
  FROM product.brand b
  INNER JOIN product.resource i ON i.owner_id = b.id
  WHERE (
    (name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
    (description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL)
  )
  GROUP BY b.id
)
SELECT *
FROM filtered_brands
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateBrand :one
WITH inserted_brand AS (
    INSERT INTO product.brand (name, description)
    VALUES ($1, $2)
    RETURNING *
),
inserted_resources AS (
    INSERT INTO product.resource (owner_id, s3_id)
    SELECT id, unnest(sqlc.arg('resources')::text[]) FROM inserted_brand
    RETURNING s3_id
)
SELECT 
    b.id,
    COALESCE(array_agg(res.s3_id), '{}')::text[] as resources
FROM inserted_brand b
LEFT JOIN inserted_resources res ON true
GROUP BY b.id;

-- name: UpdateBrand :exec
UPDATE product.brand
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description)
WHERE id = $1;

-- name: DeleteBrand :exec
DELETE FROM product.brand WHERE id = $1;

-- name: GetProductModel :one
SELECT 
    pm.*,
    COALESCE(array_agg(i.s3_id) FILTER (WHERE i.s3_id IS NOT NULL), '{}')::text[] as resources,
    COALESCE(array_agg(t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}')::text[] as tags
FROM product.model pm
LEFT JOIN product.resource i ON i.owner_id = pm.id
LEFT JOIN product.tag_on_product t ON t.product_model_id = pm.id
WHERE pm.id = $1
GROUP BY pm.id;

-- name: CountProductModels :one
WITH filtered_models AS (
    SELECT pm.id
    FROM product.model pm
    WHERE (
        (pm.brand_id = sqlc.narg('brand_id') OR sqlc.narg('brand_id') IS NULL) AND
        (pm.name ILIKE '%' || sqlc.narg('name') || '%' OR sqlc.narg('name') IS NULL) AND
        (pm.description ILIKE '%' || sqlc.narg('description') || '%' OR sqlc.narg('description') IS NULL) AND
        (pm.list_price >= sqlc.narg('list_price_from') OR sqlc.narg('list_price_from') IS NULL) AND
        (pm.list_price <= sqlc.narg('list_price_to') OR sqlc.narg('list_price_to') IS NULL) AND
        (pm.date_manufactured >= sqlc.narg('date_manufactured_from') OR sqlc.narg('date_manufactured_from') IS NULL) AND
        (pm.date_manufactured <= sqlc.narg('date_manufactured_to') OR sqlc.narg('date_manufactured_to') IS NULL)
    )
)
SELECT COUNT(id)
FROM filtered_models;

-- name: ListProductModels :many
SELECT 
    pm.*,
    COALESCE(array_agg(DISTINCT i.s3_id) FILTER (WHERE i.s3_id IS NOT NULL), '{}')::text[] as resources,
    COALESCE(array_agg(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}')::text[] as tags
FROM product.model pm
LEFT JOIN product.resource i ON i.owner_id = pm.id
LEFT JOIN product.tag_on_product t ON t.product_model_id = pm.id
WHERE (
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
        brand_id, name, description, list_price, date_manufactured
    ) VALUES (
        $1, $2, $3, $4, $5
    ) RETURNING *
),
inserted_resources AS (
    INSERT INTO product.resource (owner_id, s3_id)
    SELECT id, unnest(sqlc.arg('resources')::text[]) FROM inserted_model
    RETURNING s3_id
),
inserted_tags AS (
    INSERT INTO product.tag_on_product (product_model_id, tag_name)
    SELECT id, unnest(sqlc.arg('tags')::text[]) FROM inserted_model
    RETURNING tag_name
)
SELECT 
    m.id,
    COALESCE(array_agg(res.s3_id), '{}')::text[] as resources,
    COALESCE(array_agg(t.tag_name), '{}')::text[] as tags
FROM inserted_model m
LEFT JOIN inserted_resources res ON true
LEFT JOIN inserted_tags t ON true
GROUP BY m.id;

-- name: UpdateProductModel :exec
UPDATE product.model
SET 
    brand_id = COALESCE(sqlc.narg('brand_id'), brand_id),
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    list_price = COALESCE(sqlc.narg('list_price'), list_price),
    date_manufactured = COALESCE(sqlc.narg('date_manufactured'), date_manufactured)
WHERE id = $1;

-- name: DeleteProductModel :exec
DELETE FROM product.model WHERE id = $1;

-- name: GetProduct :one
SELECT *
FROM product.base
WHERE (
    id = sqlc.narg('id') OR 
    serial_id = sqlc.narg('serial_id')
);

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

-- name: CreateSale :one
INSERT INTO product.sale (
    tag_name,
    product_model_id,
    date_started,
    date_ended,
    quantity,
    used,
    is_active,
    discount_percent,
    discount_price
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: DeleteSale :exec
DELETE FROM product.sale WHERE id = $1;

-- name: CreateTag :exec
INSERT INTO product.tag (
    tag_name,
    description
) VALUES (
    $1, $2
);

-- name: UpdateTag :exec
UPDATE product.tag
SET 
    description = COALESCE(sqlc.narg('description'), description)
WHERE tag_name = $1;

-- name: DeleteTag :exec
DELETE FROM product.tag WHERE tag_name = $1;

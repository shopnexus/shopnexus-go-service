-- name: GetBrand :one
SELECT 
    b.*,
    COALESCE(array_agg(i.url) FILTER (WHERE i.url IS NOT NULL), '{}')::TEXT[] as images
FROM product.brand b
LEFT JOIN product.image i ON i.brand_id = b.id
WHERE b.id = $1
GROUP BY b.id;

-- name: CountBrands :one
WITH filtered_brands AS (
  SELECT b.id
  FROM product.brand b
  WHERE (
    (name ILIKE sqlc.narg('name') OR sqlc.narg('name') IS NULL) AND
    (description ILIKE sqlc.narg('description') OR sqlc.narg('description') IS NULL)
  )
)
SELECT COUNT(id)
FROM filtered_brands;

-- name: ListBrands :many
WITH filtered_brands AS (
  SELECT
    b.*, 
    COALESCE(array_agg(i.url) FILTER (WHERE i.url IS NOT NULL), '{}')::TEXT[] as images
  FROM product.brand b
  INNER JOIN product.image i ON i.brand_id = b.id
  WHERE (
    (name ILIKE sqlc.narg('name') OR sqlc.narg('name') IS NULL) AND
    (description ILIKE sqlc.narg('description') OR sqlc.narg('description') IS NULL)
  )
  GROUP BY b.id
)
SELECT *
FROM filtered_brands
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateBrand :one
WITH inserted_brand AS (
    INSERT INTO product.brand (id, name, description)
    VALUES ($1, $2, $3)
    RETURNING *
),
inserted_images AS (
    INSERT INTO product.image (brand_id, url)
    SELECT $1, unnest(sqlc.arg('images')::text[])
    RETURNING url
)
SELECT 
    b.*,
    COALESCE(array_agg(i.url), '{}') as images
FROM inserted_brand b
LEFT JOIN inserted_images i ON true
GROUP BY b.id;

-- name: UpdateBrand :exec
UPDATE product.brand
SET
    name = COALESCE($2, name),
    description = COALESCE($3, description)
WHERE id = $1;

-- name: DeleteBrand :exec
DELETE FROM product.brand WHERE id = $1;

-- name: GetProductModel :one
SELECT 
    pm.*,
    COALESCE(array_agg(i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images,
    COALESCE(array_agg(t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') as tags
FROM product.model pm
LEFT JOIN product.image i ON i.product_model_id = pm.id
LEFT JOIN product.tag_on_product t ON t.product_model_id = pm.id
WHERE pm.id = $1
GROUP BY pm.id;


-- name: ListProductModels :many
SELECT 
    pm.*,
    COALESCE(array_agg(DISTINCT i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images,
    COALESCE(array_agg(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') as tags
FROM product.model pm
LEFT JOIN product.image i ON i.product_model_id = pm.id
LEFT JOIN product.tag_on_product t ON t.product_model_id = pm.id
WHERE (
    (pm.brand_id = $1 OR $1 IS NULL) AND
    (pm.name ILIKE '%' || $2 || '%' OR $2 IS NULL) AND
    (pm.description ILIKE '%' || $3 || '%' OR $3 IS NULL) AND
    (pm.list_price = $4 OR $4 IS NULL) AND 
    (pm.date_manufactured >= sqlc.arg('date_manufactured_from') OR sqlc.arg('date_manufactured_from') IS NULL) AND
    (pm.date_manufactured <= sqlc.arg('date_manufactured_to') OR sqlc.arg('date_manufactured_to') IS NULL)
)
GROUP BY pm.id
ORDER BY pm.id DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CreateProductModel :one
WITH inserted_model AS (
    INSERT INTO product.model (
        id, brand_id, name, description, list_price, date_manufactured
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    ) RETURNING *
),
inserted_images AS (
    INSERT INTO product.image (product_model_id, url)
    SELECT $1, unnest($7::text[])
    RETURNING url
),
inserted_tags AS (
    INSERT INTO product.tag_on_product (product_model_id, tag_name)
    SELECT $1, unnest($8::text[])
    RETURNING tag_name
)
SELECT 
    m.id,
    m.brand_id,
    m.name,
    m.description,
    m.list_price,
    COALESCE(array_agg(DISTINCT i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images,
    COALESCE(array_agg(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') as tags
FROM inserted_model m
LEFT JOIN inserted_images i ON true
LEFT JOIN inserted_tags t ON true
GROUP BY m.id;

-- name: UpdateProductModel :one
WITH updated_model AS (
    UPDATE product.model
    SET brand_id = $2,
        name = $3,
        description = $4,
        list_price = $5,
        date_manufactured = $6
    WHERE id = $1
    RETURNING *
),
deleted_images AS (
    DELETE FROM product.image
    WHERE product_model_id = $1
),
deleted_tags AS (
    DELETE FROM product.tag_on_product
    WHERE product_model_id = $1
),
inserted_images AS (
    INSERT INTO product.image (product_model_id, url)
    SELECT $1, unnest($7::text[])
    RETURNING url
),
inserted_tags AS (
    INSERT INTO product.tag_on_product (product_model_id, tag_name)
    SELECT $1, unnest($8::text[])
    RETURNING tag_name
)
SELECT 
    m.id,
    m.brand_id,
    m.name,
    m.description,
    m.list_price,
    COALESCE(array_agg(DISTINCT i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images,
    COALESCE(array_agg(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') as tags
FROM updated_model m
LEFT JOIN inserted_images i ON true
LEFT JOIN inserted_tags t ON true
GROUP BY m.id;

-- name: DeleteProductModel :exec
DELETE FROM product.model WHERE id = $1;

-- name: GetProduct :one
SELECT 
    serial_id,
    product_model_id,
    EXTRACT(EPOCH FROM date_created)::bigint as date_created,
    EXTRACT(EPOCH FROM date_update)::bigint as date_update
FROM product.base
WHERE serial_id = $1;

-- name: ListProducts :many
SELECT 
    serial_id,
    product_model_id,
    EXTRACT(EPOCH FROM date_created)::bigint as date_created,
    EXTRACT(EPOCH FROM date_update)::bigint as date_update
FROM product.base
WHERE ($1::bytea IS NULL OR product_model_id = $1)
    AND ($2::timestamp IS NULL OR date_created >= $2)
    AND ($3::timestamp IS NULL OR date_created <= $3)
ORDER BY date_created DESC
LIMIT $4 OFFSET $5;

-- name: CreateProduct :one
INSERT INTO product.base (
    serial_id,
    product_model_id,
    date_created,
    date_update
) VALUES (
    $1, $2, NOW(), NOW()
) RETURNING 
    serial_id,
    product_model_id,
    EXTRACT(EPOCH FROM date_created)::bigint as date_created,
    EXTRACT(EPOCH FROM date_update)::bigint as date_update;

-- name: UpdateProduct :one
UPDATE product.base
SET 
    product_model_id = $2,
    date_update = NOW()
WHERE serial_id = $1
RETURNING 
    serial_id,
    product_model_id,
    EXTRACT(EPOCH FROM date_created)::bigint as date_created,
    EXTRACT(EPOCH FROM date_update)::bigint as date_update;

-- name: DeleteProduct :exec
DELETE FROM product.base WHERE serial_id = $1;

-- name: CreateSale :one
INSERT INTO product.sale (
    id,
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
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: DeleteSale :exec
DELETE FROM product.sale WHERE id = $1;

-- name: CreateTag :one
INSERT INTO product.tag (
    tag_name,
    description
) VALUES (
    $1, $2
) RETURNING *;

-- name: DeleteTag :exec
DELETE FROM product.tag WHERE tag_name = $1;

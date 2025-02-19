-- name: GetBrand :one
SELECT 
    b.id,
    b.name,
    b.description,
    COALESCE(array_agg(i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images
FROM product.brand b
LEFT JOIN product.image i ON i.brand_id = b.id
WHERE b.id = $1
GROUP BY b.id;

-- name: ListBrands :many
SELECT 
    b.id,
    b.name,
    b.description,
    COALESCE(array_agg(i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images
FROM product.brand b
LEFT JOIN product.image i ON i.brand_id = b.id
WHERE ($1::text IS NULL OR b.name ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR b.description ILIKE '%' || $2 || '%')
GROUP BY b.id
ORDER BY b.name
LIMIT $3 OFFSET $4;

-- name: CreateBrand :one
WITH inserted_brand AS (
    INSERT INTO product.brand (id, name, description)
    VALUES ($1, $2, $3)
    RETURNING *
),
inserted_images AS (
    INSERT INTO product.image (brand_id, url)
    SELECT $1, unnest($4::text[])
    RETURNING url
)
SELECT 
    b.id,
    b.name,
    b.description,
    COALESCE(array_agg(i.url), '{}') as images
FROM inserted_brand b
LEFT JOIN inserted_images i ON true
GROUP BY b.id;

-- name: UpdateBrand :one
WITH updated_brand AS (
    UPDATE product.brand
    SET name = $2,
        description = $3
    WHERE id = $1
    RETURNING *
),
deleted_images AS (
    DELETE FROM product.image
    WHERE brand_id = $1
),
inserted_images AS (
    INSERT INTO product.image (brand_id, url)
    SELECT $1, unnest($4::text[])
    RETURNING url
)
SELECT 
    b.id,
    b.name,
    b.description,
    COALESCE(array_agg(i.url), '{}') as images
FROM updated_brand b
LEFT JOIN inserted_images i ON true
GROUP BY b.id;

-- name: DeleteBrand :exec
DELETE FROM product.brand WHERE id = $1;

-- name: GetProductModel :one
SELECT 
    pm.id,
    pm.brand_id,
    pm.name,
    pm.description,
    pm.list_price,
    COALESCE(array_agg(DISTINCT i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images,
    COALESCE(array_agg(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') as tags
FROM product.model pm
LEFT JOIN product.image i ON i.product_model_id = pm.id
LEFT JOIN product.tag_on_product t ON t.product_model_id = pm.id
WHERE pm.id = $1
GROUP BY pm.id;

-- name: ListProductModels :many
SELECT 
    pm.id,
    pm.brand_id,
    pm.name,
    pm.description,
    pm.list_price,
    COALESCE(array_agg(DISTINCT i.url) FILTER (WHERE i.url IS NOT NULL), '{}') as images,
    COALESCE(array_agg(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') as tags
FROM product.model pm
LEFT JOIN product.image i ON i.product_model_id = pm.id
LEFT JOIN product.tag_on_product t ON t.product_model_id = pm.id
WHERE ($1::bytea IS NULL OR pm.brand_id = $1)
    AND ($2::text IS NULL OR pm.name ILIKE '%' || $2 || '%')
    AND ($3::text IS NULL OR pm.description ILIKE '%' || $3 || '%')
    AND ($4::decimal IS NULL OR pm.list_price = $4)
    AND ($5::timestamp IS NULL OR pm.date_manufactured >= $5)
    AND ($6::timestamp IS NULL OR pm.date_manufactured <= $6)
GROUP BY pm.id
ORDER BY pm.name
LIMIT $7 OFFSET $8;

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

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
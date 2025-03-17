-- name: GetTag :one
SELECT * FROM product.tag WHERE tag = $1;

-- name: CountTags :one
SELECT COUNT(*) FROM product.tag
WHERE
    (sqlc.narg('tag')::text IS NULL OR tag ILIKE '%' || sqlc.narg('tag') || '%') AND
    (sqlc.narg('description')::text IS NULL OR description ILIKE '%' || sqlc.narg('description') || '%');

-- name: ListTags :many
SELECT * FROM product.tag
WHERE
    (sqlc.narg('tag')::text IS NULL OR tag ILIKE '%' || sqlc.narg('tag') || '%') AND
    (sqlc.narg('description')::text IS NULL OR description ILIKE '%' || sqlc.narg('description') || '%')
ORDER BY tag
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateTag :exec
INSERT INTO product.tag (
    tag,
    description
) VALUES (
    $1, $2
);

-- name: UpdateTag :exec
UPDATE product.tag
SET 
    tag = COALESCE(sqlc.narg('new_tag'), tag),
    description = COALESCE(sqlc.narg('description'), description)
WHERE tag = $1;

-- name: DeleteTag :exec
DELETE FROM product.tag WHERE tag = $1;
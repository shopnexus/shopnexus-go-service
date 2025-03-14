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
    description = COALESCE(sqlc.narg('description'), description)
WHERE tag = $1;

-- name: DeleteTag :exec
DELETE FROM product.tag WHERE tag = $1;

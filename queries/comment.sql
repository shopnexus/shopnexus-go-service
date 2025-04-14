-- name: GetComment :one
SELECT 
    c.*,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM product.comment c
LEFT JOIN product.resource res ON c.id = res.owner_id
WHERE id = $1
GROUP BY c.id;

-- name: CountComments :one
SELECT COUNT(id) FROM product.comment
WHERE
    (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
    (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
    (dest_id = sqlc.narg('dest_id') OR sqlc.narg('dest_id') IS NULL) AND
    (sqlc.narg('body') ILIKE '%' || sqlc.narg('body') || '%' OR sqlc.narg('body') IS NULL) AND
    (upvote >= sqlc.narg('upvote_from') OR sqlc.narg('upvote_from') IS NULL) AND
    (upvote <= sqlc.narg('upvote_to') OR sqlc.narg('upvote_to') IS NULL) AND
    (downvote >= sqlc.narg('downvote_from') OR sqlc.narg('downvote_from') IS NULL) AND
    (downvote <= sqlc.narg('downvote_to') OR sqlc.narg('downvote_to') IS NULL) AND
    (score >= sqlc.narg('score_from') OR sqlc.narg('score_from') IS NULL) AND
    (score <= sqlc.narg('score_to') OR sqlc.narg('score_to') IS NULL) AND
    (date_created >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
    (date_created <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL);

-- name: ListComments :many
SELECT 
    c.*,
    COALESCE(array_agg(DISTINCT res.url) FILTER (WHERE res.url IS NOT NULL), '{}')::text[] as resources
FROM product.comment c
LEFT JOIN product.resource res ON c.id = res.owner_id
WHERE
    (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL) AND
    (type = sqlc.narg('type') OR sqlc.narg('type') IS NULL) AND
    (dest_id = sqlc.narg('dest_id') OR sqlc.narg('dest_id') IS NULL) AND
    (sqlc.narg('body') ILIKE '%' || sqlc.narg('body') || '%' OR sqlc.narg('body') IS NULL) AND
    (upvote >= sqlc.narg('upvote_from') OR sqlc.narg('upvote_from') IS NULL) AND
    (upvote <= sqlc.narg('upvote_to') OR sqlc.narg('upvote_to') IS NULL) AND
    (downvote >= sqlc.narg('downvote_from') OR sqlc.narg('downvote_from') IS NULL) AND
    (downvote <= sqlc.narg('downvote_to') OR sqlc.narg('downvote_to') IS NULL) AND
    (score >= sqlc.narg('score_from') OR sqlc.narg('score_from') IS NULL) AND
    (score <= sqlc.narg('score_to') OR sqlc.narg('score_to') IS NULL) AND
    (date_created >= sqlc.narg('created_at_from') OR sqlc.narg('created_at_from') IS NULL) AND
    (date_created <= sqlc.narg('created_at_to') OR sqlc.narg('created_at_to') IS NULL)
GROUP BY c.id
ORDER BY date_created DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateComment :exec
INSERT INTO product.comment (
    account_id, type, dest_id, body, upvote, downvote, score
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: UpdateComment :exec
UPDATE product.comment
SET 
    body = COALESCE(sqlc.narg('body'), body),
    upvote = COALESCE(sqlc.narg('upvote'), upvote),
    downvote = COALESCE(sqlc.narg('downvote'), downvote),
    score = COALESCE(sqlc.narg('score'), score)
WHERE 
    id = $1 AND 
    (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL);


-- name: DeleteComment :exec
DELETE FROM product.comment 
WHERE (
  id = $1
  AND (account_id = sqlc.narg('account_id') OR sqlc.narg('account_id') IS NULL)
);
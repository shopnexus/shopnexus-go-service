-- name: GetAccountBase :one
SELECT * FROM "account".base
WHERE id = $1;

-- name: GetAccountAdmin :one
SELECT a.*, b.*
FROM "account".admin a
INNER JOIN "account".base b ON a.id = b.id
WHERE (
  a.id = sqlc.narg('id') OR
  b.username = sqlc.narg('username')
);

-- name: GetAccountUser :one
SELECT u.*, b.*
FROM "account".user u
INNER JOIN "account".base b ON u.id = b.id
WHERE (
  u.id = sqlc.narg('id') OR
  u.email = sqlc.narg('email') OR
  u.phone = sqlc.narg('phone') OR
  b.username = sqlc.narg('username')
);

-- name: CreateAccountUser :one
WITH base AS (
  INSERT INTO "account".base (username, password, role)
  VALUES ($1, $2, 'USER')
  RETURNING id
)
INSERT INTO "account".user (id, email, phone, gender, full_name)
SELECT id, $3, $4, $5, $6
FROM base
RETURNING id;

-- name: CreateAccountAdmin :one
WITH base AS (
  INSERT INTO "account".base (username, password, role)
  VALUES ($1, $2, 'ADMIN')
  RETURNING id
)
INSERT INTO "account".admin (id)
SELECT id
FROM base
RETURNING id;

-- name: AddCartItem :one
INSERT INTO "account".item_on_cart (cart_id, product_model_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_model_id)
DO UPDATE SET quantity = "account".item_on_cart.quantity + $3
RETURNING quantity;

-- name: UpdateCartItem :one
UPDATE "account".item_on_cart
SET quantity = $3
WHERE cart_id = $1 AND product_model_id = $2
RETURNING quantity;

-- name: RemoveCartItem :exec
DELETE FROM "account".item_on_cart
WHERE cart_id = $1 AND product_model_id = $2;

-- name: GetCartItems :many
SELECT * FROM "account".item_on_cart
WHERE cart_id = $1;

-- name: GetCart :one
SELECT * FROM "account".cart
WHERE id = $1;

-- name: CreateCart :exec
INSERT INTO "account".cart (id)
VALUES ($1);

-- name: ClearCart :exec
DELETE FROM "account".item_on_cart
WHERE cart_id = $1;
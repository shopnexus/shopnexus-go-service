-- name: AddCartItem :one
INSERT INTO "account".item_on_cart (cart_id, product_model_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_model_id)
DO UPDATE SET quantity = "account".item_on_cart.quantity + $3
RETURNING quantity;

-- name: DeductCartItem :one
UPDATE "account".item_on_cart
SET quantity = quantity - $3
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

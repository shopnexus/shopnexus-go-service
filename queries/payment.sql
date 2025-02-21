-- name: CreatePayment :one
INSERT INTO payment.base (
    user_id,
    address,
    payment_method,
    total,
    status,
    date_created
)
VALUES (
    $1, $2, $3, $4, $5, NOW()
) 
RETURNING *;

-- name: CreatePaymentProducts :copyfrom
INSERT INTO "payment".product_on_payment (
    payment_id,
    product_serial_id,
    quantity,
    price,
    total_price
)
VALUES (
    $1, $2, $3, $4, $5
);

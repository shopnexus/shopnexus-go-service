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

-- name: GetAccountStaff :one
SELECT s.*, b.*
FROM "account".staff s
INNER JOIN "account".base b ON s.id = b.id
WHERE (
  s.id = sqlc.narg('id') OR
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

-- name: GetRolePermissions :one
SELECT permission FROM "account".permission_on_role
INNER JOIN "account".role ON permission_on_role.role = role.name
WHERE role.name = $1;

-- name: GetCustomPermissions :one
SELECT custom_permission FROM "account".base
WHERE 
id = $1;

-- name: UpdateAccount :one
UPDATE "account".base
SET 
  username = COALESCE(sqlc.narg('username'), username),
  password = COALESCE(sqlc.narg('password'), password),
  custom_permission = CASE WHEN sqlc.narg('null_custom_permission') = TRUE THEN NULL ELSE COALESCE(sqlc.narg('custom_permission'), custom_permission) END
WHERE id = $1
RETURNING *;

-- name: UpdateAccountUser :one
UPDATE "account".user
SET 
  email = COALESCE(sqlc.narg('email'), email),
  phone = COALESCE(sqlc.narg('phone'), phone),
  gender = COALESCE(sqlc.narg('gender'), gender),
  full_name = COALESCE(sqlc.narg('full_name'), full_name),
  default_address_id = CASE WHEN sqlc.narg('null_default_address_id') = TRUE THEN NULL ELSE COALESCE(sqlc.narg('default_address_id'), default_address_id) END
WHERE id = $1
RETURNING *;
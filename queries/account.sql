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
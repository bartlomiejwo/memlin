-- USERS & LINKED ACCOUNTS --

-- name: CreateUser :one
INSERT INTO users (email, name, avatar, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByProviderID :one
SELECT * FROM users u
JOIN linked_accounts la ON u.id = la.user_id
WHERE la.provider = $1 AND la.provider_id = $2;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserLogin :exec
UPDATE users
SET last_login = $1
WHERE email = $2;

-- name: CreateLinkedAccount :one
INSERT INTO linked_accounts (user_id, provider, provider_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetLinkedAccount :one
SELECT * FROM linked_accounts WHERE provider = $1 AND provider_id = $2;

-- name: GetLinkedAccountsByEmail :many
SELECT la.id, la.user_id, la.provider, la.provider_id, la.created_at
FROM linked_accounts la
JOIN users u ON u.id = la.user_id
WHERE u.email = $1;


-- GROUPS, PERMISSIONS & RELATIONSHIPS --

-- name: GetUserGroups :many
SELECT g.id, g.name
FROM groups g
JOIN user_groups ug ON g.id = ug.group_id
WHERE ug.user_id = $1;

-- name: GetGroupPermissions :many
SELECT p.id, p.codename, p.name
FROM permissions p
JOIN group_permissions gp ON p.id = gp.permission_id
WHERE gp.group_id = $1;

-- name: GetUserPermissions :many
SELECT DISTINCT p.id, p.codename, p.name
FROM permissions p
JOIN group_permissions gp ON p.id = gp.permission_id
JOIN user_groups ug ON ug.group_id = gp.group_id
WHERE ug.user_id = $1;

-- name: AddUserToGroup :exec
INSERT INTO user_groups (user_id, group_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: CreateGroup :one
INSERT INTO groups (name)
VALUES ($1)
RETURNING *;

-- name: CreatePermission :one
INSERT INTO permissions (codename, name)
VALUES ($1, $2)
RETURNING *;

-- name: AssignPermissionToGroup :exec
INSERT INTO group_permissions (group_id, permission_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

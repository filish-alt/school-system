-- name: GetUserByUsername :one
SELECT u.id, u.tenant_id, u.username, u.password_hash, u.email, u.role_id, u.status, r.name AS role_name
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
WHERE u.username = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (id, tenant_id, username, password_hash, email, role_id, status)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT u.id, u.tenant_id, u.username, u.password_hash, u.email, u.role_id, u.status, r.name AS role_name
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
WHERE u.id = ? LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = ? WHERE id = ?;

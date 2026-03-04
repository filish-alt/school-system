-- name: GetRoleByName :one
SELECT id, name FROM roles WHERE name = ? LIMIT 1;


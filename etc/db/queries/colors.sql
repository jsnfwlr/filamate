-- name: FindColors :many
SELECT * FROM colors;

-- name: DeleteColor :exec
DELETE FROM colors WHERE id = @id;

-- name: CreateColor :one
INSERT INTO colors (label, hex_code, alias) VALUES (@label, @hex_code, @alias) RETURNING *;

-- name: UpdateColor :one
UPDATE colors SET label = @label, hex_code = @hex_code, alias = @alias WHERE id = @id RETURNING *;

-- name: GetColorByID :one
SELECT * FROM colors WHERE id = @id;

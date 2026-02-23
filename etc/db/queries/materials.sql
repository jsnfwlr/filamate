-- name: FindMaterials :many
SELECT * FROM materials ORDER BY label ASC;

-- name: DeleteMaterial :exec
DELETE FROM materials WHERE id = @id;

-- name: CreateMaterial :one
INSERT INTO materials (label, class, description, special) VALUES (@label, @class, @description, @special) RETURNING *;

-- name: UpdateMaterial :one
UPDATE materials SET label = @label, class = @class, description = @description, special = @special WHERE id = @id RETURNING *;

-- name: GetMaterialByID :one
SELECT * FROM materials WHERE id = @id;

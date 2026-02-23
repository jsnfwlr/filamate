-- name: FindSpools :many
SELECT * FROM spools WHERE deleted_at IS NULL;

-- name: DeleteSpool :exec
UPDATE spools SET deleted_at = NOW() WHERE id = $1;

-- name: CreateSpool :one
INSERT INTO spools (location_id, material_id, brand_id, store_id, empty, weight, combined_weight, current_weight, price) VALUES (@location, @material, @brand, @store, @empty, @weight, @combined_weight, @current_weight, @price) RETURNING *;

-- name: UpdateSpool :one
UPDATE spools SET location_id = @location, material_id = @material, brand_id = @brand, store_id = @store, empty = @empty, weight = @weight, combined_weight = @combined_weight, current_weight = @current_weight, price = @price, updated_at = NOW() WHERE id = @id RETURNING *;

-- name: GetSpoolByID :one
SELECT * FROM spools WHERE id = $1;

-- name: GetSpoolColors :many
SELECT colors.* FROM colors JOIN spool_colors ON colors.id = spool_colors.color_id WHERE spool_colors.spool_id = @spool_id;

-- name: ResetSpoolColor :exec
DELETE FROM spool_colors WHERE spool_id = @spool_id;

-- name: AddSpoolColors :exec
WITH color_set AS (
    SELECT jsonb_array_elements(to_jsonb(@color_ids::bigint[]))::bigint as color_values
)
INSERT INTO spool_colors (spool_id, color_id)
SELECT @spool_id, color_values::bigint
FROM color_set;


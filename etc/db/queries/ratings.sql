-- name: FindRatings :many
SELECT ratings.*
FROM ratings
JOIN spools ON spools.id = ratings.spool_id
JOIN brands ON brands.id = spools.brand_id
JOIN materials ON materials.id = spools.material_id
WHERE spools.deleted_at IS NULL
ORDER BY brands.label, materials.label ASC;

-- name: DeleteRating :exec
DELETE FROM ratings WHERE id = @id;

-- name: CreateRating :one
INSERT INTO ratings (rating, spool_id) VALUES (@rating::bigint, @spool_id) RETURNING *;

-- name: UpdateRating :one
UPDATE ratings SET rating = @rating::bigint, updated_at = NOW() WHERE id = @id RETURNING *;

-- name: GetRatingByID :one
SELECT * FROM ratings WHERE id = @id;

-- name: GetRatingsByBrandIDAndMaterialID :many
SELECT ratings.*
FROM ratings
JOIN spools ON spools.id = ratings.spool_id
JOIN brands ON brands.id = spools.brand_id
JOIN materials ON materials.id = spools.material_id
WHERE spools.brand_id = @brand_id
AND spools.material_id = @material_id
AND spools.deleted_at IS NULL;

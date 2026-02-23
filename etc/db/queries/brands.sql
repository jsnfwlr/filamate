-- name: FindBrands :many
SELECT *
    FROM brands
    ORDER BY label ASC;

-- name: DeleteBrand :exec
DELETE FROM brands WHERE id = @id;

-- name: CreateBrand :one
INSERT INTO brands (label, active, store_id) VALUES (@label, @active, @store_id) RETURNING *;

-- name: UpdateBrand :one
UPDATE brands SET label = @label, active = @active, store_id = @store WHERE id = @id RETURNING *;

-- name: GetBrandByID :one
SELECT * FROM brands WHERE id = @id;

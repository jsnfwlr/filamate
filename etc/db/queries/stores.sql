-- name: FindStores :many
SELECT * FROM stores ORDER BY label ASC;

-- name: DeleteStore :exec
DELETE FROM stores WHERE id = @id;

-- name: CreateStore :one
INSERT INTO stores (label, url) VALUES (@label, @url) RETURNING *;

-- name: UpdateStore :one
UPDATE stores SET label = @label, url = @url WHERE id = @id RETURNING *;

-- name: GetStoreByID :one
SELECT * FROM stores WHERE id = @id;

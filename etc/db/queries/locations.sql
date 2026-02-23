-- name: FindLocations :many
SELECT * FROM locations;

-- name: DeleteLocation :exec
DELETE FROM locations WHERE id = @id;

-- name: CreateLocation :one
INSERT INTO locations (label, description, capacity, printable, tally) VALUES (@label, @description, @capacity, @printable, @tally) RETURNING *;

-- name: UpdateLocation :one
UPDATE locations SET label = @label, description = @description, capacity = @capacity, printable = @printable, tally = @tally WHERE id = @id RETURNING *;

-- name: GetLocationByID :one
SELECT * FROM locations WHERE id = @id;

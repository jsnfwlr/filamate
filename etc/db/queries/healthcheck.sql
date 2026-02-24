-- name: CheckMigration :one
SELECT version::integer FROM db_version;

-- name: CheckDemoData :one
SELECT done FROM demo_data;

-- name: SetDemoData :exec
UPDATE demo_data SET done = true WHERE done = false;

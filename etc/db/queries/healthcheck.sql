-- name: CheckMigration :one
SELECT version::integer FROM db_version;

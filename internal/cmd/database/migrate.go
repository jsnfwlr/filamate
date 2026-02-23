package database

import (
	"context"
	"fmt"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/etc/db/migrations"
	"github.com/jsnfwlr/filamate/internal/db"

	"github.com/spf13/cobra"
)

func init() {
	BaseCmd.AddCommand(MigrateCmd)

	MigrateCmd.Flags().Int32("version", -1, "what version to migrate to, -1 for latest")
}

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "run migrations",
	RunE:  MigrateRun,
}

func MigrateRun(cmd *cobra.Command, args []string) (fault error) {
	ctx := cmd.Context()

	cfg, err := db.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var version int32 = -1
	if cmd.Flags().Changed("version") {
		version, _ = cmd.Flags().GetInt32("version")
	}

	err = DoMigration(ctx, cfg, version)
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}

func DoMigration(ctx context.Context, cfg db.ConfigProvider, version int32) (fault error) {
	ctx, o := go11y.Get(ctx)

	col, err := migrations.New()
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	m, err := db.NewMigrator(ctx, cfg, col)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if version != -1 {
		o.Info("Migrating to specific version", "version", version)
		if i, err := m.MigrateTo(ctx, version); err != nil {
			return fmt.Errorf("failed to set migration version: %w", err)
		} else {
			o.Info("Migration complete", db.SequenceKey, i.Sequence, db.MigrationFileKey, i.File, db.DirectionKey, i.Direction)
		}

		return nil
	}

	if i, err := m.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	} else {
		o.Info("Migration complete", db.SequenceKey, i.Sequence, db.MigrationFileKey, i.File, db.DirectionKey, i.Direction)
	}

	return nil
}

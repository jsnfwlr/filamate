package db

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	migrate "github.com/jackc/tern/v2/migrate"
)

type DBMigrator struct {
	connection    *pgx.Conn
	migrator      *migrate.Migrator
	configuration ConfigProvider
	steps         int32
}

type ConfigProvider interface {
	GetURI() string
	GetRedactedURI() string
	GetHost() string
	GetPort() string
	GetDatabase() string
	GetVersionTable() string
	GetIncDemo() bool
}

type FilesystemProvider interface {
	ReadDir(name string) ([]fs.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	Open(name string) (fs.File, error)
}

type MigratorError struct {
	Connect    bool
	Create     bool
	Filesystem bool
	URL        string
	Err        error
}

func (e MigratorError) Error() string {
	switch {
	case e.Connect:
		return fmt.Sprintf("could not connect to database %q: %v", e.URL, e.Err)
	case e.Create:
		return fmt.Sprintf("could not create migrator: %v", e.Err)
	case e.Filesystem:
		return fmt.Sprintf("could not load migrations: %v", e.Err)
	default:
		return fmt.Sprintf("unknown migration error: %v", e.Err)
	}
}

func NewMigrator(ctx context.Context, dbParams ConfigProvider, fs FilesystemProvider) (migrator DBMigrator, fault error) {
	conn, err := pgx.Connect(ctx, dbParams.GetURI())
	if err != nil {
		return DBMigrator{}, MigratorError{Connect: true, URL: dbParams.GetRedactedURI(), Err: err}
	}

	opts := &migrate.MigratorOptions{
		DisableTx: false,
	}

	mig, err := migrate.NewMigratorEx(ctx, conn, dbParams.GetVersionTable(), opts)
	if err != nil {
		return DBMigrator{}, MigratorError{Create: true, Err: err}
	}

	if fs == nil {
		return DBMigrator{}, MigratorError{Filesystem: true, Err: fmt.Errorf("nil filesystem provided")}
	}

	err = mig.LoadMigrations(fs)
	if err != nil {
		return DBMigrator{}, MigratorError{Filesystem: true, Err: err}
	}

	return DBMigrator{
		connection:    conn,
		migrator:      mig,
		steps:         int32(len(mig.Migrations)),
		configuration: dbParams,
	}, nil
}

type Info struct {
	Hostname   string
	Port       string
	Database   string
	Migrations MigrationInfo
}

type MigrationInfo struct {
	CurrentVersion int32
	TargetVersion  int32
	Stages         []Stage
	Summary        string
}

type Stage struct {
	Sequence  int32
	File      string
	Direction string
	Migrated  bool
}

func (m DBMigrator) Steps() int32 {
	return m.steps
}

func (m DBMigrator) Info(ctx context.Context, stopAfter int32) (information Info, fault error) {
	var err error

	i := Info{
		Hostname:   m.configuration.GetHost(),
		Port:       m.configuration.GetPort(),
		Database:   m.configuration.GetDatabase(),
		Migrations: MigrationInfo{},
	}

	i.Migrations.CurrentVersion, err = m.migrator.GetCurrentVersion(ctx)
	if err != nil {
		return Info{}, fmt.Errorf("could not get current version: %w", err)
	}

	if stopAfter < 0 {
		stopAfter = m.migrator.Migrations[len(m.migrator.Migrations)-1].Sequence
	}

	for _, mig := range m.migrator.Migrations {
		ind := "  "

		s := Stage{
			Sequence: mig.Sequence,
			File:     mig.Name,
			Migrated: mig.Sequence <= i.Migrations.CurrentVersion,
		}
		i.Migrations.Stages = append(i.Migrations.Stages, s)

		if mig.Sequence == stopAfter {
			ind = "> "
		}

		if mig.Sequence == i.Migrations.CurrentVersion {
			ind = "@ "
		}

		i.Migrations.Summary += fmt.Sprintf("%2s %3d %s\n", ind, mig.Sequence, mig.Name)
	}

	return i, nil
}

func (m *DBMigrator) Migrate(ctx context.Context) (details Stage, fault error) {
	var info Stage
	m.migrator.OnStart = func(sequence int32, name string, direction string, sql string) {
		info = Stage{
			Sequence:  sequence,
			File:      name,
			Direction: direction,
		}
	}

	err := m.migrator.Migrate(ctx)
	if err != nil {
		return Stage{}, fmt.Errorf("could not migrate: %w", err)
	}

	return info, nil
}

func (m *DBMigrator) MigrateTo(ctx context.Context, sequence int32) (details Stage, fault error) {
	var info Stage
	m.migrator.OnStart = func(sequence int32, name string, direction string, _ string) {
		info = Stage{
			Sequence:  sequence,
			File:      name,
			Direction: direction,
		}
	}

	err := m.migrator.MigrateTo(ctx, sequence)
	if err != nil {
		return Stage{}, fmt.Errorf("could not migrate to %d: %w", sequence, err)
	}

	return info, nil
}

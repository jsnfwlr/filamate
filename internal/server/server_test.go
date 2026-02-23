package server_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server"
	"github.com/jsnfwlr/filamate/internal/test/containers"
	. "github.com/jsnfwlr/filamate/internal/types"
)

func TestServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})

	gCfg, err := go11y.LoadConfig()
	if err != nil {
		t.Fatalf("could not load go11y config: %v", err)
	}

	ctx, _, err = go11y.Initialise(ctx, gCfg, os.Stdout)
	if err != nil {
		t.Fatalf("could not initialise go11y: %v", err)
	}

	ctr, dbCfg, err := containers.Postgres(t, ctx, "db_version", "17", PointerOf("server-test"))
	if err != nil {
		t.Fatalf("could not start the Postgres container: %v", err)
	}

	t.Cleanup(
		func() {
			ctr.Cleanup(t)
		},
	)

	dbClient, err := db.Connect(ctx, dbCfg)
	if err != nil {
		t.Fatalf("could not connect to the database: %v", err)
	}

	cfg, err := server.LoadConfig()
	if err != nil {
		t.Fatalf("could not load server config: %v", err)
	}

	srv, err := server.New(ctx, cfg, dbClient)
	if err != nil {
		t.Fatalf("could not create server: %v", err)
	}

	// Test server start and stop
	go func() {
		err = srv.Start(ctx)
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("could not start the server: %v", err)
		}
	}()

	t.Cleanup(
		func() {
			srv.Close(ctx)
			dbClient.Close()
		},
	)
}

package daemon

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/cmd/database"
	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server"
	"github.com/jsnfwlr/filamate/internal/server/log"

	"github.com/spf13/cobra"
)

func init() {
	BaseCmd.AddCommand(StartCmd)
}

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "run the daemon",
	RunE:  StartRun,
}

func StartRun(cmd *cobra.Command, args []string) error {
	ctx, o := go11y.Span(cmd.Context(), tracer, "StartRun", go11y.SpanKindClient)
	defer o.End()

	dbConfig, err := db.LoadConfig()
	if err != nil {
		o.Error("could not load database config", err, go11y.SeverityHighest)
		return err
	}

	err = database.DoMigration(ctx, dbConfig, -1)
	if err != nil {
		o.Error("could not perform database migration", err, go11y.SeverityHighest)
		return err
	}

	dbClient, err := db.Connect(ctx, dbConfig)
	if err != nil {
		o.Error("could not connect to database", err, go11y.SeverityHighest)
		return err
	}

	defer dbClient.Close()

	apiConfig, err := server.LoadConfig()
	if err != nil {
		o.Error("could not load API config", err, go11y.SeverityHighest)
		return err
	}

	srvr, err := server.New(ctx, apiConfig, dbClient)
	if err != nil {
		o.Error("could not instantiate web and API handlers", err, go11y.SeverityHighest)
		return err
	}

	sigChan := make(chan os.Signal, 1)

	go func() {
		o.Info("starting filamate daemon", log.ProcessIDKey, os.Getpid())
		err = srvr.Start(ctx)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				o.Error("server error", err, go11y.SeverityHighest)
			}
			sigChan <- syscall.SIGTERM
		}
		o.Info("web and API server stopped accepting new connections")
	}()

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, server.DefaultShutdownTimeout)
	defer shutdownRelease()

	if err := srvr.Close(shutdownCtx); err != nil {
		o.Error("could not shutdown server", err, go11y.SeverityMedium)
		return err
	}

	return nil
}

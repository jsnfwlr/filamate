// Package cmd provides the command-line interface for Filamate.
package cmd

import (
	"context"
	"io"
	"os"

	"github.com/jsnfwlr/filamate/internal/cmd/daemon"
	"github.com/jsnfwlr/filamate/internal/cmd/database"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Version: "v0.1.0",
	Use:     "filamate",
	Short:   "filamate CLI",
}

func Execute(ctx context.Context, errDst, outDst io.Writer, baseCmd *cobra.Command, args ...string) (fault error) {
	RootCmd.SetContext(ctx)

	if baseCmd != nil {
		RootCmd.AddCommand(baseCmd)
	}

	if len(args) > 0 {
		RootCmd.SetArgs(args)
	}

	if errDst == nil {
		errDst = os.Stderr
	}

	RootCmd.SetErr(errDst)

	if outDst == nil {
		outDst = os.Stdout
	}

	RootCmd.SetOut(outDst)

	return RootCmd.Execute()
}

func init() {
	_ = godotenv.Load("dev.env", ".env")

	cobra.EnableCaseInsensitive = true
	// cobra.EnableCommandSorting = false

	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.Flags().SortFlags = false
	RootCmd.PersistentFlags().SortFlags = false

	RootCmd.AddCommand(daemon.BaseCmd)
	RootCmd.AddCommand(database.BaseCmd)
}

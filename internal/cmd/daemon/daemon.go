// Package daemon provides commands to run Filamate as a daemon.
package daemon

import (
	"github.com/jsnfwlr/go11y"

	"github.com/spf13/cobra"
)

var tracer = go11y.NewTracer("github.com/jsnfwlr/filamate/cmd/daemon")

var BaseCmd = &cobra.Command{
	Use:   "daemon",
	Short: "start filamate as a daemon",
}

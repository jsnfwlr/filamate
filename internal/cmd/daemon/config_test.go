package daemon_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/cmd/daemon"
)

func TestConfigRun(t *testing.T) {
	gCfg, err := go11y.LoadConfig()
	if err != nil {
		t.Fatalf("could not load go11y config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
	})

	ctx, _, err = go11y.Initialise(ctx, gCfg, os.Stdout)
	if err != nil {
		t.Fatalf("could not initialise go11y: %v", err)
	}

	cmd := daemon.ConfigCmd
	cmd.SetContext(ctx)
	daemon.ConfigRun(cmd, []string{})
}

func TestGetConfig(t *testing.T) {
	output := daemon.GetConfig()
	expected := fmt.Sprintf(daemon.CfgFormat, "filament", "REDACTED_SECRET", "localhost", "5432", "filament", "0.0.0.0", "9766", "embedded")

	if output != expected {
		t.Errorf("GetConfig() output does not match expected.\nGot:\n%s\nExpected:\n%s", output, expected)
	}
}

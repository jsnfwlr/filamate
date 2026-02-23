package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server"
	"github.com/jsnfwlr/filamate/internal/server/log"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/spf13/cobra"
)

func init() {
	BaseCmd.AddCommand(CheckCmd)
}

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "check the daemon is running",
	RunE:  CheckRun,
}

func CheckRun(cmd *cobra.Command, args []string) error {
	return Healthcheck(cmd.Context())
}

func Healthcheck(ctx context.Context) error {
	_, o := go11y.Span(ctx, tracer, "healthcheck", go11y.SpanKindClient)
	defer o.End()

	o.Debug("checking filamate daemon")

	client := &http.Client{
		Transport: http.DefaultTransport,
	}

	cfg, err := server.LoadConfig()
	if err != nil {
		o.Error("could not load server config", err, go11y.SeverityHigh)
		return err
	}

	err = go11y.AddLoggingToHTTPClient(client)
	if err != nil {
		o.Error("could not add logging roundtripper", err, go11y.SeverityHigh)
		return err
	}

	host := cfg.Host()
	if host == "0.0.0.0" {
		host = "localhost"
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/healthcheck", host, cfg.Port()), nil)
	if err != nil {
		o.Error("could not create request", err, go11y.SeverityHigh)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		o.Error("could not complete request", err, go11y.SeverityHigh)
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody := oapi.HealthCheck500JSONResponse{}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			o.Error("could not read response body", err, go11y.SeverityHigh)

			os.Exit(1)
		}

		err = json.Unmarshal(b, &respBody)
		if err != nil {
			o.Error("could not unmarshal body", err, go11y.SeverityHigh)
			os.Exit(1)
		}

		o.Error(respBody.Message, errors.New("health check failed"), go11y.SeverityHigh, log.StatusCodeKey, resp.StatusCode)

		os.Exit(1)
	}

	o.Debug("filamate daemon is running", log.StatusCodeKey, resp.StatusCode)
	return nil
}

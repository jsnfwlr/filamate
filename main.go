//go:generate go tool oapi-codegen -config ./etc/openapi/server.yaml ./etc/openapi/spec.yaml
//go:generate go tool sqlc generate -f ./etc/db/sqlc.yaml

package main

import (
	"context"
	"os"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/cmd"

	_ "google.golang.org/genproto/googleapis/api/httpbody"
	_ "google.golang.org/genproto/googleapis/rpc/status"
)

func main() {
	cfg, err := go11y.LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx, o, err := go11y.Initialise(context.Background(), cfg, os.Stdout)
	if err != nil {
		panic(err)
	}

	defer o.Close()

	err = cmd.Execute(ctx, os.Stderr, os.Stdout, nil)
	if err != nil {
		os.Exit(1)
	}
}

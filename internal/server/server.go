// Package server provides the web and API server for Filamate.
package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/handlers"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

const (
	DefaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	router *http.ServeMux
	core   *http.Server
}

var tracer = go11y.NewTracer("github.com/jsnfwlr/filamate/internal/server")

func New(ctx context.Context, cfg Config, dbClient *db.Client) (server Server, fault error) {
	ctx, o := go11y.Get(ctx)

	mux := http.NewServeMux()

	o.Debug("creating base router")

	h, err := handlers.New(ctx, dbClient, cfg.StaticType())
	if err != nil {
		return Server{}, err
	}

	api := oapi.NewStrictHandlerWithOptions(h, nil, oapi.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		ResponseErrorHandlerFunc: responseErrorHandler,
	})

	opts := oapi.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []oapi.MiddlewareFunc{
			recoverPanic,
			go11y.SetRequestID,
			o.LogRequest,
		},
	}

	oh := oapi.HandlerWithOptions(api, opts)

	// oh := oapi.HandlerFromMux(api, r)

	mux.HandleFunc("GET /", h.UI)

	core := &http.Server{
		Addr:                         net.JoinHostPort(cfg.Host(), cfg.Port()),
		DisableGeneralOptionsHandler: true,
		Handler:                      oh,
	}

	s := Server{
		router: mux,
		core:   core,
	}

	return s, nil
}

func (srvr *Server) Start(ctx context.Context) (fault error) {
	_, o := go11y.Get(ctx)

	o.Info("starting web and API server", "Address", srvr.core.Addr)
	err := srvr.core.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			o.Error("could not start Web and API server", err, go11y.SeverityHighest)
			return err
		}
	}
	return nil
}

func (srvr *Server) Close(ctx context.Context) (fault error) {
	_, o := go11y.Get(ctx)
	o.Info("shutting web and API server down")

	err := srvr.core.Close()
	if err != nil {
		o.Error("could not shut web and API server down gracefully", err, go11y.SeverityMedium)
		return err
	}
	return nil
}

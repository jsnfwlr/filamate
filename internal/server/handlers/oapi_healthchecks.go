package handlers

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/healthchecks"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// HealthCheck HealthCheck checks the health of the server (GET /healthcheck)
// (GET /health/check)
func (h Handlers) HealthCheck(ctx context.Context, r oapi.HealthCheckRequestObject) (response oapi.HealthCheckResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.HealthCheck500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := healthchecks.HealthCheck(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.HealthCheck500JSONResponse{}, err
	}

	return resp, nil
}

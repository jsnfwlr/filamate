package handlers

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/spools"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// KillSpool Kill spool
// (DELETE /spool/{id})
func (h Handlers) KillSpool(ctx context.Context, r oapi.KillSpoolRequestObject) (response oapi.KillSpoolResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.KillSpool500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := spools.Kill(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.KillSpool500JSONResponse{}, err
	}

	return resp, nil
}

// CheckSpool CORS preflight for spool by ID
// (OPTIONS /spool/{id})
func (h Handlers) CheckSpool(ctx context.Context, r oapi.CheckSpoolRequestObject) (response oapi.CheckSpoolResponseObject, fault error) {
	return spools.Check(ctx, h.DBClient.Queries, r)
}

// CheckSpools CORS preflight for spools
// (OPTIONS /spools)
func (h Handlers) CheckSpools(ctx context.Context, r oapi.CheckSpoolsRequestObject) (response oapi.CheckSpoolsResponseObject, fault error) {
	return spools.Checks(ctx, h.DBClient.Queries, r)
}

// UpdateSpool Update spool
// (PATCH /spool/{id})
func (h Handlers) UpdateSpool(ctx context.Context, r oapi.UpdateSpoolRequestObject) (response oapi.UpdateSpoolResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.UpdateSpool500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := spools.Update(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.UpdateSpool500JSONResponse{}, err
	}

	return resp, nil
}

// FindSpools Find spools
// (GET /spools)
func (h Handlers) FindSpools(ctx context.Context, r oapi.FindSpoolsRequestObject) (response oapi.FindSpoolsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.FindSpools500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := spools.Find(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.FindSpools500JSONResponse{}, err
	}

	return resp, nil
}

// CreateSpool Create spool
// (POST /spools)
func (h Handlers) CreateSpool(ctx context.Context, r oapi.CreateSpoolRequestObject) (response oapi.CreateSpoolResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.CreateSpool500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := spools.Create(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.CreateSpool500JSONResponse{}, err
	}

	return resp, nil
}

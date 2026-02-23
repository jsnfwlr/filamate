package handlers

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/stores"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// KillStore Kill store
// (DELETE /store/{id})
func (h Handlers) KillStore(ctx context.Context, r oapi.KillStoreRequestObject) (response oapi.KillStoreResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.KillStore500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stores.Kill(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.KillStore500JSONResponse{}, err
	}

	return resp, nil
}

// CheckStore CORS preflight for store by ID
// (OPTIONS /store/{id})
func (h Handlers) CheckStore(ctx context.Context, r oapi.CheckStoreRequestObject) (response oapi.CheckStoreResponseObject, fault error) {
	return stores.Check(ctx, h.DBClient.Queries, r)
}

// CheckStores CORS preflight for store by ID
// (OPTIONS /store/{id})
func (h Handlers) CheckStores(ctx context.Context, r oapi.CheckStoresRequestObject) (response oapi.CheckStoresResponseObject, fault error) {
	return stores.Checks(ctx, h.DBClient.Queries, r)
}

// UpdateStore Update store
// (PATCH /store/{id})
func (h Handlers) UpdateStore(ctx context.Context, r oapi.UpdateStoreRequestObject) (response oapi.UpdateStoreResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.UpdateStore500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stores.Update(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.UpdateStore500JSONResponse{}, err
	}

	return resp, nil
}

// FindStores Find stores
// (GET /stores)
func (h Handlers) FindStores(ctx context.Context, r oapi.FindStoresRequestObject) (response oapi.FindStoresResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.FindStores500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stores.Find(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.FindStores500JSONResponse{}, err
	}

	return resp, nil
}

// CreateStore Create store
// (POST /stores)
func (h Handlers) CreateStore(ctx context.Context, r oapi.CreateStoreRequestObject) (response oapi.CreateStoreResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.CreateStore500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stores.Create(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.CreateStore500JSONResponse{}, err
	}

	return resp, nil
}

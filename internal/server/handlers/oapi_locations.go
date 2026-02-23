package handlers

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/locations"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// KillLocation Kill location
// (DELETE /location/{id})
func (h Handlers) KillLocation(ctx context.Context, r oapi.KillLocationRequestObject) (response oapi.KillLocationResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.KillLocation500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := locations.Kill(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.KillLocation500JSONResponse{}, err
	}

	return resp, nil
}

// CheckLocation CORS preflight for location by ID
// (OPTIONS /location/{id})
func (h Handlers) CheckLocation(ctx context.Context, r oapi.CheckLocationRequestObject) (response oapi.CheckLocationResponseObject, fault error) {
	return locations.Check(ctx, h.DBClient.Queries, r)
}

// CheckLocations CORS preflight for location by ID
// (OPTIONS /location/{id})
func (h Handlers) CheckLocations(ctx context.Context, r oapi.CheckLocationsRequestObject) (response oapi.CheckLocationsResponseObject, fault error) {
	return locations.Checks(ctx, h.DBClient.Queries, r)
}

// UpdateLocation Update location
// (PATCH /location/{id})
func (h Handlers) UpdateLocation(ctx context.Context, r oapi.UpdateLocationRequestObject) (response oapi.UpdateLocationResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.UpdateLocation500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := locations.Update(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.UpdateLocation500JSONResponse{}, err
	}

	return resp, nil
}

// FindLocations Find locations
// (GET /locations)
func (h Handlers) FindLocations(ctx context.Context, r oapi.FindLocationsRequestObject) (response oapi.FindLocationsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.FindLocations500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := locations.Find(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.FindLocations500JSONResponse{}, err
	}

	return resp, nil
}

// CreateLocation Create location
// (POST /locations)
func (h Handlers) CreateLocation(ctx context.Context, r oapi.CreateLocationRequestObject) (response oapi.CreateLocationResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.CreateLocation500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := locations.Create(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.CreateLocation500JSONResponse{}, err
	}

	return resp, nil
}

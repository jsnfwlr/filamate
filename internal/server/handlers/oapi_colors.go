package handlers

import (
	"context"
	"net/http"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/colors"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// KillColor Kill color
// (DELETE /color/{id})
func (h Handlers) KillColor(ctx context.Context, r oapi.KillColorRequestObject) (response oapi.KillColorResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.KillColor500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := colors.Kill(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.KillColor500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// CheckColor CORS preflight for color by ID
// (OPTIONS /color/{id})
func (h Handlers) CheckColor(ctx context.Context, r oapi.CheckColorRequestObject) (response oapi.CheckColorResponseObject, fault error) {
	return colors.Check(ctx, h.DBClient.Queries, r)
}

// CheckColors CORS preflight for color by ID
// (OPTIONS /color/{id})
func (h Handlers) CheckColors(ctx context.Context, r oapi.CheckColorsRequestObject) (response oapi.CheckColorsResponseObject, fault error) {
	return colors.Checks(ctx, h.DBClient.Queries, r)
}

// UpdateColor Update color
// (PATCH /color/{id})
func (h Handlers) UpdateColor(ctx context.Context, r oapi.UpdateColorRequestObject) (response oapi.UpdateColorResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.UpdateColor500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := colors.Update(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.UpdateColor500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// FindColors Find colors
// (GET /colors)
func (h Handlers) FindColors(ctx context.Context, r oapi.FindColorsRequestObject) (response oapi.FindColorsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.FindColors500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := colors.Find(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.FindColors500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// CreateColor Create color
// (POST /colors)
func (h Handlers) CreateColor(ctx context.Context, r oapi.CreateColorRequestObject) (response oapi.CreateColorResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.CreateColor500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := colors.Create(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.CreateColor500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

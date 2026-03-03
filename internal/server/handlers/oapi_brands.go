package handlers

import (
	"context"
	"net/http"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/brands"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// KillBrand Kill brand
// (DELETE /brand/{id})
func (h Handlers) KillBrand(ctx context.Context, r oapi.KillBrandRequestObject) (response oapi.KillBrandResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.KillBrand500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := brands.Kill(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.KillBrand500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// CheckBrand CORS preflight for brand by ID
// (OPTIONS /brand/{id})
func (h Handlers) CheckBrand(ctx context.Context, r oapi.CheckBrandRequestObject) (response oapi.CheckBrandResponseObject, fault error) {
	return brands.Check(ctx, h.DBClient.Queries, r)
}

// CheckBrands CORS preflight for brand by ID
// (OPTIONS /brand/{id})
func (h Handlers) CheckBrands(ctx context.Context, r oapi.CheckBrandsRequestObject) (response oapi.CheckBrandsResponseObject, fault error) {
	return brands.Checks(ctx, h.DBClient.Queries, r)
}

// UpdateBrand Update brand
// (PATCH /brand/{id})
func (h Handlers) UpdateBrand(ctx context.Context, r oapi.UpdateBrandRequestObject) (response oapi.UpdateBrandResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.UpdateBrand500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := brands.Update(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.UpdateBrand500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// FindBrands Find brands
// (GET /brands)
func (h Handlers) FindBrands(ctx context.Context, r oapi.FindBrandsRequestObject) (response oapi.FindBrandsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.FindBrands500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := brands.Find(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, nil
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.FindBrands500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// CreateBrand Create brand
// (POST /brands)
func (h Handlers) CreateBrand(ctx context.Context, r oapi.CreateBrandRequestObject) (response oapi.CreateBrandResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.CreateBrand500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := brands.Create(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.CreateBrand500JSONResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

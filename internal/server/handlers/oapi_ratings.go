package handlers

import (
	"context"
	"net/http"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/ratings"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// KillRating Kill rating

// (DELETE /rating/{id})

func (h Handlers) KillRating(ctx context.Context, r oapi.KillRatingRequestObject) (response oapi.KillRatingResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)

		return oapi.KillRating500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := ratings.Kill(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)

		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}

		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)

		return oapi.KillRating500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// CheckRating CORS preflight for rating by ID

// (OPTIONS /rating/{id})

func (h Handlers) CheckRating(ctx context.Context, r oapi.CheckRatingRequestObject) (response oapi.CheckRatingResponseObject, fault error) {
	return ratings.Check(ctx, h.DBClient.Queries, r)
}

// CheckRatings CORS preflight for rating by ID

// (OPTIONS /rating/{id})

func (h Handlers) CheckRatings(ctx context.Context, r oapi.CheckRatingsRequestObject) (response oapi.CheckRatingsResponseObject, fault error) {
	return ratings.Checks(ctx, h.DBClient.Queries, r)
}

// UpdateRating Update rating

// (PATCH /rating/{id})

func (h Handlers) UpdateRating(ctx context.Context, r oapi.UpdateRatingRequestObject) (response oapi.UpdateRatingResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)

		return oapi.UpdateRating500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := ratings.Update(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)

		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}

		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)

		return oapi.UpdateRating500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// FindRatings Find ratings

// (GET /ratings)

func (h Handlers) FindRatings(ctx context.Context, r oapi.FindRatingsRequestObject) (response oapi.FindRatingsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)

		return oapi.FindRatings500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := ratings.Find(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)

		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}

		return resp, nil
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)

		return oapi.FindRatings500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

// CreateRating Create rating

// (POST /ratings)

func (h Handlers) CreateRating(ctx context.Context, r oapi.CreateRatingRequestObject) (response oapi.CreateRatingResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)

		return oapi.CreateRating500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := ratings.Create(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)

		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}

		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)

		return oapi.CreateRating500JSONResponse{
			Message: err.Error(),

			Code: http.StatusInternalServerError,
		}, nil
	}

	return resp, nil
}

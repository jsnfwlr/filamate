package handlers

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/handlers/stats"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

// CheckUsageStats CORS preflight for usage stats by ID
// (OPTIONS /api/stats/usage)
func (h Handlers) CheckUsageStats(ctx context.Context, r oapi.CheckUsageStatsRequestObject) (response oapi.CheckUsageStatsResponseObject, fault error) {
	return stats.CheckUsage(ctx, h.DBClient.Queries, r)
}

// GetUsageStats gets usage stats
// (GET /api/stats/usage)
func (h Handlers) GetUsageStats(ctx context.Context, r oapi.GetUsageStatsRequestObject) (response oapi.GetUsageStatsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.GetUsageStats500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stats.GetUsageStats(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.GetUsageStats500JSONResponse{}, err
	}

	return resp, nil
}

// CheckStorageStats CORS preflight for storage stats by ID
// (OPTIONS /api/stats/storage)
func (h Handlers) CheckStorageStats(ctx context.Context, r oapi.CheckStorageStatsRequestObject) (response oapi.CheckStorageStatsResponseObject, fault error) {
	return stats.CheckStorage(ctx, h.DBClient.Queries, r)
}

// GetStorageStats gets storage stats
// (GET /api/stats/storage)
func (h Handlers) GetStorageStats(ctx context.Context, r oapi.GetStorageStatsRequestObject) (response oapi.GetStorageStatsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.GetStorageStats500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stats.GetStorageStats(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.GetStorageStats500JSONResponse{}, err
	}

	return resp, nil
}

// CheckRatingStats CORS preflight for rating stats by ID
// (OPTIONS /api/stats/rating)
func (h Handlers) CheckRatingStats(ctx context.Context, r oapi.CheckRatingStatsRequestObject) (response oapi.CheckRatingStatsResponseObject, fault error) {
	return stats.CheckRating(ctx, h.DBClient.Queries, r)
}

// GetRatingStats gets rating stats
// (GET /api/stats/rating)
func (h Handlers) GetRatingStats(ctx context.Context, r oapi.GetRatingStatsRequestObject) (response oapi.GetRatingStatsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.GetRatingStats500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stats.GetRatingStats(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.GetRatingStats500JSONResponse{}, err
	}

	return resp, nil
}

// CheckStorageChart CORS preflight for storage chart by ID
// (OPTIONS /api/chart/storage)
func (h Handlers) CheckStorageChart(ctx context.Context, r oapi.CheckStorageChartRequestObject) (response oapi.CheckStorageChartResponseObject, fault error) {
	return stats.CheckStorageChart(ctx, h.DBClient.Queries, r)
}

// GetStorageChart gets storage chart
// (GET /api/chart/storage)
func (h Handlers) GetStorageChart(ctx context.Context, r oapi.GetStorageChartRequestObject) (response oapi.GetStorageChartResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.GetStorageChart500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stats.GetStorageChart(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.GetStorageChart500JSONResponse{}, err
	}

	return resp, nil
}

// CheckMaterialChart CORS preflight for material chart by ID
// (OPTIONS /api/chart/material)
func (h Handlers) CheckMaterialChart(ctx context.Context, r oapi.CheckMaterialChartRequestObject) (response oapi.CheckMaterialChartResponseObject, fault error) {
	return stats.CheckMaterialChart(ctx, h.DBClient.Queries, r)
}

// GetMaterialChart gets material chart
// (GET /api/chart/material)
func (h Handlers) GetMaterialChart(ctx context.Context, r oapi.GetMaterialChartRequestObject) (response oapi.GetMaterialChartResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	tx, err := h.DBClient.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		o.Error("could not begin db transaction", err, go11y.SeverityHigh)
		return oapi.GetMaterialChart500JSONResponse{}, err
	}

	txQuerier := h.DBClient.Queries.WithTx(tx)

	resp, err := stats.GetMaterialChart(ctx, txQuerier, r)
	if err != nil {
		o.Error("request failed", err, go11y.SeverityHigh)
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			o.Error("could not rollback transaction", rbErr, go11y.SeverityHigh)
		}
		return resp, err
	}

	if err := tx.Commit(ctx); err != nil {
		o.Error("could not commit transaction", err, go11y.SeverityHigh)
		return oapi.GetMaterialChart500JSONResponse{}, err
	}

	return resp, nil
}

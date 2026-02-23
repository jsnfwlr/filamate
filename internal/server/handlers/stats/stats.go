// Package stats contains handlers for the dashboard data
package stats

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type statsQuerier interface {
	GetStorageStats(ctx context.Context) ([]db.GetStorageStatsRow, error)
	GetUsageStats(ctx context.Context) ([]db.GetUsageStatsRow, error)
}

// CheckUsage does CORS preflight for spool by ID
// (OPTIONS /api/stats/usage)
func CheckUsage(ctx context.Context, dbq statsQuerier, r oapi.CheckUsageStatsRequestObject) (response oapi.CheckUsageStatsResponseObject, fault error) {
	return oapi.CheckUsageStats204Response{}, nil
}

// CheckStorage does CORS preflight for spool by ID
// (OPTIONS /api/stats/storage)
func CheckStorage(ctx context.Context, dbq statsQuerier, r oapi.CheckStorageStatsRequestObject) (response oapi.CheckStorageStatsResponseObject, fault error) {
	return oapi.CheckStorageStats204Response{}, nil
}

// Find finds spool records
// (GET /spools)
// TODO add filtering, pagination, etc.
func GetUsageStats(ctx context.Context, dbq statsQuerier, r oapi.GetUsageStatsRequestObject) (response oapi.GetUsageStatsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	stats, err := dbq.GetUsageStats(ctx)
	if err != nil {
		o.Error("failed to find stats", err, go11y.SeverityHigh)

		return oapi.GetUsageStats500JSONResponse{
			Message: "Failed to find stats",
			Code:    500,
		}, err
	}

	var results []oapi.UsageStatsItem
	for _, s := range stats {
		result := oapi.UsageStatsItem{
			Color:    s.Color,
			Material: s.Material,
			Used:     s.Used,
			Ordered:  s.Ordered,
		}
		results = append(results, result)
	}
	return oapi.GetUsageStats200JSONResponse(results), nil
}

// GetStorageStats retrieves storage stats
// (GET /api/stats/storage)
func GetStorageStats(ctx context.Context, dbq statsQuerier, r oapi.GetStorageStatsRequestObject) (response oapi.GetStorageStatsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	stats, err := dbq.GetStorageStats(ctx)
	if err != nil {
		o.Error("failed to find stats", err, go11y.SeverityHigh)

		return oapi.GetStorageStats500JSONResponse{
			Message: "Failed to find stats",
			Code:    500,
		}, err
	}

	var results []oapi.StorageStatsItem
	for _, s := range stats {
		result := oapi.StorageStatsItem{
			Label: s.TallyLabel,
			Max:   s.Max,
			Used:  s.Used,
			Free:  s.Free,
		}
		results = append(results, result)
	}
	return oapi.GetStorageStats200JSONResponse(results), nil
}

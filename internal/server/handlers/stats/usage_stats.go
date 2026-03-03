// Package stats contains handlers for the dashboard data
package stats

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type usageStatsQuerier interface {
	GetUsageStats(ctx context.Context) ([]db.GetUsageStatsRow, error)
	GetStorageStats(ctx context.Context) ([]db.GetStorageStatsRow, error)
	FindSpools(ctx context.Context) ([]db.Spool, error)
	GetSpoolColors(ctx context.Context, spoolID int64) ([]db.Color, error)
}

// CheckUsage does CORS preflight for spool by ID
// (OPTIONS /api/stats/usage)
func CheckUsage(ctx context.Context, dbq usageStatsQuerier, r oapi.CheckUsageStatsRequestObject) (response oapi.CheckUsageStatsResponseObject, fault error) {
	return oapi.CheckUsageStats204Response{}, nil
}

// GetUsageStats retrieves usage stats
// (GET /api/stats/usage)
func GetUsageStats(ctx context.Context, dbq usageStatsQuerier, r oapi.GetUsageStatsRequestObject) (response oapi.GetUsageStatsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	stats, err := dbq.GetUsageStats(ctx)
	if err != nil {
		o.Error("failed to find stats", err, go11y.SeverityHigh)

		return oapi.GetUsageStats500JSONResponse{
			Message: fmt.Sprintf("failed to find stats: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	var resp []oapi.UsageStat
	for _, s := range stats {
		result := oapi.UsageStat{
			ID:       s.ID,
			Color:    s.Color,
			Material: s.Material,
			Used:     s.Used,
			Ordered:  s.Ordered,
		}
		resp = append(resp, result)
	}
	return oapi.GetUsageStats200JSONResponse(resp), nil
}

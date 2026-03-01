// Package stats contains handlers for the dashboard data
package stats

import (
	"context"
	"fmt"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type statsQuerier interface {
	GetStorageStats(ctx context.Context) ([]db.GetStorageStatsRow, error)
	GetUsageStats(ctx context.Context) ([]db.GetUsageStatsRow, error)
	FindSpools(ctx context.Context) ([]db.Spool, error)
	GetSpoolColors(ctx context.Context, spoolID int64) ([]db.Color, error)
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

// GetUsageStats retrieves usage stats
// (GET /api/stats/usage)
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

	var resp []oapi.UsageStat
	for _, s := range stats {
		result := oapi.UsageStat{
			Color:    s.Color,
			Material: s.Material,
			Used:     s.Used,
			Ordered:  s.Ordered,
		}
		resp = append(resp, result)
	}
	return oapi.GetUsageStats200JSONResponse(resp), nil
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

	spools, err := dbq.FindSpools(ctx)
	if err != nil {
		o.Error("failed to find spools", err, go11y.SeverityHigh)

		return oapi.GetStorageStats500JSONResponse{
			Message: "Failed to find spools",
			Code:    500,
		}, err
	}

	var resp []oapi.StorageStat
	for _, stat := range stats {
		if stat.ID == nil {
			continue
		}
		details := []oapi.StorageStatsDetails{}
		for _, spool := range spools {
			if spool.LocationID == *stat.ID {
				colors, err := dbq.GetSpoolColors(ctx, spool.ID)
				if err != nil {
					o.Error("failed to find spool colors", err, go11y.SeverityHigh)

					return oapi.GetStorageStats500JSONResponse{
						Message: "Failed to find spool colors",
						Code:    500,
					}, err
				}

				var colorsHex []string
				var colorsLabel []string
				for _, color := range colors {
					colorsHex = append(colorsHex, color.HexCode)
					colorsLabel = append(colorsLabel, color.Label)
				}

				rCurrentWeight, _ := spool.CurrentWeight.Float64Value()

				detail := oapi.StorageStatsDetails{
					Brand:         spool.BrandID,
					Material:      spool.MaterialID,
					CurrentWeight: fmt.Sprintf("%0.2f", rCurrentWeight.Float64),
					ColorsHex:     colorsHex,
					ColorsLabel:   colorsLabel,
				}

				details = append(details, detail)
			}
		}

		r := oapi.StorageStat{
			Label:   stat.TallyLabel,
			Max:     stat.Max,
			Used:    stat.Used,
			Free:    stat.Free,
			Details: details,
		}
		resp = append(resp, r)
	}
	return oapi.GetStorageStats200JSONResponse(resp), nil
}

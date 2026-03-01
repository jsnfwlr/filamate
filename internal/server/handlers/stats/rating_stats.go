// Package stats contains handlers for the dashboard data
package stats

import (
	"context"
	"encoding/json"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type ratingStatsQuerier interface {
	GetRatingStats(ctx context.Context) ([]db.GetRatingStatsRow, error)
}

// CheckRating does CORS preflight for spool by ID
// (OPTIONS /api/stats/rating)
func CheckRating(ctx context.Context, dbq ratingStatsQuerier, r oapi.CheckRatingStatsRequestObject) (response oapi.CheckRatingStatsResponseObject, fault error) {
	return oapi.CheckRatingStats204Response{}, nil
}

// GetRatingStats retrieves rating stats
// (GET /api/stats/rating)
func GetRatingStats(ctx context.Context, dbq ratingStatsQuerier, r oapi.GetRatingStatsRequestObject) (response oapi.GetRatingStatsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	stats, err := dbq.GetRatingStats(ctx)
	if err != nil {
		o.Error("failed to find stats", err, go11y.SeverityHigh)
		return oapi.GetRatingStats500JSONResponse{
			Message: "Failed to find stats",
			Code:    500,
		}, err
	}

	var resp []oapi.RatingStat
	for _, s := range stats {
		result := oapi.RatingStat{
			ID:            s.ID,
			Brand:         s.Brand,
			Material:      s.Material,
			RatingAverage: float32(s.AverageRating),
			RatingCount:   s.RatingCount,
		}

		details := []oapi.RatingStatDetail{}
		err := json.Unmarshal(s.Details, &details)
		if err != nil {
			o.Error("failed to unmarshal rating stat details", err, go11y.SeverityHigh)
			return oapi.GetRatingStats500JSONResponse{
				Message: "Failed to unmarshal rating stat details",
				Code:    500,
			}, err
		}
		result.Details = details

		resp = append(resp, result)
	}
	return oapi.GetRatingStats200JSONResponse(resp), nil
}

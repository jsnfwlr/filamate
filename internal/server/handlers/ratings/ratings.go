// Package ratings contains handlers for rating related API endpoints
package ratings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

type ratingsQuerier interface {
	CreateRating(ctx context.Context, rating int64, spoolID int64) (db.Rating, error)
	DeleteRating(ctx context.Context, id int64) error
	FindRatings(ctx context.Context) ([]db.Rating, error)
	GetRatingsByBrandIDAndMaterialID(ctx context.Context, brandID int64, materialID int64) ([]db.Rating, error)
	UpdateRating(ctx context.Context, rating int64, iD int64) (db.Rating, error)
	GetRatingByID(ctx context.Context, id int64) (db.Rating, error)
}

// Check does CORS preflight for rating by ID
// (OPTIONS /rating/{id})
func Check(ctx context.Context, dbq ratingsQuerier, r oapi.CheckRatingRequestObject) (response oapi.CheckRatingResponseObject, fault error) {
	return oapi.CheckRating204Response{}, nil
}

// Checks does CORS preflight for rating
// (OPTIONS /rating)
func Checks(ctx context.Context, dbq ratingsQuerier, r oapi.CheckRatingsRequestObject) (response oapi.CheckRatingsResponseObject, fault error) {
	return oapi.CheckRatings204Response{}, nil
}

// Find finds rating records
// (GET /ratings)
func Find(ctx context.Context, dbq ratingsQuerier, r oapi.FindRatingsRequestObject) (response oapi.FindRatingsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	var err error
	var ratings []db.Rating

	if r.Params.BrandID != nil && r.Params.MaterialID != nil {
		ratings, err = dbq.GetRatingsByBrandIDAndMaterialID(ctx, *r.Params.BrandID, *r.Params.MaterialID)
		if err != nil {
			o.Error("failed to get rating by brand id and material id", err, go11y.SeverityHigh, "brand_id", r.Params.BrandID, "material_id", r.Params.MaterialID)

			return oapi.FindRatings500JSONResponse{
				Message: fmt.Sprintf("failed to get rating by brand id and material id: %s", err.Error()),
				Code:    http.StatusInternalServerError,
			}, nil
		}
	} else {
		ratings, err = dbq.FindRatings(ctx)
		if err != nil {
			o.Error("failed to find ratings", err, go11y.SeverityHigh)

			return oapi.FindRatings500JSONResponse{
				Message: fmt.Sprintf("failed to find ratings: %s", err.Error()),
				Code:    http.StatusInternalServerError,
			}, nil
		}
	}

	if len(ratings) == 0 {
		o.Info("no ratings found")
		return oapi.FindRatings200JSONResponse([]oapi.Rating{}), nil
	}

	var resp []oapi.Rating

	for _, rating := range ratings {
		resp = append(resp, oapi.Rating{
			ID:        rating.ID,
			SpoolID:   rating.SpoolID,
			Rating:    rating.Rating,
			CreatedAt: rating.CreatedAt,
			UpdatedAt: rating.UpdatedAt,
		})
	}

	return oapi.FindRatings200JSONResponse(resp), nil
}

// Update updates a rating record
// (PATCH /rating/{id})
func Update(ctx context.Context, dbq ratingsQuerier, r oapi.UpdateRatingRequestObject) (response oapi.UpdateRatingResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	_, err := dbq.GetRatingByID(ctx, r.ID)
	if err != nil {
		o.Error("failed to get rating by id", err, go11y.SeverityHigh, "rating_id", r.ID)

		if err == pgx.ErrNoRows {
			return oapi.UpdateRating404JSONResponse{
				Message: fmt.Sprintf("failed to get record by id: %s", err.Error()),
				Code:    http.StatusNotFound,
			}, nil
		}

		return oapi.UpdateRating500JSONResponse{
			Message: fmt.Sprintf("failed to get record by id: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	rating, err := dbq.UpdateRating(ctx, r.Body.Rating, r.ID)
	if err != nil {
		o.Error("failed to update rating", err, go11y.SeverityHigh, "rating_id", r.ID)

		return oapi.UpdateRating500JSONResponse{
			Message: fmt.Sprintf("failed to update rating: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	resp := oapi.Rating{
		ID:        rating.ID,
		SpoolID:   rating.SpoolID,
		Rating:    rating.Rating,
		CreatedAt: rating.CreatedAt,
		UpdatedAt: rating.UpdatedAt,
	}

	return oapi.UpdateRating200JSONResponse(resp), nil
}

// Create creates a rating record
// (POST /ratings)
func Create(ctx context.Context, dbq ratingsQuerier, r oapi.CreateRatingRequestObject) (response oapi.CreateRatingResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	rating, err := dbq.CreateRating(ctx, r.Body.Rating, r.Body.SpoolID)
	if err != nil {
		o.Error("failed to create rating", err, go11y.SeverityHigh)

		return oapi.CreateRating500JSONResponse{
			Message: fmt.Sprintf("failed to create rating: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	o.Info("rating created", go11y.SeverityLow, "rating_id", rating.ID)
	resp := oapi.Rating{
		ID:        rating.ID,
		SpoolID:   rating.SpoolID,
		Rating:    rating.Rating,
		CreatedAt: rating.CreatedAt,
		UpdatedAt: rating.UpdatedAt,
	}

	return oapi.CreateRating201JSONResponse(resp), nil
}

// Kill deletes a rating record
// (DELETE /rating/{id})
func Kill(ctx context.Context, dbq ratingsQuerier, r oapi.KillRatingRequestObject) (response oapi.KillRatingResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	err := dbq.DeleteRating(ctx, r.ID)
	if err != nil {
		o.Error("failed to delete rating", err, go11y.SeverityHigh, "rating_id", r.ID)

		return oapi.KillRating500JSONResponse{
			Message: fmt.Sprintf("failed to delete rating: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return oapi.KillRating204Response{}, nil
}

// Package spools contains handlers for spool and spool-color related API endpoints
package spools

import (
	"context"
	"fmt"
	"time"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5/pgtype"
)

type spoolsQuerier interface {
	CreateSpool(ctx context.Context, params db.CreateSpoolParams) (db.Spool, error)
	DeleteSpool(ctx context.Context, id int64) error
	FindSpools(ctx context.Context) ([]db.Spool, error)
	GetSpoolByID(ctx context.Context, id int64) (db.Spool, error)
	UpdateSpool(ctx context.Context, params db.UpdateSpoolParams) (db.Spool, error)
	GetSpoolColors(ctx context.Context, spoolID int64) (colors []db.Color, fault error)
	ResetSpoolColor(ctx context.Context, spool int64) (fault error)
	AddSpoolColors(ctx context.Context, spool int64, colors []int64) (fault error)
}

// Kill deletes a spool record
// (DELETE /spool/{id})
func Kill(ctx context.Context, dbq spoolsQuerier, r oapi.KillSpoolRequestObject) (response oapi.KillSpoolResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	err := dbq.DeleteSpool(ctx, r.ID)
	if err != nil {
		o.Error("failed to delete spool", err, go11y.SeverityHigh, "spool_id", r.ID)

		return oapi.KillSpool500JSONResponse{
			Message: "Failed to delete spool",
			Code:    500,
		}, err
	}

	return oapi.KillSpool204Response{}, nil
}

// Check does CORS preflight for spool by ID
// (OPTIONS /spool/{id})
func Check(ctx context.Context, dbq spoolsQuerier, r oapi.CheckSpoolRequestObject) (response oapi.CheckSpoolResponseObject, fault error) {
	return oapi.CheckSpool204Response{}, nil
}

// Checks does CORS preflight for spools
// (OPTIONS /spools)
func Checks(ctx context.Context, dbq spoolsQuerier, r oapi.CheckSpoolsRequestObject) (response oapi.CheckSpoolsResponseObject, fault error) {
	return oapi.CheckSpools204Response{}, nil
}

// Update updates a spool record
// (PATCH /spool/{id})
func Update(ctx context.Context, dbq spoolsQuerier, r oapi.UpdateSpoolRequestObject) (response oapi.UpdateSpoolResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	og, err := dbq.GetSpoolByID(ctx, r.ID)
	if err != nil {
		o.Error("failed to get spool by id", err, go11y.SeverityHigh, "spool_id", r.ID)

		return oapi.UpdateSpool500JSONResponse{
			Message: "Failed to get spool",
			Code:    500,
		}, err
	}

	price := og.Price
	if r.Body.Price != nil {
		err = price.Scan(*r.Body.Price)
		if err != nil {
			o.Error("failed to scan price", err, go11y.SeverityHigh, "price", r.Body.Price)

			return oapi.UpdateSpool500JSONResponse{
				Message: "Failed to process price",
				Code:    500,
			}, err
		}
	}

	params := db.UpdateSpoolParams{
		ID:             r.ID,
		Location:       og.LocationID,
		Brand:          og.BrandID,
		Material:       og.MaterialID,
		Store:          og.StoreID,
		Weight:         og.Weight,
		CurrentWeight:  og.CurrentWeight,
		CombinedWeight: og.CombinedWeight,
		Price:          price,
		EmptiedAt:      og.EmptiedAt,
	}

	if r.Body.Location != nil {
		params.Location = *r.Body.Location
	}
	if r.Body.Brand != nil {
		params.Brand = *r.Body.Brand
	}
	if r.Body.Material != nil {
		params.Material = *r.Body.Material
	}
	if r.Body.Store != nil {
		params.Store = *r.Body.Store
	}
	if r.Body.Weight != nil {
		_ = params.Weight.Scan(*r.Body.Weight)
	}
	if r.Body.CurrentWeight != nil {
		_ = params.CurrentWeight.Scan(*r.Body.CurrentWeight)
	}
	if r.Body.CombinedWeight != nil {
		_ = params.CombinedWeight.Scan(*r.Body.CombinedWeight)
	}
	if r.Body.Empty != nil {
		if *r.Body.Empty {
			params.EmptiedAt = new(time.Now())
		} else {
			params.EmptiedAt = nil
		}
	}

	s, err := dbq.UpdateSpool(ctx, params)
	if err != nil {
		o.Error("failed to update spool", err, go11y.SeverityHigh, "spool_id", r.ID)

		return oapi.UpdateSpool500JSONResponse{
			Message: "Failed to update spool",
			Code:    500,
		}, err
	}

	if r.Body.Colors != nil {
		err = swapColors(ctx, dbq, s.ID, r.Body.Colors)
		if err != nil {
			o.Error("failed to swap spool colors", err, go11y.SeverityHigh, "spool_id", s.ID)

			return oapi.UpdateSpool500JSONResponse{
				Message: "Failed to update spool colors",
				Code:    500,
			}, err
		}
	}

	rPrice, _ := s.Price.Float64Value()
	rWeight, _ := s.Weight.Float64Value()
	rCurrentWeight, _ := s.CurrentWeight.Float64Value()
	rCombinedWeight, _ := s.CombinedWeight.Float64Value()

	resp := oapi.Spool{
		ID:             s.ID,
		Location:       s.LocationID,
		Brand:          s.BrandID,
		Material:       s.MaterialID,
		Store:          s.StoreID,
		Weight:         fmt.Sprintf("%.2f", rWeight.Float64),
		CurrentWeight:  fmt.Sprintf("%.2f", rCurrentWeight.Float64),
		CombinedWeight: fmt.Sprintf("%.2f", rCombinedWeight.Float64),
		Price:          fmt.Sprintf("%.2f", rPrice.Float64),
		Empty:          (s.EmptiedAt != nil),
		EmptiedAt:      s.EmptiedAt,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
		DeletedAt:      s.DeletedAt,
	}

	cSet, err := getColors(ctx, dbq, s.ID)
	if err != nil {
		o.Error("failed to get spool colors", err, go11y.SeverityHigh, "spool_id", s.ID)

		return oapi.UpdateSpool500JSONResponse{
			Message: "Failed to get spool colors",
			Code:    500,
		}, err
	}

	colors := make([]int64, 0, len(cSet))
	for _, c := range cSet {
		colors = append(colors, c.ID)
	}

	resp.Colors = colors

	return oapi.UpdateSpool200JSONResponse(resp), nil
}

// Find finds spool records
// (GET /spools)
// TODO add filtering, pagination, etc.
func Find(ctx context.Context, dbq spoolsQuerier, r oapi.FindSpoolsRequestObject) (response oapi.FindSpoolsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	spools, err := dbq.FindSpools(ctx)
	if err != nil {
		o.Error("failed to find spools", err, go11y.SeverityHigh)

		return oapi.FindSpools500JSONResponse{
			Message: "Failed to find spools",
			Code:    500,
		}, err
	}

	var resp []oapi.Spool
	for _, s := range spools {
		cSet, err := getColors(ctx, dbq, s.ID)
		if err != nil {
			o.Error("failed to get spool colors", err, go11y.SeverityHigh, "spool_id", s.ID)

			return oapi.FindSpools500JSONResponse{
				Message: "Failed to get spool colors",
				Code:    500,
			}, err
		}

		colors := make([]int64, 0, len(cSet))
		for _, c := range cSet {
			colors = append(colors, c.ID)
		}

		rPrice, _ := s.Price.Float64Value()
		rWeight, _ := s.Weight.Float64Value()
		rCurrentWeight, _ := s.CurrentWeight.Float64Value()
		rCombinedWeight, _ := s.CombinedWeight.Float64Value()

		r := oapi.Spool{
			ID:             s.ID,
			Colors:         colors,
			Location:       s.LocationID,
			Brand:          s.BrandID,
			Material:       s.MaterialID,
			Store:          s.StoreID,
			Weight:         fmt.Sprintf("%.2f", rWeight.Float64),
			CurrentWeight:  fmt.Sprintf("%.2f", rCurrentWeight.Float64),
			CombinedWeight: fmt.Sprintf("%.2f", rCombinedWeight.Float64),
			Price:          fmt.Sprintf("%.2f", rPrice.Float64),
			Empty:          (s.EmptiedAt != nil),
			EmptiedAt:      s.EmptiedAt,
			CreatedAt:      s.CreatedAt,
			UpdatedAt:      s.UpdatedAt,
			DeletedAt:      s.DeletedAt,
		}
		resp = append(resp, r)
	}
	return oapi.FindSpools200JSONResponse(resp), nil
}

// Create creates a spool record
// (POST /spools)
func Create(ctx context.Context, dbq spoolsQuerier, r oapi.CreateSpoolRequestObject) (response oapi.CreateSpoolResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	price := pgtype.Numeric{}
	err := price.Scan(r.Body.Price)
	if err != nil {
		o.Error("failed to scan price", err, go11y.SeverityHigh, "price", r.Body.Price)

		return oapi.CreateSpool500JSONResponse{
			Message: "Failed to process price",
			Code:    500,
		}, err
	}

	params := db.CreateSpoolParams{
		Location:  r.Body.Location,
		Brand:     r.Body.Brand,
		Material:  r.Body.Material,
		Store:     r.Body.Store,
		Price:     price,
		EmptiedAt: nil,
	}

	if r.Body.Empty != nil && *r.Body.Empty {
		params.EmptiedAt = new(time.Now())
	}

	err = params.Weight.Scan(r.Body.Weight)
	if err != nil {
		o.Error("failed to scan weight", err, go11y.SeverityHigh, "weight", r.Body.Weight)

		return oapi.CreateSpool500JSONResponse{
			Message: "Failed to process weight",
			Code:    500,
		}, err
	}

	err = params.CurrentWeight.Scan(r.Body.CurrentWeight)
	if err != nil {
		o.Error("failed to scan current weight", err, go11y.SeverityHigh, "current_weight", r.Body.CurrentWeight)

		return oapi.CreateSpool500JSONResponse{
			Message: "Failed to process current weight",
			Code:    500,
		}, err
	}

	err = params.CombinedWeight.Scan(r.Body.CombinedWeight)
	if err != nil {
		o.Error("failed to scan combined weight", err, go11y.SeverityHigh, "combined_weight", r.Body.CombinedWeight)

		return oapi.CreateSpool500JSONResponse{
			Message: "Failed to process combined weight",
			Code:    500,
		}, err
	}

	spool, err := dbq.CreateSpool(ctx, params)
	if err != nil {
		o.Error("failed to create spool", err, go11y.SeverityHigh)

		return oapi.CreateSpool500JSONResponse{
			Message: "Failed to create spool",
			Code:    500,
		}, err
	}

	if r.Body.Colors != nil {
		err = swapColors(ctx, dbq, spool.ID, r.Body.Colors)
		if err != nil {
			o.Error("failed to add color to spool", err, go11y.SeverityHigh, "spool_id", spool.ID, "color_id", r.Body.Colors)

			return oapi.CreateSpool500JSONResponse{
				Message: "Failed to add color to spool",
				Code:    500,
			}, err
		}
	}

	rPrice, _ := spool.Price.Float64Value()
	rWeight, _ := spool.Weight.Float64Value()
	rCurrentWeight, _ := spool.CurrentWeight.Float64Value()
	rCombinedWeight, _ := spool.CombinedWeight.Float64Value()

	resp := oapi.Spool{
		ID:             spool.ID,
		Location:       spool.LocationID,
		Brand:          spool.BrandID,
		Material:       spool.MaterialID,
		Store:          spool.StoreID,
		Weight:         fmt.Sprintf("%.2f", rWeight.Float64),
		CurrentWeight:  fmt.Sprintf("%.2f", rCurrentWeight.Float64),
		CombinedWeight: fmt.Sprintf("%.2f", rCombinedWeight.Float64),
		Price:          fmt.Sprintf("%.2f", rPrice.Float64),
		Empty:          (spool.EmptiedAt != nil),
		EmptiedAt:      spool.EmptiedAt,
		CreatedAt:      spool.CreatedAt,
		UpdatedAt:      spool.UpdatedAt,
		DeletedAt:      spool.DeletedAt,
	}

	if r.Body.Colors != nil {
		c, err := dbq.GetSpoolColors(ctx, spool.ID)
		if err != nil {
			o.Error("failed to get spool colors", err, go11y.SeverityHigh, "spool_id", spool.ID)

			return oapi.CreateSpool500JSONResponse{
				Message: "Failed to get spool colors",
				Code:    500,
			}, err
		}

		var colors []int64
		for _, c := range c {
			colors = append(colors, c.ID)
		}

		resp.Colors = colors
	}

	return oapi.CreateSpool201JSONResponse(resp), nil
}

func swapColors(ctx context.Context, dbq spoolsQuerier, spoolID int64, colorIDs *[]int64) error {
	ctx, o := go11y.Get(ctx)

	if colorIDs == nil {
		return nil
	}

	og, err := dbq.GetSpoolColors(ctx, spoolID)
	if err != nil {
		o.Error("failed to get original spool colors", err, go11y.SeverityHigh, "spool_id", spoolID)

		return err
	}

	// no point in resetting if there were no original colors
	if len(og) > 0 {
		err = dbq.ResetSpoolColor(ctx, spoolID)
		if err != nil {
			o.Error("failed to reset spool colors", err, go11y.SeverityHigh, "spool_id", spoolID)

			return err
		}
	}

	err = dbq.AddSpoolColors(ctx, spoolID, *colorIDs) // previously established colorIDs is not nil
	if err != nil {
		o.Error("failed to add spool colors", err, go11y.SeverityHigh, "spool_id", spoolID)

		return err
	}

	return nil
}

func getColors(ctx context.Context, dbq spoolsQuerier, spoolID int64) ([]oapi.Color, error) {
	ctx, o := go11y.Get(ctx)
	cs, err := dbq.GetSpoolColors(ctx, spoolID)
	if err != nil {
		o.Error("failed to get spool colors", err, go11y.SeverityHigh, "spool_id", spoolID)
		return nil, err
	}

	colors := make([]oapi.Color, 0, len(cs))
	for _, c := range cs {
		color := oapi.Color{
			ID:    c.ID,
			Label: c.Label,
			Hex:   c.HexCode,
			Alias: c.Alias,
		}

		colors = append(colors, color)
	}

	return colors, nil
}

// Package brands contains handlers for brand related API endpoints
package brands

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type brandsQuerier interface {
	CreateBrand(ctx context.Context, params db.CreateBrandParams) (db.Brand, error)
	DeleteBrand(ctx context.Context, id int64) error
	GetBrandByID(ctx context.Context, id int64) (db.Brand, error)
	UpdateBrand(ctx context.Context, params db.UpdateBrandParams) (db.Brand, error)
	FindBrands(ctx context.Context) ([]db.Brand, error)
}

// Check does CORS preflight for brand by ID
// (OPTIONS /brand/{id})
func Check(ctx context.Context, dbq brandsQuerier, r oapi.CheckBrandRequestObject) (response oapi.CheckBrandResponseObject, fault error) {
	return oapi.CheckBrand204Response{}, nil
}

// Checks does CORS preflight for brand
// (OPTIONS /brand)
func Checks(ctx context.Context, dbq brandsQuerier, r oapi.CheckBrandsRequestObject) (response oapi.CheckBrandsResponseObject, fault error) {
	return oapi.CheckBrands204Response{}, nil
}

// Find finds brand records
// (GET /brands)
func Find(ctx context.Context, dbq brandsQuerier, r oapi.FindBrandsRequestObject) (response oapi.FindBrandsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	brands, err := dbq.FindBrands(ctx)
	if err != nil {
		o.Error("failed to find brands", err, go11y.SeverityHigh)

		return oapi.FindBrands500JSONResponse{
			Message: "Failed to find brands",
			Code:    500,
		}, err
	}

	if len(brands) == 0 {
		o.Info("no brands found")

		return oapi.FindBrands200JSONResponse([]oapi.BrandItem{}), nil
	}

	var brandItems []oapi.BrandItem
	for _, b := range brands {
		brandItems = append(brandItems, oapi.BrandItem{
			ID:      b.ID,
			Label:   b.Label,
			Active:  b.Active,
			StoreID: b.StoreID,
		})
	}

	return oapi.FindBrands200JSONResponse(brandItems), nil
}

// Update updates a brand record
// (PATCH /brand/{id})
func Update(ctx context.Context, dbq brandsQuerier, r oapi.UpdateBrandRequestObject) (response oapi.UpdateBrandResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	og, err := dbq.GetBrandByID(ctx, r.ID)
	if err != nil {
		o.Error("failed to get brand by id", err, go11y.SeverityHigh, "brand_id", r.ID)

		return oapi.UpdateBrand500JSONResponse{
			Message: "Failed to get brand by id",
			Code:    500,
		}, err
	}

	params := db.UpdateBrandParams{
		ID:     r.ID,
		Label:  *r.Body.Label,
		Active: *r.Body.Active,
	}

	params.Store = og.StoreID

	if r.Body.StoreID != nil {
		params.Store = r.Body.StoreID
	}

	b, err := dbq.UpdateBrand(ctx, params)
	if err != nil {
		o.Error("failed to update brand", err, go11y.SeverityHigh, "brand_id", r.ID)

		return oapi.UpdateBrand500JSONResponse{
			Message: "Failed to update brand",
			Code:    500,
		}, err
	}

	resp := oapi.BrandItem{
		ID:      b.ID,
		Label:   b.Label,
		StoreID: b.StoreID,
	}

	return oapi.UpdateBrand200JSONResponse(resp), nil
}

// Create creates a brand record
// (POST /brands)
func Create(ctx context.Context, dbq brandsQuerier, r oapi.CreateBrandRequestObject) (response oapi.CreateBrandResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	params := db.CreateBrandParams{Label: r.Body.Label, Active: true, StoreID: r.Body.StoreID}

	b, err := dbq.CreateBrand(ctx, params)
	if err != nil {
		o.Error("failed to create brand", err, go11y.SeverityHigh)

		return oapi.CreateBrand500JSONResponse{
			Message: "Failed to create brand",
			Code:    500,
		}, err
	}

	o.Info("brand created", go11y.SeverityLow, "brand_id", b.ID)

	bi := oapi.BrandItem{
		ID:      b.ID,
		Label:   b.Label,
		StoreID: b.StoreID,
		Active:  b.Active,
	}

	return oapi.CreateBrand201JSONResponse(bi), nil
}

// Kill deletes a brand record
// (DELETE /brand/{id})
func Kill(ctx context.Context, dbq brandsQuerier, r oapi.KillBrandRequestObject) (response oapi.KillBrandResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)
	err := dbq.DeleteBrand(ctx, r.ID)
	if err != nil {
		o.Error("failed to delete brand", err, go11y.SeverityHigh, "brand_id", r.ID)

		return oapi.KillBrand500JSONResponse{
			Message: "Failed to delete brand",
			Code:    500,
		}, err
	}

	return oapi.KillBrand204Response{}, nil
}

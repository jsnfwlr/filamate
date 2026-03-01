// Package materials contains handlers for material related API endpoints
package materials

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type materialsQuerier interface {
	CreateMaterial(ctx context.Context, params db.CreateMaterialParams) (db.Material, error)
	DeleteMaterial(ctx context.Context, id int64) error
	FindMaterials(ctx context.Context) ([]db.Material, error)
	GetMaterialByID(ctx context.Context, id int64) (db.Material, error)
	UpdateMaterial(ctx context.Context, params db.UpdateMaterialParams) (db.Material, error)
}

// Kill deletes a material record
// (DELETE /material/{id})
func Kill(ctx context.Context, dbq materialsQuerier, r oapi.KillMaterialRequestObject) (response oapi.KillMaterialResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	err := dbq.DeleteMaterial(ctx, r.ID)
	if err != nil {
		o.Error("failed to delete material", err, go11y.SeverityHigh, "material_id", r.ID)

		return oapi.KillMaterial500JSONResponse{
			Message: "Failed to delete material",
			Code:    500,
		}, err
	}

	return oapi.KillMaterial204Response{}, nil
}

// Check does CORS preflight for material by ID
// (OPTIONS /material/{id})
func Check(ctx context.Context, dbq materialsQuerier, r oapi.CheckMaterialRequestObject) (response oapi.CheckMaterialResponseObject, fault error) {
	return oapi.CheckMaterial204Response{}, nil
}

// Checks does CORS preflight for materials
// (OPTIONS /materials)
func Checks(ctx context.Context, dbq materialsQuerier, r oapi.CheckMaterialsRequestObject) (response oapi.CheckMaterialsResponseObject, fault error) {
	return oapi.CheckMaterials204Response{}, nil
}

// Update updates a material record
// (PATCH /material/{id})
func Update(ctx context.Context, dbq materialsQuerier, r oapi.UpdateMaterialRequestObject) (response oapi.UpdateMaterialResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	_, og := dbq.GetMaterialByID(ctx, r.ID)
	if og != nil {
		o.Error("failed to get material by id", og, go11y.SeverityHigh, "material_id", r.ID)

		return oapi.UpdateMaterial500JSONResponse{
			Message: "Failed to get material",
			Code:    500,
		}, og
	}

	params := db.UpdateMaterialParams{
		ID:          r.ID,
		Label:       r.Body.Label,
		Description: r.Body.Description,
		Class:       r.Body.Class,
		Special:     r.Body.Special,
	}

	m, err := dbq.UpdateMaterial(ctx, params)
	if err != nil {
		o.Error("failed to update material", err, go11y.SeverityHigh)

		return oapi.UpdateMaterial500JSONResponse{
			Message: "Failed to update material",
			Code:    500,
		}, err
	}

	resp := oapi.Material{
		ID:          m.ID,
		Label:       m.Label,
		Description: m.Description,
		Class:       m.Class,
		Special:     m.Special,
	}

	return oapi.UpdateMaterial200JSONResponse(resp), nil
}

// Find finds material records
// (GET /materials)
// TODO add filtering, pagination, etc.
func Find(ctx context.Context, dbq materialsQuerier, r oapi.FindMaterialsRequestObject) (response oapi.FindMaterialsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	materials, err := dbq.FindMaterials(ctx)
	if err != nil {
		o.Error("failed to find materials", err, go11y.SeverityHigh)

		return oapi.FindMaterials500JSONResponse{
			Message: "Failed to find materials",
			Code:    500,
		}, err
	}

	var resp []oapi.Material
	for _, m := range materials {
		resp = append(resp, oapi.Material{
			ID:          m.ID,
			Class:       m.Class,
			Label:       m.Label,
			Description: m.Description,
			Special:     m.Special,
		})
	}

	return oapi.FindMaterials200JSONResponse(resp), nil
}

// Create creates a material record
// (POST /materials)
func Create(ctx context.Context, dbq materialsQuerier, r oapi.CreateMaterialRequestObject) (response oapi.CreateMaterialResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	params := db.CreateMaterialParams{
		Label:       r.Body.Label,
		Description: r.Body.Description,
		Class:       r.Body.Class,
		Special:     r.Body.Special,
	}

	m, err := dbq.CreateMaterial(ctx, params)
	if err != nil {
		o.Error("failed to create material", err, go11y.SeverityHigh)

		return oapi.CreateMaterial500JSONResponse{
			Message: "Failed to create material",
			Code:    500,
		}, err
	}

	resp := oapi.Material{
		ID:          m.ID,
		Label:       m.Label,
		Description: m.Description,
		Class:       m.Class,
		Special:     m.Special,
	}

	return oapi.CreateMaterial201JSONResponse(resp), nil
}

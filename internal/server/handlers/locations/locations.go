// Package locations contains handlers for location related API endpoints
package locations

import (
	"context"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

type locationsQuerier interface {
	CreateLocation(ctx context.Context, params db.CreateLocationParams) (db.Location, error)
	DeleteLocation(ctx context.Context, id int64) error
	FindLocations(ctx context.Context) ([]db.Location, error)
	GetLocationByID(ctx context.Context, id int64) (db.Location, error)
	UpdateLocation(ctx context.Context, params db.UpdateLocationParams) (db.Location, error)
}

// Kill deletes a location record
// (DELETE /location/{id})
func Kill(ctx context.Context, dbq locationsQuerier, r oapi.KillLocationRequestObject) (response oapi.KillLocationResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	err := dbq.DeleteLocation(ctx, r.ID)
	if err != nil {
		o.Error("failed to delete location", err, go11y.SeverityHigh, "location_id", r.ID)

		return oapi.KillLocation500JSONResponse{
			Message: "Failed to delete location",
			Code:    500,
		}, err
	}

	return oapi.KillLocation204Response{}, nil
}

// Check does CORS preflight for location by ID
// (OPTIONS /location/{id})
func Check(ctx context.Context, dbq locationsQuerier, r oapi.CheckLocationRequestObject) (response oapi.CheckLocationResponseObject, fault error) {
	return oapi.CheckLocation204Response{}, nil
}

// Checks does CORS preflight for locations
// (OPTIONS /locations)
func Checks(ctx context.Context, dbq locationsQuerier, r oapi.CheckLocationsRequestObject) (response oapi.CheckLocationsResponseObject, fault error) {
	return oapi.CheckLocations204Response{}, nil
}

// Update updates a location record
// (PATCH /location/{id})
func Update(ctx context.Context, dbq locationsQuerier, r oapi.UpdateLocationRequestObject) (response oapi.UpdateLocationResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	_, err := dbq.GetLocationByID(ctx, r.ID)
	if err != nil {
		o.Error("failed to get location by id", err, go11y.SeverityHigh, "location_id", r.ID)

		return oapi.UpdateLocation500JSONResponse{
			Message: "Failed to get location",
			Code:    500,
		}, err
	}

	params := db.UpdateLocationParams{
		ID:          r.ID,
		Label:       r.Body.Label,
		Description: r.Body.Description,
		Printable:   r.Body.Printable,
		Tally:       r.Body.Tally,
		Capacity:    int32(r.Body.Capacity),
	}

	l, err := dbq.UpdateLocation(ctx, params)
	if err != nil {
		o.Error("failed to update location", err, go11y.SeverityHigh, "location_id", r.ID)

		return oapi.UpdateLocation500JSONResponse{
			Message: "Failed to update location",
			Code:    500,
		}, err
	}

	return oapi.UpdateLocation200JSONResponse{
		ID:          l.ID,
		Label:       l.Label,
		Description: l.Description,
		Printable:   l.Printable,
		Tally:       l.Tally,
		Capacity:    int(l.Capacity),
	}, nil
}

// Find finds location records
// (GET /locations)
// TODO add filtering, pagination, etc.
func Find(ctx context.Context, dbq locationsQuerier, r oapi.FindLocationsRequestObject) (response oapi.FindLocationsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	locations, err := dbq.FindLocations(ctx)
	if err != nil {
		o.Error("failed to find locations", err, go11y.SeverityHigh)

		return oapi.FindLocations500JSONResponse{
			Message: "Failed to find locations",
			Code:    500,
		}, err
	}

	var resp []oapi.Location
	for _, l := range locations {
		resp = append(resp, oapi.Location{
			ID:          l.ID,
			Label:       l.Label,
			Description: l.Description,
			Printable:   l.Printable,
			Tally:       l.Tally,
			Capacity:    int(l.Capacity),
		})
	}

	return oapi.FindLocations200JSONResponse(resp), nil
}

// Create creates a location record
// (POST /locations)
func Create(ctx context.Context, dbq locationsQuerier, r oapi.CreateLocationRequestObject) (response oapi.CreateLocationResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	params := db.CreateLocationParams{
		Label:       r.Body.Label,
		Description: r.Body.Description,
		Printable:   r.Body.Printable,
		Tally:       r.Body.Tally,
		Capacity:    int32(r.Body.Capacity),
	}

	l, err := dbq.CreateLocation(ctx, params)
	if err != nil {
		o.Error("failed to create location", err, go11y.SeverityHigh)

		return oapi.CreateLocation500JSONResponse{
			Message: "Failed to create location",
			Code:    500,
		}, err
	}

	return oapi.CreateLocation201JSONResponse{
		ID:          l.ID,
		Label:       l.Label,
		Description: l.Description,
		Printable:   l.Printable,
		Tally:       l.Tally,
		Capacity:    int(l.Capacity),
	}, nil
}

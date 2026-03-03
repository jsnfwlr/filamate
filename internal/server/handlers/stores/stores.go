// Package stores contains handlers for store related API endpoints
package stores

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

type storesQuerier interface {
	CreateStore(ctx context.Context, label string, url *string) (db.Store, error)
	DeleteStore(ctx context.Context, id int64) error
	FindStores(ctx context.Context) ([]db.Store, error)
	GetStoreByID(ctx context.Context, id int64) (db.Store, error)
	UpdateStore(ctx context.Context, params db.UpdateStoreParams) (db.Store, error)
}

// Kill deletes a store record
// (DELETE /store/{id})
func Kill(ctx context.Context, dbq storesQuerier, r oapi.KillStoreRequestObject) (response oapi.KillStoreResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	err := dbq.DeleteStore(ctx, r.ID)
	if err != nil {
		o.Error("failed to delete store", err, go11y.SeverityHigh, "store_id", r.ID)

		return oapi.KillStore500JSONResponse{
			Message: fmt.Sprintf("failed to delete store: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return oapi.KillStore204Response{}, nil
}

// Check does CORS preflight for store by ID
// (OPTIONS /store/{id})
func Check(ctx context.Context, dbq storesQuerier, r oapi.CheckStoreRequestObject) (response oapi.CheckStoreResponseObject, fault error) {
	return oapi.CheckStore204Response{}, nil
}

// Checks does CORS preflight for stores
// (OPTIONS /stores)
func Checks(ctx context.Context, dbq storesQuerier, r oapi.CheckStoresRequestObject) (response oapi.CheckStoresResponseObject, fault error) {
	return oapi.CheckStores204Response{}, nil
}

// Update updates a store record
// (PATCH /store/{id})
func Update(ctx context.Context, dbq storesQuerier, r oapi.UpdateStoreRequestObject) (response oapi.UpdateStoreResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	og, err := dbq.GetStoreByID(ctx, r.ID)
	if err != nil {
		o.Error("failed to get store by id", err, go11y.SeverityHigh, "store_id", r.ID)

		if err == pgx.ErrNoRows {
			return oapi.UpdateStore404JSONResponse{
				Message: fmt.Sprintf("failed to get record by id: %s", err.Error()),
				Code:    http.StatusNotFound,
			}, nil
		}

		return oapi.UpdateStore500JSONResponse{
			Message: fmt.Sprintf("failed to get record by id: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	params := db.UpdateStoreParams{
		ID:    r.ID,
		Label: r.Body.Label,
	}

	params.URL = og.URL

	if r.Body.URL != nil {
		params.URL = r.Body.URL
	}

	s, err := dbq.UpdateStore(ctx, params)
	if err != nil {
		o.Error("failed to update store", err, go11y.SeverityHigh, "store_id", r.ID)

		return oapi.UpdateStore500JSONResponse{
			Message: fmt.Sprintf("failed to update store: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	resp := oapi.UpdateStore200JSONResponse{
		ID:    s.ID,
		Label: s.Label,
		URL:   s.URL,
	}

	return resp, nil
}

// Find finds store records
// (GET /stores)
// TODO add filtering, pagination, etc.
func Find(ctx context.Context, dbq storesQuerier, r oapi.FindStoresRequestObject) (response oapi.FindStoresResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	stores, err := dbq.FindStores(ctx)
	if err != nil {
		o.Error("failed to find stores", err, go11y.SeverityHigh)

		return oapi.FindStores500JSONResponse{
			Message: fmt.Sprintf("failed to find stores: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	var resp []oapi.Store
	for _, s := range stores {
		store := oapi.Store{
			ID:    s.ID,
			Label: s.Label,
			URL:   s.URL,
		}

		resp = append(resp, store)
	}

	return oapi.FindStores200JSONResponse(resp), nil
}

// Create creates a store record
// (POST /stores)
func Create(ctx context.Context, dbq storesQuerier, r oapi.CreateStoreRequestObject) (response oapi.CreateStoreResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	s, err := dbq.CreateStore(ctx, r.Body.Label, r.Body.URL)
	if err != nil {
		o.Error("failed to create store", err, go11y.SeverityHigh, "label", r.Body.Label)

		return oapi.CreateStore500JSONResponse{
			Message: fmt.Sprintf("failed to create store: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	resp := oapi.CreateStore201JSONResponse{
		ID:    s.ID,
		Label: s.Label,
		URL:   s.URL,
	}

	return resp, nil
}

// Package colors contains handlers for color related API endpoints
package colors

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/oapi"

	"github.com/jackc/pgx/v5"
)

type colorsQuerier interface {
	CreateColor(ctx context.Context, params db.CreateColorParams) (db.Color, error)
	DeleteColor(ctx context.Context, id int64) error
	FindColors(ctx context.Context) ([]db.Color, error)
	GetColorByID(ctx context.Context, id int64) (db.Color, error)
	UpdateColor(ctx context.Context, params db.UpdateColorParams) (db.Color, error)
}

// Kill deletes a color record
// (DELETE /color/{id})
func Kill(ctx context.Context, dbq colorsQuerier, r oapi.KillColorRequestObject) (response oapi.KillColorResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	err := dbq.DeleteColor(ctx, r.ID)
	if err != nil {
		o.Error("failed to delete color", err, go11y.SeverityHigh, "color_id", r.ID)

		return oapi.KillColor500JSONResponse{
			Message: fmt.Sprintf("failed to delete color: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	return oapi.KillColor204Response{}, nil
}

// Check does CORS preflight for color by ID
// (OPTIONS /color/{id})
func Check(ctx context.Context, dbq colorsQuerier, r oapi.CheckColorRequestObject) (response oapi.CheckColorResponseObject, fault error) {
	return oapi.CheckColor204Response{}, nil
}

// Checks does CORS preflight for colors
// (OPTIONS /colors)
func Checks(ctx context.Context, dbq colorsQuerier, r oapi.CheckColorsRequestObject) (response oapi.CheckColorsResponseObject, fault error) {
	return oapi.CheckColors204Response{}, nil
}

// Update updates a color record
// (PATCH /color/{id})
func Update(ctx context.Context, dbq colorsQuerier, r oapi.UpdateColorRequestObject) (response oapi.UpdateColorResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	hex, err := validateHexCode(r.Body.Hex)
	if err != nil {
		o.Error("invalid hex code provided", err, go11y.SeverityMedium, "hex_code", r.Body.Hex)

		return oapi.UpdateColor500JSONResponse{
			Message: fmt.Sprintf("invalid hex code provided: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	og, err := dbq.GetColorByID(ctx, r.ID)
	if err != nil {
		o.Error("failed to get color by id", err, go11y.SeverityHigh, "color_id", r.ID)

		if err == pgx.ErrNoRows {
			return oapi.UpdateColor404JSONResponse{
				Message: fmt.Sprintf("failed to get record by id: %s", err.Error()),
				Code:    http.StatusNotFound,
			}, nil
		}

		return oapi.UpdateColor500JSONResponse{
			Message: fmt.Sprintf("failed to get record by id: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	params := db.UpdateColorParams{
		ID:      r.ID,
		Label:   r.Body.Label,
		HexCode: hex,
	}

	params.Alias = og.Alias

	if r.Body.Alias != nil {
		params.Alias = r.Body.Alias
	}

	c, err := dbq.UpdateColor(ctx, params)
	if err != nil {
		o.Error("failed to update color", err, go11y.SeverityHigh, "color_id", r.ID)

		return oapi.UpdateColor500JSONResponse{
			Message: fmt.Sprintf("failed to update color: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	resp := oapi.UpdateColor200JSONResponse{
		ID:    c.ID,
		Label: c.Label,
		Hex:   c.HexCode,
		Alias: c.Alias,
	}

	return resp, nil
}

// Find finds color records
// (GET /colors)
// TODO add filtering, pagination, etc.
func Find(ctx context.Context, dbq colorsQuerier, r oapi.FindColorsRequestObject) (response oapi.FindColorsResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	colors, err := dbq.FindColors(ctx)
	if err != nil {
		o.Error("failed to find colors", err, go11y.SeverityHigh)

		return oapi.FindColors500JSONResponse{
			Message: fmt.Sprintf("failed to find colors: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	var resp oapi.FindColors200JSONResponse
	for _, c := range colors {
		color := oapi.Color{
			ID:    c.ID,
			Label: c.Label,
			Hex:   c.HexCode,
			Alias: c.Alias,
		}

		resp = append(resp, color)
	}

	return resp, nil
}

// Create creates a color record
// (POST /colors)
func Create(ctx context.Context, dbq colorsQuerier, r oapi.CreateColorRequestObject) (response oapi.CreateColorResponseObject, fault error) {
	ctx, o := go11y.Get(ctx)

	hex, err := validateHexCode(r.Body.Hex)
	if err != nil {
		o.Error("invalid hex code provided", err, go11y.SeverityMedium, "hex_code", r.Body.Hex)

		return oapi.CreateColor400JSONResponse{
			Message: fmt.Sprintf("invalid hex code provided: %s", err.Error()),
			Code:    400,
		}, nil
	}

	params := db.CreateColorParams{
		Label:   r.Body.Label,
		HexCode: hex,
		Alias:   r.Body.Alias,
	}

	c, err := dbq.CreateColor(ctx, params)
	if err != nil {
		o.Error("failed to create color", err, go11y.SeverityHigh)

		return oapi.CreateColor500JSONResponse{
			Message: fmt.Sprintf("failed to create color: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}, nil
	}

	resp := oapi.CreateColor201JSONResponse{
		ID:    c.ID,
		Label: c.Label,
		Hex:   c.HexCode,
		Alias: c.Alias,
	}

	return resp, nil
}

func validateHexCode(hex string) (validHex string, fault error) {
	rex := regexp.MustCompile(`^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$`)

	if !rex.MatchString(hex) {
		return "", fmt.Errorf("invalid hex code format")
	}

	if len(hex) == 4 {
		// Expand shorthand hex code (#RGB to #RRGGBB)
		expanded := "#" + string(hex[1]) + string(hex[1]) + string(hex[2]) + string(hex[2]) + string(hex[3]) + string(hex[3])
		return strings.ToUpper(expanded), nil
	}

	return strings.ToUpper(hex), nil
}

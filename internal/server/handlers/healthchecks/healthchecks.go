// Package healthchecks contains handlers for health check API endpoints
package healthchecks

import (
	"context"
	"fmt"

	"github.com/jsnfwlr/go11y"

	migrations "github.com/jsnfwlr/filamate/etc/db/migrations"
	"github.com/jsnfwlr/filamate/internal/db"
	"github.com/jsnfwlr/filamate/internal/server/log"
	"github.com/jsnfwlr/filamate/internal/server/oapi"
)

var tracer = go11y.NewTracer("github.com/jsnfwlr/filamate/internal/server/handlers/healthchecks")

type healthchecksQuerier interface {
	CheckMigration(ctx context.Context) (int32, error)
}

// HealthCheck handles the GET request for the health check endpoint
func HealthCheck(ctx context.Context, dbc healthchecksQuerier, r oapi.HealthCheckRequestObject) (res oapi.HealthCheckResponseObject, fault error) {
	ctx, o := go11y.Span(ctx, tracer, "doHealthCheck", go11y.SpanKindServer)
	defer o.End()

	o.Debug("healthCheck request", log.RequestIDKey, go11y.GetRequestID(ctx))

	coll, err := migrations.New()
	if err != nil {
		o.Error("health check failure - unable to read database migrations", err, go11y.SeverityHighest)

		return oapi.HealthCheck500JSONResponse{
			Message: "unable to read database migrations",
			Code:    500,
		}, err
	}

	// Check that filamate is connected to the database by running a query - auth should fail to start if it cannot connect to the database
	// but if the database goes down after startup this will catch it
	version, err := dbc.CheckMigration(ctx)
	if err != nil {
		o.Error("health check failure - unable to query database", err, go11y.SeverityHighest)

		return oapi.HealthCheck500JSONResponse{
			Message: "unable to query database",
			Code:    500,
		}, err
	}

	if version != coll.Steps() {
		err = fmt.Errorf("database migrations (%d) not up to date (%d)", version, coll.Steps())
		o.Error("health check failure", err, go11y.SeverityHighest, db.DBVersionKey, version, db.TotalStepsKey, coll.Steps())

		return oapi.HealthCheck500JSONResponse{
			Message: fmt.Sprintf("database migrations (%d) not up to date (%d)", version, coll.Steps()),
			Code:    500,
		}, err
	}

	return oapi.HealthCheck200Response{}, nil
}

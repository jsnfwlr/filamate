package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jsnfwlr/go11y"

	"github.com/jsnfwlr/filamate/internal/server/log"
)

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code      int
	Err       error
	RequestID string
}

// Error allows StatusError to satisfy the error interface.
func (statusErr StatusError) Error() string {
	return statusErr.Err.Error()
}

func (statusErr StatusError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code      int    `json:"code"`
		Error     string `json:"error"`
		RequestID string `json:"request_id"`
	}{
		Code:      statusErr.Code,
		Error:     statusErr.Error(),
		RequestID: statusErr.RequestID,
	})
}

// Status returns our HTTP status code.
func (statusErr StatusError) Status() int {
	return statusErr.Code
}

// String returns our HTTP status code.
func (statusErr StatusError) String() string {
	return fmt.Sprintf("{ \"code\": %d, \"error\": \"%s\", \"request_id\": \"%s\"}", statusErr.Code, statusErr.Err.Error(), statusErr.RequestID)
}

func NewStatusError(ctx context.Context, code int, err error) StatusError {
	r := go11y.GetRequestID(ctx)

	return StatusError{
		RequestID: r,
		Code:      code,
		Err:       err,
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var body []byte
	ctx, o := go11y.Span(r.Context(), tracer, "errorHandler", go11y.SpanKindServer)
	defer o.End()

	e := StatusError{}

	if !errors.As(err, &e) {
		e = NewStatusError(ctx, http.StatusInternalServerError, err)
	} else {
		e = err.(StatusError)
	}

	body, _ = json.Marshal(e)

	o.Error("error handling request", e, log.RequestBodyKey, string(body))

	http.Error(w, string(body), e.Code)
}

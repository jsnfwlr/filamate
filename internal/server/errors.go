package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jsnfwlr/go11y"
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

func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx, o := go11y.Span(r.Context(), tracer, "requestErrorHandler", go11y.SpanKindServer)
	defer o.End()

	e := NewStatusError(ctx, http.StatusInternalServerError, err)

	switch {
	case strings.Contains(err.Error(), "can't decode JSON body") || strings.Contains(err.Error(), "invalid character") || strings.Contains(err.Error(), "cannot unmarshal"):
		e.Code = http.StatusBadRequest
	}

	o.Error("error handling request", err, go11y.SeverityMedium)

	body, _ := json.Marshal(e)
	http.Error(w, string(body), e.Code)
}

func responseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx, o := go11y.Span(r.Context(), tracer, "responseErrorHandler", go11y.SpanKindServer)
	defer o.End()

	e := NewStatusError(ctx, http.StatusInternalServerError, err)

	rbByte, _ := io.ReadAll(r.Body)
	var requestBody map[string]interface{}
	if err := json.Unmarshal(rbByte, &requestBody); err != nil {
		requestBody = map[string]interface{}{
			"raw": string(rbByte),
		}
	}

	o.Error("error handling response", err, go11y.SeverityMedium, "request_body", requestBody)

	body, _ := json.Marshal(e)
	http.Error(w, string(body), http.StatusInternalServerError)
}

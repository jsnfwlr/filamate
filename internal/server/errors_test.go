package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jsnfwlr/go11y"
)

func TestStatusError(t *testing.T) {
	tests := []struct {
		name           string
		statusError    StatusError
		expectedError  string
		expectedStatus int
		expectedString string
	}{
		{
			name: "basic error",
			statusError: StatusError{
				Code:      http.StatusBadRequest,
				Err:       errors.New("test error"),
				RequestID: "req-123",
			},
			expectedError:  "test error",
			expectedStatus: http.StatusBadRequest,
			expectedString: `{ "code": 400, "error": "test error", "request_id": "req-123"}`,
		},
		{
			name: "internal server error",
			statusError: StatusError{
				Code:      http.StatusInternalServerError,
				Err:       errors.New("internal server error"),
				RequestID: "req-456",
			},
			expectedError:  "internal server error",
			expectedStatus: http.StatusInternalServerError,
			expectedString: `{ "code": 500, "error": "internal server error", "request_id": "req-456"}`,
		},
		{
			name: "not found error",
			statusError: StatusError{
				Code:      http.StatusNotFound,
				Err:       errors.New("resource not found"),
				RequestID: "req-789",
			},
			expectedError:  "resource not found",
			expectedStatus: http.StatusNotFound,
			expectedString: `{ "code": 404, "error": "resource not found", "request_id": "req-789"}`,
		},
		{
			name: "empty request id",
			statusError: StatusError{
				Code:      http.StatusUnauthorized,
				Err:       errors.New("unauthorized"),
				RequestID: "",
			},
			expectedError:  "unauthorized",
			expectedStatus: http.StatusUnauthorized,
			expectedString: `{ "code": 401, "error": "unauthorized", "request_id": ""}`,
		},
		{
			name: "special characters in error",
			statusError: StatusError{
				Code:      http.StatusBadRequest,
				Err:       errors.New(`error with "quotes" and \backslashes`),
				RequestID: "req-special",
			},
			expectedError:  `error with "quotes" and \backslashes`,
			expectedStatus: http.StatusBadRequest,
			expectedString: `{ "code": 400, "error": "error with "quotes" and \backslashes", "request_id": "req-special"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Error() method
			if got := tt.statusError.Error(); got != tt.expectedError {
				t.Errorf("StatusError.Error() = %q, want %q", got, tt.expectedError)
			}

			// Test Status() method
			if got := tt.statusError.Status(); got != tt.expectedStatus {
				t.Errorf("StatusError.Status() = %d, want %d", got, tt.expectedStatus)
			}

			// Test String() method
			if got := tt.statusError.String(); got != tt.expectedString {
				t.Errorf("Received\n%q\nExpected\n%q", got, tt.expectedString)
			}
		})
	}
}

func TestStatusErrorMarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		statusError  StatusError
		expectedJSON string
	}{
		{
			name: "basic json marshal",
			statusError: StatusError{
				Code:      http.StatusBadRequest,
				Err:       errors.New("validation failed"),
				RequestID: "req-json-123",
			},
			expectedJSON: `{"code":400,"error":"validation failed","request_id":"req-json-123"}`,
		},
		{
			name: "json marshal with empty fields",
			statusError: StatusError{
				Code:      http.StatusOK,
				Err:       errors.New(""),
				RequestID: "",
			},
			expectedJSON: `{"code":200,"error":"","request_id":""}`,
		},
		{
			name: "json marshal with special characters",
			statusError: StatusError{
				Code:      http.StatusTeapot,
				Err:       errors.New("I'm a teapot with \"quotes\" and newlines\n"),
				RequestID: "req-teapot",
			},
			expectedJSON: `{"code":418,"error":"I'm a teapot with \"quotes\" and newlines\n","request_id":"req-teapot"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := tt.statusError.MarshalJSON()
			if err != nil {
				t.Errorf("StatusError.MarshalJSON() error = %v", err)
				return
			}

			if got := string(jsonBytes); got != tt.expectedJSON {
				t.Errorf("StatusError.MarshalJSON() = %q, want %q", got, tt.expectedJSON)
			}

			// Verify JSON is valid by unmarshaling
			var unmarshaled map[string]any
			if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
				t.Errorf("Generated JSON is invalid: %v", err)
			}
		})
	}
}

func TestNewStatusError(t *testing.T) {
	tests := []struct {
		name          string
		code          int
		err           error
		expectedCode  int
		expectedError string
	}{
		{
			name: "with request id in context",

			code:          http.StatusBadRequest,
			err:           errors.New("test error"),
			expectedCode:  http.StatusBadRequest,
			expectedError: "test error",
		},
		{
			name: "without request id in context",

			code:          http.StatusInternalServerError,
			err:           errors.New("server error"),
			expectedCode:  http.StatusInternalServerError,
			expectedError: "server error",
		},
		{
			name: "nil error",

			code:          http.StatusBadRequest,
			err:           nil,
			expectedCode:  http.StatusBadRequest,
			expectedError: "",
		},
		{
			name: "different status codes",

			code:          http.StatusNotFound,
			err:           errors.New("not found"),
			expectedCode:  http.StatusNotFound,
			expectedError: "not found",
		},
	}

	gCfg, err := go11y.LoadConfig()
	if err != nil {
		t.Fatalf("could not load go11y config: %v", err)
	}

	ctx, _, err := go11y.Initialise(context.Background(), gCfg, os.Stdout)
	if err != nil {
		t.Fatalf("could not initialise go11y: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusErr := NewStatusError(ctx, tt.code, tt.err)

			if statusErr.Code != tt.expectedCode {
				t.Errorf("NewStatusError() Code = %d, want %d", statusErr.Code, tt.expectedCode)
			}

			if tt.err == nil {
				if statusErr.Err != nil {
					t.Errorf("NewStatusError() Err = %v, want nil", statusErr.Err)
				}
			} else {
				if statusErr.Err == nil || statusErr.Err.Error() != tt.expectedError {
					t.Errorf("NewStatusError() Err = %v, want %v", statusErr.Err, tt.err)
				}
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name               string
		inputError         error
		expectedStatusCode int
		expectedBodyCheck  func(body string) bool
		description        string
	}{
		{
			name: "status error handling",

			inputError: StatusError{
				Code:      http.StatusBadRequest,
				Err:       errors.New("custom status error"),
				RequestID: "handler-test-1",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBodyCheck: func(body string) bool {
				return strings.Contains(body, "custom status error") &&
					strings.Contains(body, "handler-test-1") &&
					strings.Contains(body, "400")
			},
			description: "should handle StatusError directly",
		},
		{
			name: "regular error handling",

			inputError:         errors.New("regular error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyCheck: func(body string) bool {
				return strings.Contains(body, "regular error") &&
					strings.Contains(body, "500")
			},
			description: "should wrap regular error as StatusError",
		},
		{
			name: "error without request id",

			inputError:         errors.New("no request id error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyCheck: func(body string) bool {
				return strings.Contains(body, "no request id error") &&
					strings.Contains(body, "500")
			},
			description: "should handle error without request ID in context",
		},
		{
			name: "complex status error",

			inputError: StatusError{
				Code:      http.StatusUnprocessableEntity,
				Err:       errors.New("validation failed: field 'name' is required"),
				RequestID: "complex-test",
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBodyCheck: func(body string) bool {
				return strings.Contains(body, "validation failed") &&
					strings.Contains(body, "field 'name' is required") &&
					strings.Contains(body, "complex-test") &&
					strings.Contains(body, "422")
			},
			description: "should handle complex validation errors",
		},
	}

	gCfg, err := go11y.LoadConfig()
	if err != nil {
		t.Fatalf("could not load go11y config: %v", err)
	}

	ctx, _, err := go11y.Initialise(context.Background(), gCfg, os.Stdout)
	if err != nil {
		t.Fatalf("could not initialise go11y: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			w := httptest.NewRecorder()

			// Create a request with the test context
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req = req.WithContext(ctx)

			// Call the error handler
			errorHandler(w, req, tt.inputError)

			// Check status code
			if w.Code != tt.expectedStatusCode {
				t.Errorf("errorHandler() status code = %d, want %d", w.Code, tt.expectedStatusCode)
			}

			// Check content type
			if contentType := w.Header().Get("Content-Type"); !strings.Contains(contentType, "text/plain") {
				t.Errorf("errorHandler() Content-Type = %q, want to contain 'text/plain'", contentType)
			}

			// Check body content
			body := w.Body.String()
			if tt.expectedBodyCheck != nil && !tt.expectedBodyCheck(body) {
				t.Errorf("errorHandler() body check failed for %s. Body: %s", tt.description, body)
			}

			// Verify the body contains valid JSON
			if body != "" {
				var jsonObj map[string]any
				if err := json.Unmarshal([]byte(body), &jsonObj); err != nil {
					t.Errorf("errorHandler() body is not valid JSON: %v. Body: %s", err, body)
				}

				// Verify required JSON fields exist
				if _, ok := jsonObj["code"]; !ok {
					t.Error("errorHandler() response missing 'code' field")
				}
				if _, ok := jsonObj["error"]; !ok {
					t.Error("errorHandler() response missing 'error' field")
				}
				if _, ok := jsonObj["request_id"]; !ok {
					t.Error("errorHandler() response missing 'request_id' field")
				}
			}
		})
	}
}

func TestStatusErrorImplementsError(t *testing.T) {
	// Ensure StatusError implements the error interface
	var _ error = StatusError{}

	statusErr := StatusError{
		Code:      http.StatusBadRequest,
		Err:       errors.New("interface test"),
		RequestID: "interface-req",
	}

	// Should be usable as an error
	var err error = statusErr
	if err.Error() != "interface test" {
		t.Errorf("StatusError as error interface Error() = %q, want %q", err.Error(), "interface test")
	}
}

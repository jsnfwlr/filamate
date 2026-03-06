package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/jsnfwlr/go11y"
)

type StackTraceEntry struct {
	Call     string `json:"call"`
	Location string `json:"location"`
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if cause := recover(); cause != nil {
				_, o := go11y.Get(r.Context())

				stack := strings.Split(strings.TrimSpace(string(debug.Stack())), "\n")

				stackTrace := []StackTraceEntry{
					{
						Call:     stack[0],
						Location: "panic",
					},
				}

				for i := 1; i < len(stack)-1; i += 2 {
					stackTrace = append(stackTrace, StackTraceEntry{
						Call:     stack[i],
						Location: stack[i+1],
					})
				}

				o.Error("Panic recovered", fmt.Errorf("panic: %v", cause), go11y.SeverityHighest, "stack_track", stackTrace, "full_trace", stack)

				// Set appropriate headers and write a 500 error response
				w.Header().Set("Connection", "close")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

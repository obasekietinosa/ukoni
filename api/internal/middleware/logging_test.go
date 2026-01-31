package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Setup logger to write to buffer
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	mw := NewLoggingMiddleware(logger)

	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
		expectedLog    map[string]interface{}
		expectError    bool
	}{
		{
			name: "Success 200",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			},
			expectedStatus: http.StatusOK,
			expectedLog: map[string]interface{}{
				"level":  "INFO",
				"msg":    "request completed",
				"status": float64(200),
			},
		},
		{
			name: "Error 400 with UserID",
			handler: func(w http.ResponseWriter, r *http.Request) {
				logData := GetRequestLogData(r.Context())
				if logData != nil {
					logData.UserID = "user-123"
				}
				http.Error(w, "bad request", http.StatusBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectedLog: map[string]interface{}{
				"level":         "WARN",
				"msg":           "request failed",
				"status":        float64(400),
				"user_id":       "user-123",
				"error_details": "bad request\n",
			},
		},
		{
			name: "Error 500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal error", http.StatusInternalServerError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedLog: map[string]interface{}{
				"level":         "ERROR",
				"msg":           "request failed",
				"status":        float64(500),
				"error_details": "internal error\n",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			req := httptest.NewRequest("GET", "/test", nil)
			rr := httptest.NewRecorder()

			handler := mw.Handler(tc.handler)
			handler.ServeHTTP(rr, req)

			// Check Response
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check X-Request-ID header
			if rr.Header().Get("X-Request-ID") == "" {
				t.Error("expected X-Request-ID header")
			}

			// Check Log
			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("failed to parse log output: %v", err)
			}

			for k, v := range tc.expectedLog {
				if logEntry[k] != v {
					t.Errorf("expected log field %s to be %v, got %v", k, v, logEntry[k])
				}
			}

			// Check RequestID present in log
			if _, ok := logEntry["request_id"]; !ok {
				t.Error("expected request_id in log")
			}
		})
	}
}

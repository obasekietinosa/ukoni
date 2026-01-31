package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type key int

const (
	requestLogDataKey key = iota
)

// RequestLogData holds data that can be populated by other middleware/handlers
// to be included in the request log.
type RequestLogData struct {
	UserID string
}

// GetRequestLogData retrieves the RequestLogData from the context.
func GetRequestLogData(ctx context.Context) *RequestLogData {
	if v, ok := ctx.Value(requestLogDataKey).(*RequestLogData); ok {
		return v
	}
	return nil
}

type LoggingMiddleware struct {
	Logger *slog.Logger
}

func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{Logger: logger}
}

// responseWriter wraps http.ResponseWriter to capture status and size
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int64
	body   []byte // Capture partial body for error details
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// If WriteHeader hasn't been called, it defaults to 200
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)

	// Capture body if status >= 400 (errors) and buffer isn't full
	if rw.status >= 400 && len(rw.body) < 512 {
		appendLen := 512 - len(rw.body)
		if appendLen > len(b) {
			appendLen = len(b)
		}
		rw.body = append(rw.body, b[:appendLen]...)
	}

	return size, err
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := uuid.New().String()

		// Add RequestID to header
		w.Header().Set("X-Request-ID", reqID)

		// Setup Log Data
		logData := &RequestLogData{}
		ctx := context.WithValue(r.Context(), requestLogDataKey, logData)

		rw := &responseWriter{ResponseWriter: w}

		// Process request
		next.ServeHTTP(rw, r.WithContext(ctx))

		duration := time.Since(start)

		// Determine log level and attributes
		attrs := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("request_id", reqID),
			slog.Int("status", rw.status),
			slog.Duration("duration", duration),
			slog.String("ip", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		}

		if logData.UserID != "" {
			attrs = append(attrs, slog.String("user_id", logData.UserID))
		}

		if rw.status >= 500 {
			if len(rw.body) > 0 {
				attrs = append(attrs, slog.String("error_details", string(rw.body)))
			}
			m.Logger.Error("request failed", attrs...)
		} else if rw.status >= 400 {
			if len(rw.body) > 0 {
				attrs = append(attrs, slog.String("error_details", string(rw.body)))
			}
			m.Logger.Warn("request failed", attrs...)
		} else {
			m.Logger.Info("request completed", attrs...)
		}
	})
}

package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

// statusColor returns the ANSI color for a given HTTP status code
func statusColor(status int) string {
	switch {
	case status >= 500:
		return colorRed
	case status >= 400:
		return colorYellow
	case status >= 300:
		return colorCyan
	case status >= 200:
		return colorGreen
	default:
		return colorReset
	}
}

// methodColor returns the ANSI color for a given HTTP method
func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return colorBlue
	case http.MethodPost:
		return colorGreen
	case http.MethodPut:
		return colorYellow
	case http.MethodDelete:
		return colorRed
	case http.MethodPatch:
		return colorCyan
	default:
		return colorReset
	}
}

// RequestLogger returns a slog-based request logging middleware
// In debug mode, it outputs colored human-readable logs
// In production mode, it outputs structured JSON logs
func RequestLogger(debugMode bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			status := wrapped.Status()

			if debugMode {
				// Human-readable colored output for local development
				reqID := middleware.GetReqID(r.Context())
				reqIDStr := ""
				if reqID != "" {
					reqIDStr = fmt.Sprintf(" %s[%s]%s", colorGray, reqID, colorReset)
				}

				fmt.Printf("%s%s%s %s%3d%s %s%-7s%s %s %s%s%s%s\n",
					colorGray, time.Now().Format("15:04:05"), colorReset,
					statusColor(status), status, colorReset,
					methodColor(r.Method), r.Method, colorReset,
					r.URL.Path,
					colorGray, duration, colorReset,
					reqIDStr,
				)
			} else {
				// Structured JSON logging for production
				level := slog.LevelInfo
				if status >= 500 {
					level = slog.LevelError
				} else if status >= 400 {
					level = slog.LevelWarn
				}

				slog.Log(r.Context(), level, "HTTP request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", status,
					"duration", duration.String(),
					"ip", r.RemoteAddr,
				)
			}
		})
	}
}

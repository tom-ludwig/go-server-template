package middleware

import (
	"io"
	"log"
	"net/http"
	"os"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// simpleLogger implements chi's LoggerInterface
type simpleLogger struct {
	logger *log.Logger
}

func (l *simpleLogger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

// ColoredLogger returns a chi logger middleware with colors in debug mode
func ColoredLogger(debugMode bool) func(next http.Handler) http.Handler {
	logger := &simpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	// Use colored logger for debug mode, plain for production
	return chimiddleware.RequestLogger(&chimiddleware.DefaultLogFormatter{
		Logger:  logger,
		NoColor: !debugMode,
	})
}

// ColoredLoggerWithWriter returns a chi logger middleware with colors in debug mode and custom writer
func ColoredLoggerWithWriter(w io.Writer, debugMode bool) func(next http.Handler) http.Handler {
	logger := &simpleLogger{
		logger: log.New(w, "", log.LstdFlags),
	}

	// Use colored logger for debug mode, plain for production
	return chimiddleware.RequestLogger(&chimiddleware.DefaultLogFormatter{
		Logger:  logger,
		NoColor: !debugMode,
	})
}

// Package httpclient provides an HTTP client preconfigured for distributed
// tracing.
package httpclient

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// New returns an *http.Client whose transport emits OTel client spans and
// propagates trace context.
func New() *http.Client {
	return &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
}

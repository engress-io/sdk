package trace

import (
	"net/http"
	"os"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Init starts the Datadog APM tracer when DD_AGENT_HOST is set.
func Init(service string) {
	if os.Getenv("DD_AGENT_HOST") == "" {
		return
	}
	env := os.Getenv("DD_ENV")
	if env == "" {
		env = "production"
	}
	tracer.Start(
		tracer.WithService(service),
		tracer.WithEnv(env),
	)
}

// Stop flushes and stops the tracer.
func Stop() {
	tracer.Stop()
}

// WrapHandler wraps an http.Handler with Datadog tracing.
func WrapHandler(service, resource string, handler http.Handler) http.Handler {
	if os.Getenv("DD_AGENT_HOST") == "" {
		return handler
	}
	return httptrace.WrapHandler(handler, service, resource)
}

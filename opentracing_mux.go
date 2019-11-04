package opentracing_mux

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/opentracing/opentracing-go"
)

const (
	UnknownMuxRoute = "UNKNOWN"
)

func TracingMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var span opentracing.Span
		name := getName(r)

		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header),
		)

		if err != nil {
			// TODO do something with the errors, in particular opentracing.ErrorSpanContextNotFound
		}

		span = opentracing.StartSpan(
			name,
			ext.RPCServerOption(wireContext),
		)
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(r.Context(), span)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getName(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if name := route.GetName(); name != "" {
		return name
	}
	if regex, err := route.GetPathRegexp(); err != nil {
		return regex
	}
	return UnknownMuxRoute
}

# omux

A simple opentracing middleware for gorilla/mux.

This middleware automatically reads incoming HTTP headers and generates a Span from them if opentracing headers are present. If not headers are found, it starts a new Trace.

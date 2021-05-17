package httptransport

import (
	"testing"

	"github.com/openzipkin/zipkin-go"
)

func TestNewTransport(t *testing.T) {
	tracer, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
	transport, err := NewTransport(tracer)
	if err != nil {
		t.Error(err)
	}

	if transport == nil {
		t.Error("expected transport non nil")
	}
}

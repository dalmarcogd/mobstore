package httptransport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

func NewTransport(tracer *zipkin.Tracer) (http.RoundTripper, error) {
	return zipkinhttp.NewTransport(
		tracer,
		zipkinhttp.TransportTrace(true),
		zipkinhttp.RoundTripper(&http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}),
	)
}

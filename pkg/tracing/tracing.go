package tracing

import (
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitTracer инициализирует трассировщик OpenTracing
func InitTracer(serviceName string) (opentracing.Tracer, io.Closer) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Fatalf("cfg.NewTracer: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}

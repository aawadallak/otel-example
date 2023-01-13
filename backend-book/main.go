package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aawadallak/otlp-example/pkg/app/books"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider(ctx context.Context) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "localhost:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(time.Second))),
	)

	global.SetMeterProvider(provider)

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

// HTTP middleware setting a value on the request context
func metricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		meter := global.Meter("metric_middleware")

		counter, err := meter.SyncInt64().Counter(
			"test_my_counter",
			instrument.WithDescription("Just a test counter"),
		)
		if err != nil {
			panic(err)
		}

		// Increment the counter.
		counter.Add(ctx, 1, attribute.String("foo", "bar"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HTTP middleware setting a value on the request context
func tracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("test-tracer")

		ctx, span := tracer.Start(r.Context(), r.RequestURI)
		defer span.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	ctx := context.Background()

	shutdown, err := initProvider(ctx)
	if err != nil {
		panic(err)
	}

	//nolint:errcheck
	defer shutdown(ctx)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(tracingMiddleware)
	r.Use(metricMiddleware)
	r.Use(middleware.Logger)

	r.Get("/books", books.GetBooks)

	http.ListenAndServe(":5000", r)
}

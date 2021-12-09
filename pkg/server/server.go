package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func InitServer(port *string) {
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.Use(otelmux.Middleware("otel-jaeger-demo-api"))

	router.HandleFunc("/", handleHelloWorld).
		Methods(http.MethodGet)
	router.HandleFunc("/error", handleError).
		Methods(http.MethodGet)

	fmt.Printf("Listening on port %s\n", *port)
	err := http.ListenAndServe(*port, router)
	if err != nil {
		panic(err)
	}
}

func handleError(w http.ResponseWriter, r *http.Request) {
    tp := otel.GetTracerProvider()
    tc := tp.Tracer("handleHelloWorld")
    _, span1 := tc.Start(r.Context(), "slow call 1")
    defer span1.End()

    err := errors.New("Fatal error")
    span1.RecordError(err)
    span1.SetStatus(codes.Error, err.Error())

    w.WriteHeader(http.StatusInternalServerError)
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) {
    tp := otel.GetTracerProvider()
    tc := tp.Tracer("handleHelloWorld")
    ctx, span1 := tc.Start(r.Context(), "slow call 1")
    defer span1.End()

    time.Sleep(time.Microsecond * 400)

    ctx, span2 := tc.Start(ctx, "slow call 2")
    defer span2.End()

    time.Sleep(time.Microsecond * 800)

	w.Write([]byte("Hello World"))
}

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("otel-jaeger-demo-api"),
		)),
	)
	return tp, nil
}

func ConfigOtelJaeger() *tracesdk.TracerProvider {
	tp, err := tracerProvider("http://jeager:14268/api/traces")
	if err != nil {
        panic(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}

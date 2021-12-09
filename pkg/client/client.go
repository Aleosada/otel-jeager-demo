package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Run(requests, interval *int, withError string) {
	cli := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	wg := &sync.WaitGroup{}
	wg.Add(*requests)
	for i := 1; i <= *requests; i++ {
		go func(requestId int, cli http.Client, wg *sync.WaitGroup) {
			ctx := context.Background()
			tp := otel.GetTracerProvider()
			tc := tp.Tracer("handleHelloWorld")
			ctx, span1 := tc.Start(ctx, "Request")
			defer span1.End()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:3000/"+withError, nil)
			if err != nil {
                span1.RecordError(err)
                span1.SetStatus(codes.Error, err.Error())
				wg.Done()
				panic(err)
			}

			res, err := cli.Do(req)
			if err != nil {
                span1.RecordError(err)
                span1.SetStatus(codes.Error, err.Error())
				wg.Done()
				panic(err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
                span1.RecordError(err)
                span1.SetStatus(codes.Error, err.Error())
				wg.Done()
				panic(err)
			}

			fmt.Printf("Request %d: %s\n", requestId, string(body))
			wg.Done()
		}(i, cli, wg)

		if *requests > 1 && *interval > 1 {
			time.Sleep(time.Second * time.Duration(*interval))
		}
	}
	wg.Wait()
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
			semconv.ServiceNameKey.String("otel-jaeger-demo-cli"),
		)),
	)
	return tp, nil
}

func ConfigOtelJaeger() *tracesdk.TracerProvider {
	tp, err := tracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		panic(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}

package opentelemetry

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	serviceName_gin = "gin-demo"
)

var tracer = otel.Tracer("gin-server")

func TestGin(t *testing.T) {
	ctx := context.Background()
	// 初始化tracer
	tp, err := initTracer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			t.Fatal(err)
		}
	}()
	r := gin.New()
	// 设置 otelgin 中间件
	r.Use(otelgin.Middleware(serviceName_gin))
	// 在响应头添加 Trace-Id
	r.Use(func(c *gin.Context) {
		c.Header("Trace-Id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String())
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := getUser(c, id)
		c.JSON(http.StatusOK, gin.H{
			"name": name,
			"id":   id,
		})
	})
	r.Run(":8080")
}

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tp, err := newJaegerTraceProvider(ctx, serviceName_gin)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}

func getUser(c *gin.Context, id string) string {
	_, span := tracer.Start(
		c.Request.Context(), "getUser", trace.WithAttributes(attribute.String("id", id)),
	)
	defer span.End()

	name := "unknown"
	if id == "1" {
		name = "Alice"
	} else if id == "2" {
		name = "Bob"
	}
	return name
}

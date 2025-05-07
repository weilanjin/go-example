package opentelemetry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// 实现http client的链路追踪

// Http client
const (
	serviceName_http = "httpClient-Demo"
	peerServiceName  = "baidu"
	blogURL          = "https://baidu.com"
)

func TestHttpClient(t *testing.T) {
	// init
	ctx := context.Background()
	tp, err := newJaegerTraceProvider(ctx, serviceName_http)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	//  --------------------------------------

	tr := otel.Tracer("http-client")
	// 开启 span，PeerService 指要连接的目标服务
	ctx, span := tr.Start(ctx, "baidu", trace.WithAttributes(semconv.PeerService(peerServiceName)))
	defer span.End()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, blogURL, nil)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	// 发起请求
	res, _ := client.Do(req)
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	fmt.Println("Response Body: ", len(body))
}

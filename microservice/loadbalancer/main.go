package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var servers = []Server{
	{
		ID:              "1",
		Name:            "Server 1",
		Protocol:        "HTTP",
		Host:            "localhost",
		Port:            8080,
		URL:             "http://localhost:8080",
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
	}, {
		ID:              "2",
		Name:            "Server 2",
		Protocol:        "HTTP",
		Host:            "localhost",
		Port:            8081,
		URL:             "http://localhost:8081",
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
	}, {
		ID:              "3",
		Name:            "Server 3",
		Protocol:        "HTTP",
		Host:            "localhost",
		Port:            8082,
		URL:             "http://localhost:8082",
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
	},
}

func main() {
	// 创建服务器池
	serverPool := NewServerPool()
	for _, server := range servers {
		serverPool.AddServer(&server)
	}

	// 创建负载均衡器
	rrb := NewRoundRobinBalancer(serverPool)
	lb := NewLoadBalancer(rrb)

	distributeLoad(8080, lb)

	select {}
}

func distributeLoad(port int, lb *LoadBalancer) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", lb.Serve)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Starting load balancer on port %d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}

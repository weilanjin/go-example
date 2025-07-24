package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HealthChecker struct {
	interval time.Duration
	pool     *ServerPool
}

func NewHealthChecker(interval time.Duration, pool *ServerPool) *HealthChecker {
	return &HealthChecker{
		interval: interval,
		pool:     pool,
	}
}

func (hc *HealthChecker) Check() {
	servers := hc.pool.GetAllServers()
	for {
		log.Println("Starting Health Check")
		for _, server := range servers {
			if server.HealthCheckURL != "" {
				doHttpRequest(server)
			} else {
				log.Printf("Health Check URL not specified for server: %s\n", server.URL)
			}
		}
		fmt.Println()
		fmt.Println()
		log.Println("Server Status")
		time.Sleep(hc.interval)
	}
}
func doHttpRequest(server *Server) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, server.HealthCheckURL, nil)
	if err != nil {
		updateServerUnhealthyStatus(server)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		updateServerUnhealthyStatus(server)
		log.Printf("Health check failed for server %s: %v\n", server.Name, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		updateServerUnhealthyStatus(server)
		log.Printf("Health check failed for server %s: %s\n", server.Name, resp.Status)
		return
	}
	if !server.IsHealthy {
		updateServerHealthyStatus(server)
	}
}

func updateServerUnhealthyStatus(server *Server) {
	log.Printf("Health Check Failed: Server did not responded with 200 status code. Server: %v\n", server.HealthCheckURL)

	if server.IsHealthy {
		server.FailureCount++
	}

	if server.FailureCount >= server.UnhealthyAfter && server.IsHealthy {
		server.IsHealthy = false
	}
}

func updateServerHealthyStatus(server *Server) {
	server.SuccessCount++

	if server.SuccessCount >= server.HealthyAfter && !server.IsHealthy {
		server.IsHealthy = true
		server.SuccessCount = 0
		server.FailureCount = 0
	}
}

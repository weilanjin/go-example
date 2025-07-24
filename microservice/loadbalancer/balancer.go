package main

import (
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BalancerStrategy interface {
	GetNextServer() (*Server, error) // 获取下一个服务器
}

type LoadBalancer struct {
	strategy BalancerStrategy // 负载均衡策略
}

func NewLoadBalancer(strategy BalancerStrategy) *LoadBalancer {
	return &LoadBalancer{
		strategy: strategy,
	}
}

func (lb *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	server, err := lb.strategy.GetNextServer()
	if err != nil {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}
	proxyURL, err := url.Parse(server.URL)
	if err != nil {
		http.Error(w, "Invalid server URL", http.StatusInternalServerError)
		return
	}

	proxyPath := strings.TrimRight(proxyURL.String(), "/") + r.URL.Path
	if r.URL.RawQuery != "" {
		proxyPath += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, proxyPath, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header
	req.Header.Set("X-Forwarded-For", r.RemoteAddr) // 转发原始客户端地址
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	byteResp, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	maps.Copy(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	fmt.Fprint(w, string(byteResp))
}

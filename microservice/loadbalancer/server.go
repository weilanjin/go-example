package main

import "time"

type Server struct {
	ID              string        // 服务器唯一标识
	Name            string        // 服务器人类可读的名称
	Protocol        string        // 协议（HTTP, HTTPS, TCP等）
	Host            string        // 服务器地址
	Port            uint16        // 端口号
	URL             string        // 包含协议、主机和端口的完整URL
	HealthCheckURL  string        // 健康检查的URL
	Timeout         time.Duration // 请求超时时间
	IsHealthy       bool          // 健康状态
	LastHealthCheck time.Time     // 上次健康检查时间
	FailureCount    int           // 失败计数
	SuccessCount    int           // 成功计数
	HealthyAfter    int           // 健康状态持续时间（秒）
	UnhealthyAfter  int           // 不健康状态持续时间（秒）
	RetryCount      int           // 重试次数
}

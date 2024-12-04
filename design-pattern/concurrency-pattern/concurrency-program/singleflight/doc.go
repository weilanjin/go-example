// Package singleflight
// 相同key合并执行一次
// sync.Once 主要被应用在单次初始的场景中
// SingleFlight 主要被应用在合并并发请求的场景中
package singleflight
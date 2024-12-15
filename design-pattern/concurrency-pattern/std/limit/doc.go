// Package limit
// 🪣令牌桶 、🧺漏桶
// 网络流量整形(traffic shaping) 和速率限制(rate limiting) 常用的算法
// [漏桶算法]能够强行限制请求的处理速率, 任何突发请求都会被平滑处理.
// [令牌桶算法]能够在限制请求处理速率的同时允许某种程度的突发请求.
package limit
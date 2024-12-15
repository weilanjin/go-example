package project

import (
	"errors"
	"fmt"
	"github.com/sony/gobreaker/v2"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

// Circuit Breaker Pattern | Microsoft Learn https://learn.microsoft.com/en-us/previous-versions/msp-n-p/dn589784(v=pandp.10)?redirectedfrom=MSDN
// pkg sony/gobreaker https://github.com/sony/gobreaker
//
// 断路器 circuit breaker
// 分布式系统中常用故障保护机制, 用于防止在服务调用过程中出现故障. 从而保护系统的稳定性和可靠性.
// 当下游服务器出现故障时, 断路器会自动打开, 并拒绝新的请求, 防止故障扩撒, 保护系统的其他部分
// 不受影响 等待下游服务器恢复后, 断路器会自动关闭, 开始接受新的请求.
// 断路器可以实现自我修复,在故障发生后一段时间内,断路器会尝试重新连接服务, 并检查服务是否恢复正常.
//
// 如果服务某个节点出现了负载过大问题, 则可能导致响应很慢,请求大量堆积在处理队列中,服务器的压力很大
// 后续请求还是被源源不断地发给这个节点, 造成业务堆积,处理变慢,整个系统的上下游很多节点处理都可能慢
// 下来, 服务器都有可能崩溃.

/*
	断路器有三种状态:
		闭合(Closed) - 所有的请求都会被放行
		断开(Open) - 所有的请求都不会被处理, 失败次数或者比例达到阈值 closed -> open
		半开(Half-Open) - 尝试处理请求, 处理请求成功次数或者比例达到阈值 half-open -> closed
*/

type State int

// 当读路器的状态发生或者断路处于闭合状态时, 清除这个计数
type Counts struct {
	Requests             uint32 // 请求数
	TotalSuccess         uint32 // 总成功请求数
	TotalFailure         uint32 // 总失败请求数
	ConsecutiveSuccesses uint32 // 连续成功请求数
	ConsecutiveFailures  uint32 // 连续失败请求数
}

type Settings struct {
	Name        string        // 断路器名称
	MaxRequests uint32        // 「半开状态」下允许通过的最大请求数. 如果值为0,则最多允许一个请求尝试通过
	Interval    time.Duration // 「闭合状态」下清除它的计数 Counts 的周期. 如果值为0,则不会清除计数 Counts
	Timeout     time.Duration // 处于「断开状态」下的时间,之后就会进入半开状态. 如果值为0,则表示断路器的过期时间为1分钟
	// 「闭合状态」一个请求失败后, 这个方法就会被调用, 并传入 Counts 作副本
	// 如果返回 true, 则表示进入断开状态
	ReadyToTrip   func(count Counts) bool
	OnStateChange func(name string, from, to State) // 状态发生改变时,会被调用
	IsSuccessful  func(err error) bool              // 默认, 非nil都返回false
}

var cb *gobreaker.CircuitBreaker[[]byte]

func TestCircuitBreaker(t *testing.T) {
	conf := gobreaker.Settings{
		Name: "HTTP GET",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// 计算失败比例
			failureRatio := float32(counts.TotalFailures) / float32(counts.Requests)
			// 断路器最多允许连续失败3次,并且失败比例超过60%
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	cb = gobreaker.NewCircuitBreaker[[]byte](conf)

	for i := 0; i < 10; i++ {
		_, err := HttpGet("https://www.baidu.com")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("xxx--xxx")
	}
}

func HttpGet(url string) ([]byte, error) {
	resp, err := cb.Execute(func() ([]byte, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(resp.Status)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
	return resp, err
}
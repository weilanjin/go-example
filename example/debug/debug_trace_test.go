package debug_test

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"lovec.wlj/example/broker"
	"lovec.wlj/example/debug"
	"lovec.wlj/example/debug/trace"
	"sync/atomic"
	"testing"
)

const spanUserRegister = "user_register"

func TestDebugTrace(t *testing.T) {
	debug.Trace("username1", spanUserRegister, "查询用户名是否存在", trace.WithErr(errors.New("查询用户信息失败")))
	debug.Trace("username1", spanUserRegister, "登记注册", trace.WithData("用户状态数据"))
	debug.Trace("username1", spanUserRegister, "管理员审核通过")
	debug.Trace("username1", spanUserRegister, "注册成功", trace.WithData("用户信息"))
}

type connector struct {
	status atomic.Bool
	mq     broker.Broker
}

func (conn *connector) Init(_ context.Context) error {
	conn.status.Store(true)
	conn.mq = broker.NewRedisPubSub(redis.NewClient(&redis.Options{}))
	return nil
}

func (conn *connector) Push(ctx context.Context, data ...*trace.TraceLog) error {
	if len(data) == 0 {
		return nil
	}
	var dataAny []any
	for _, v := range data {
		dataAny = append(dataAny, v)
	}

	return conn.mq.Publish("trace_log", dataAny, broker.PublishWithContext(ctx))
}

func (conn *connector) Enable() bool {
	return conn.status.Load()
}

func (conn *connector) Logger(err error) {
	log.Println(err)
}

package client_proxy

import (
	"lovec.wlj/go-case/rpc/simple2/handler"
	"net/rpc"
)

type CalcServiceStub struct {
	*rpc.Client
}

func NewCalcServiceClient(proto, address string) *CalcServiceStub {
	conn, err := rpc.Dial(proto, address)
	if err != nil {
		panic(err)
	}
	return &CalcServiceStub{conn}
}

func (c *CalcServiceStub) Add(request int, reply *int) error {
	err := c.Call(new(handler.CalcService).ServiceName()+".Add", request, reply)
	return err
}

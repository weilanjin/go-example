package server_proxy

import (
	"net/rpc"
)

type ICalcService interface {
	Add(request int, reply *int) error
	ServiceName() string
}

func RegisterCalcService(svc ICalcService) error {
	return rpc.RegisterName(svc.ServiceName(), svc)
}

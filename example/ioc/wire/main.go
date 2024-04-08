package main

import "lovec.wlj/example/ioc/wire/service"

type Application struct {
	svc *service.ABCDService
}

func NewApplication(svc *service.ABCDService) *Application {
	return &Application{
		svc: svc,
	}
}

func main() {
	app := wireApp()
	_ = app
}

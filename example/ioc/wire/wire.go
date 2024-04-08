//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"lovec.wlj/example/ioc/wire/repo"
	"lovec.wlj/example/ioc/wire/service"
)

func wireApp() *Application {
	panic(wire.Build(repo.ProviderSet, service.ProviderSet, NewApplication))
}

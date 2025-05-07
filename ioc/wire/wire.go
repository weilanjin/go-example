//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/weilanjin/go-example/ioc/wire/repo"
	"github.com/weilanjin/go-example/ioc/wire/service"
)

func wireApp() *Application {
	panic(wire.Build(repo.ProviderSet, service.ProviderSet, NewApplication))
}
package service

import "lovec.wlj/example/ioc/wire/repo"

type ABCDService struct {
	a *repo.ARepo
	b *repo.BRepo
	c *repo.CRepo
	d *repo.DRepo
}

func NewABCDService(a *repo.ARepo, b *repo.BRepo, c *repo.CRepo, d *repo.DRepo) *ABCDService {
	return &ABCDService{a: a, b: b, c: c, d: d}
}

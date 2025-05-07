package service

import "github.com/weilanjin/go-example/err/repo"

func QueryUser() error {
	return repo.FindUser()
}
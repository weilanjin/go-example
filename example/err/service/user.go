package service

import "lovec.wlj/example/err/repo"

func QueryUser() error {
	return repo.FindUser()
}
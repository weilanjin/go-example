package lock_test

import (
	"testing"

	"lovec.wlj/go-case/lock"
)

func TestDeadlock(t *testing.T) {
	lock.Deadlock()
}

func TestLivelock(t *testing.T) {
	lock.Livelock()
}

func TestStarvation(t *testing.T) {
	lock.Starvation()
}

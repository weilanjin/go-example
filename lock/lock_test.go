package lock_test

import (
	"testing"

	"github.com/weilanjin/go-example/lock"
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
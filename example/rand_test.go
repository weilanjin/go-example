package test

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestRandInt(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 10; i++ {
		log.Println(rand.Intn(3) + 1)
	}
}

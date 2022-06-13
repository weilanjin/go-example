package concurrency_pattern

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func boring(msg string) {
	for i := 0; ; i++ {
		log.Println(msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func Test1(t *testing.T) {
	go boring("boring ..")
	log.Println("I'm listening")
	time.Sleep(time.Second * 2)
	log.Println("You're boring, I'm leaving")
}

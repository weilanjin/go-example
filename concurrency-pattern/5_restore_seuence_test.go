package concurrency_pattern

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

type Message struct {
	str  string
	wait chan struct{}
}

func faninMessage(ins ...<-chan Message) <-chan Message {
	c := make(chan Message)
	for _, in := range ins {
		go func(cm <-chan Message) {
			for {
				c <- <-cm
			}
		}(in)
	}
	return c
}

func sendMsg(msg string) <-chan Message {
	c := make(chan Message)
	wait := make(chan struct{})
	go func() {
		for i := 0; ; i++ {
			c <- Message{
				str:  fmt.Sprintf("%s %d", msg, i),
				wait: wait,
			}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-wait
		}
	}()
	return c
}

func Test5(t *testing.T) {
	c := faninMessage(sendMsg("Joe"), sendMsg("Ahn"))
	for i := 0; i < 5; i++ {
		msg1 := <-c
		log.Println(msg1.str)
		msg2 := <-c
		log.Println(msg2.str)
		msg1.wait <- struct{}{}
		msg2.wait <- struct{}{}
	}
	log.Println("You're both boring. I'm leaving")
}

package limit

import (
	"fmt"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	in := []int{3, 2, 1}
	timeout := 2
	chLimit := make(chan struct{}, 1)
	chs := make([]chan string, len(in))
	limitFn := func(chLimit chan struct{}, ch chan string, task_id, sleeptime, timeout int) {
		Run(task_id, sleeptime, timeout, ch)
		<-chLimit
	}
	startTime := time.Now()
	fmt.Println("Multirun start")
	for i, sleeptime := range in {
		chs[i] = make(chan string, 1)
		chLimit <- struct{}{}
		go limitFn(chLimit, chs[i], i, sleeptime, timeout)
	}
	for _, ch := range chs {
		fmt.Println(<-ch)
	}
	endTime := time.Now()
	fmt.Printf("Multissh finished, Process time %s, Number of task is %d", endTime.Sub(startTime), len(in))
}

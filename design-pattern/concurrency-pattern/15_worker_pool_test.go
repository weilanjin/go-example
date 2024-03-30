package concurrency_pattern

import (
	"log"
	"sync"
	"testing"
	"time"
)

func workerEfficient(id int, jobs <-chan int, results chan<- int) {
	var wg sync.WaitGroup
	for job := range jobs {
		wg.Add(1)
		go func(j int) {
			log.Println("worker", id, "started job", job)
			time.Sleep(time.Second)
			log.Println("worker", id, "finished job", job)
			results <- job * 2
			wg.Done()
		}(job)
	}
	wg.Wait()
}

func Test15(t *testing.T) {
	const numbJobs = 8
	jobs := make(chan int, numbJobs)
	results := make(chan int, numbJobs)
	for i := 0; i < 3; i++ {
		go workerEfficient(i, jobs, results)
	}
	for i := 0; i < numbJobs; i++ {
		jobs <- i
	}
	close(jobs)
	log.Println("Closed job")
	for i := 0; i < numbJobs; i++ {
		<-results
	}
	close(results)
}

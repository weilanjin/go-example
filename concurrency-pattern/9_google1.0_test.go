package concurrency_pattern

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

type Result string

type Search func(query string) Result

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %s\n", kind, query))
	}
}

func Google(query string) (results []Result) {
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return
}

func Test9(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("go")
	elapsed := time.Since(start)
	log.Println(results)
	log.Println(elapsed)
}

// ======= google2.0 =========

func Google2(query string) []Result {
	c := make(chan Result)
	go func() {
		c <- Web(query)
	}()
	go func() {
		c <- Image(query)
	}()
	go func() {
		c <- Video(query)
	}()
	var results []Result
	for i := 0; i < 3; i++ {
		results = append(results, <-c)
	}
	return results
}

func Test9_1(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google2("go")
	elapsed := time.Since(start)
	log.Println(results)
	log.Println(elapsed)
}

// ======= google2.1 =========

func Google2_1(query string) []Result {
	c := make(chan Result)
	go func() {
		c <- Web(query)
	}()
	go func() {
		c <- Image(query)
	}()
	go func() {
		c <- Video(query)
	}()
	var results []Result
	timeout := time.After(50 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case r := <-c:
			results = append(results, r)
		case <-timeout:
			log.Println("timeout")
			return results
		}
	}
	return results
}

func Test9_2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google2_1("go")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

// ======= google3 =========

var (
	Web1   = fakeSearch("web1")
	Image1 = fakeSearch("image1")
	Video1 = fakeSearch("video1")
)

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	for i := range replicas {
		go func(idx int) {
			c <- replicas[idx](query)
		}(i)
	}
	return <-c
}

func Google3(query string) []Result {
	c := make(chan Result)
	go func() {
		c <- First(query, Web1, Web)
	}()
	go func() {
		c <- First(query, Image, Image1)
	}()
	go func() {
		c <- First(query, Video, Video1)
	}()
	var results []Result
	after := time.After(50 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case r := <-c:
			results = append(results, r)
		case <-after:
			log.Println("timeout")
			return results
		}
	}
	return results
}

func Test9_3(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google3("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

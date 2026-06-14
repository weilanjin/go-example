package test

import "fmt"

func Hello() string {
	return "Hello, World!"
}

// go test -v -run ExampleHello
func ExampleHello() {
	fmt.Println(Hello())
	// Output: Hello, World!
}

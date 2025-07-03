package main

import (
	"fmt"
	"unique"
)

// https://go.dev/blog/unique
// https://medium.com/google-cloud/interning-in-go-4319ea635002
// https://victoriametrics.com/blog/go-unique-package-intern-string/

func main() {
	h1 := unique.Make("Hello")
	h2 := unique.Make("Hello")
	w1 := unique.Make("World")

	fmt.Println("h1:", h1)
	fmt.Println("h2:", h2)
	fmt.Println("w1:", w1)
	fmt.Println("h1 == h2:", h1 == h2)
	fmt.Println("h1 == w1:", h1 == w1)
}

// Output:
// h1: {0x14000090030}
// h2: {0x14000090030}
// w1: {0x14000090040}
// h1 == h2: true
// h1 == w1: false

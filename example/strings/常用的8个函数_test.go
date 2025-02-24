package strings

import (
	"fmt"
	"strings"
	"testing"
)

// 8 Most used function of strings package in Go

// 1. Compare
func TestCompare(t *testing.T) {
	fmt.Println(strings.Compare("gopher19", "gopher20"))
	fmt.Println(strings.Compare("gopher19", "gopher19"))
	fmt.Println(strings.Compare("gopher20", "gopher19"))
	// Output:
	// -1
	// 0
	// 1
}

// 2. Contains
func TestContains(t *testing.T) {
	fmt.Println(strings.Contains("gopher", "go"))
	fmt.Println(strings.Contains("gopher", "php"))
	fmt.Println(strings.Contains("gopher", ""))
	fmt.Println(strings.Contains("", ""))
	// Output:
	// true
	// false
	// true
	// true
}

// 3. Count
func TestCount(t *testing.T) {
	fmt.Println(strings.Count("gopher-gopher-gopher", "g"))
	fmt.Println(strings.Count("gopher", "g"))
	// Output:
	// 3
	// 1
}

// 4. Fields
func TestField(t *testing.T) {
	fmt.Printf("%q\n", strings.Fields("gopher is sleeping"))
	// Output:
	// ["gopher" "is" "sleeping"]
}

// 5. HasPrefix
func TestHasPrefix(t *testing.T) {
	fmt.Println(strings.HasPrefix("gopher", "go"))
	fmt.Println(strings.HasPrefix("gopher", "php"))
	fmt.Println(strings.HasPrefix("gopher", ""))
	// Output:
	// true
	// false
	// true
}

// 6. Replace
func TestReplace(t *testing.T) {
	fmt.Println(strings.Replace("go go gopher", "go", "golang", 2))
	fmt.Println(strings.Replace("go go gopher", "go", "golang", -1))
	// Output:
	// golang golang gopher
	// golang golang golangpher
}

// 7. Split
func TestSplit(t *testing.T) {
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))
	fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", " "))
	// Output:
	// ["a" "b" "c"]
	// ["a" "man" "a" "plan" "a" "canal" "panama"]
}

// 8. TrimSpace
func TestTrimSpace(t *testing.T) {
	fmt.Println(strings.TrimSpace(" \t\n hello world \n\t\r"))
	// Output:
	// hello world
}
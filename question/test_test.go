package question

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	sn1 := struct {
		age  int
		name string
	}{age: 11, name: "qq"}
	sn2 := struct {
		age  int
		name string
	}{age: 11, name: "qq"}
	if sn1 == sn2 {
		fmt.Println("sn1 == sn2")
	}
}
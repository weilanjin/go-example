package question

import (
	"errors"
	"fmt"
	"testing"
)

var f = func(i int) {
	println("x")
}

func TestVariable(t *testing.T) {
	f := func(i int) {
		print(i)
		if i > 0 {
			f(i - 1) // 内部 f 声明时, 此处的 f 调用全局的f
		}
	}
	f(10)

	/*
			var x = 23
		　　func main() {
		　　	   x := 2*x-4
		   }
	*/
}

// output: 10x

func TestVar1(t *testing.T) {
	a := 1
	for i := 0; i < 3; i++ {
		// for 语句的变量a是重新声明, 它的作用范围只在for语句范围内
		a := a + 1
		a = a * 2
	}
	fmt.Println(a) // 1
}

var ErrDidNotWork = errors.New("did not work")

// 变量作用域
func DoTheThing(reallyDoIt bool) (err error) {
	if reallyDoIt {
		// if 语句块内的err变量会遮罩函数作用域内的err变量
		result, err := tryTheThing()
		if err != nil || result != "it works" {
			err = ErrDidNotWork
		}
	}
	return err
}

func tryTheThing() (string, error) {
	return "", ErrDidNotWork
}

func TestVar2(t *testing.T) {

}

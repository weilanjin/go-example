package question

import (
	"fmt"
	"testing"
)

// defer 的执行顺序是（栈）后进先出
// panic 语句出现时，会先执行defer，最后才执行panic
func TestDeferCall(t *testing.T) {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	panic("触发异常")
}

// output：
// 		打印后
// 		打印中
// 		打印前
//		panic: 触发异常

func df(i int) {
	fmt.Println(i)
}

func TestDeferFunc(t *testing.T) {
	i := 5
	// 在执行 defer 语句时时候会保存一份副本。
	defer df(i) // 5
	i += 5
}

// ------------------

func increaseA() int {
	var i int
	defer func() {
		i++
	}()
	return i
}

func increaseB() (r int) {
	defer func() {
		r++
	}()
	return r
}

func TestAnonymous(t *testing.T) {
	fmt.Println(increaseA()) // 0 返回参数匿名参数
	fmt.Println(increaseB()) // 1
}

// ----------------------------

type Person struct {
	age int
}

func TestDefer1(t *testing.T) {
	person := &Person{28}

	// 3. person.age 将28当作defer函数的参数，会把28缓存在栈中。执行defer语句时取出28
	defer fmt.Println(person.age)

	// 2.defer 缓存结构体 Person{28}的地址，Person{28}的age被重新赋值为29， defer语句最后执行时，取出age便是29
	defer func(p *Person) {
		fmt.Println(p.age)
	}(person)

	// 1.闭包引用
	defer func() {
		fmt.Println(person.age) // 29
	}()
	person.age = 29
}

// --------------------

func calc(index string, a, b int) int {
	res := a + b
	fmt.Println(index, a, b, res)
	return res
}

func TestDefer2(t *testing.T) {
	a, b := 1, 2
	defer calc("1", a, calc("10", a, b))
	a = 0
	defer calc("2", a, calc("20", a, b))
	b = 1
}

// output:
//  10 1 2 3
//  20 0 2 2
//  2  0 2 2
//  1  1 3 4

func F(n int) func() int {
	return func() int {
		n++
		return n
	}
}

func TestDefer3(t *testing.T) {
	f := F(5)
	defer func() {
		fmt.Println(f()) // 3. n = 8
	}()
	// defer() 后面的函数如果带参数，会优先计算参数，并将结果存储在栈中, 到真正执行defer()的时候取出。
	defer fmt.Println(f()) // 1. n = 6
	i := f()
	fmt.Println(i) // 2. n = 7
}

// output:
//  7 6 8

// -------------------
func deferTest1(i int) (r int) {
	r = i
	defer func() {
		r += 3 // 1 + 3
	}()
	return r // 1
}

func deferTest2(i int) (r int) {
	defer func() {
		r += 1 // 2 + 1
	}()
	return 2
}

func TestDefer4(t *testing.T) {
	fmt.Println(deferTest1(1)) // 4
	fmt.Println(deferTest2(1)) // 3
}

// ------------------

type Slice []int

func NewSlice() Slice {
	return make(Slice, 0)
}

func (s *Slice) Add(elem int) *Slice {
	*s = append(*s, elem)
	fmt.Println(elem)
	return s
}

func TestDefer5(t *testing.T) {
	s := NewSlice()
	//  defer 先执行第一个Add, 第二个放入栈执行缓存 return 后执行
	defer s.Add(1).Add(2)

	defer func() {
		s.Add(11).Add(22) // return 后执行
	}()
	s.Add(3)
}

// output: 1 3 11 12 2

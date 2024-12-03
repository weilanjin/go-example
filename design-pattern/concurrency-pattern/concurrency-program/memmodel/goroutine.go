package memmodel

// 在父goroutine中启动子goroutine的go语句的执行,一定synchronized before 子 goroutine 中代码的执行

var x string
var y int

func fg() {
	print(x)
	y = 1
}

func goroutine() {
	x = "hello world"
	go fg()
	print(y) // 可能是 0 也可能是 1
}
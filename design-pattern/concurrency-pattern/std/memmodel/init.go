package memmodel

// init 函数一定 happens before main.main 函数

var (
	a = c + b // 9
	b = f()   // 4
	c = f()   // 5
	d = 3     // 5
)

func f() int {
	d++
	return d
}
package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
	"unsafe"
)

// 参考地址 https://colobu.com/gotips/001.html

// 1. SliceToArray 切片转数组
// -----------------------------------------------
// Go 1.20
var Sli = []int{1, 2, 3}
var Arr = [2]int(Sli[:2]) // [1, 2]

// Go 1.17
var Arr2 = *(*[2]int)(Sli[:2]) // [1, 2]

// 2. 编译时接口检查
// -----------------------------------------------
// var _ Buffer = (*StringBuffer)(nil)

// 3. 数组分隔符
// -----------------------------------------------
const (
	OneBillion  = 1e9
	OneBillion1 = 1_000_000_000
	PI          = 3.14_159_265_358
)

// 4. 省略 getter 方法的Get前缀
// -----------------------------------------------
type Person struct {
	name string
}

func (p *Person) Name() string {
	return p.name
}
func (p *Person) SetName(name string) {
	p.name = name
}

// 5. Must 有意的来停止程序 失败会panic
// -----------------------------------------------
var ValidID = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
var Tmpl = template.Must(template.New("name").Parse(`Hello {{.}}!`))

// 6. 不可比较的结构体 比如 float
// -----------------------------------------------
type Point struct {
	// 非导出的: 对于你的结构体的使用者来说是隐藏的
	// 零宽度(或无成本): 因为长度为0,所以这个数组在内存中不占用任何空间
	// 不可比较: func()是一个函数类型,而函数在Go中是不可比较的
	_    [0]func()
	X, Y float64
}

func (p *Point) Equals(other Point, tolerance float64) bool {
	return math.Abs(p.X-other.X) < tolerance &&
		math.Abs(p.Y-other.Y) < tolerance
}

// 7. 针对容器环境 k8s、docker 调整GOMAXPROCS
// -----------------------------------------------
// Kubernetes CPU限制和Go https://www.ardanlabs.com/blog/2024/02/kubernetes-cpu-limits-go.html
// import _ "go.uber.org/automaxprocs"
// func main() {
// }

// 8. 将互斥放在保护的数据附近
// -----------------------------------------------
// 结构体
type UserSession struct {
	ID        string
	LastLogin time.Time

	mu         sync.Mutex // 要保护那些字段可以一目了然 Preference、Cart
	Preference map[string]any
	Cart       []string

	IsLoggedIn bool
}

// 全局变量
var (
	mu    sync.Mutex
	count int
)

// 9. 使用strings.EqualFold忽略大小写比较
// -----------------------------------------------

func init() {
	str1 := "Hello"
	str2 := "hello"
	strings.EqualFold(str1, str2)
}

// 10. 结构体中的字段按照从大到小的顺序排列
// -----------------------------------------------
// 字段填充,内存对齐
// https://github.com/dkorunic/betteralign?tab=readme-ov-file 可以加检查低效的对齐方式

// 32 字节对齐
type T1 struct {

	//  +---++++
	//  +-------
	//  ++++++++
	//  +-------

	A byte  // 1 字节 + 3  	// 8 字节对齐
	B int32 // 4 字节

	C byte // 1 字节 + 7   // 8 字节对齐
	// 8 字节对齐
	D int64 // 8 字节       // 8 字节对齐
	// 8 字节对齐
	E byte // 1 字节 + 7   // 8 字节对齐
}

// 16 字节对齐 (优化)
type T2 struct {
	// ++++++++
	// +++++++-
	D int64 // 8 字节
	B int32 // 4 字节
	A byte  // 1 字节
	C byte  // 1 字节
	E byte  // 1 字节
}

func init() {
	fmt.Println(unsafe.Sizeof(T1{})) // 32 bytes
	fmt.Println(unsafe.Sizeof(T2{})) // 16 bytes
}

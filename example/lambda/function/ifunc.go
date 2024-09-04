// Package function 函数式接口
//
// Consumer 有参，无返回值
// Function 有参, 有返回值
// Predicate 有参，返回布尔值
// Supplier 无参，有返回值
// Operator 有参, 有返回值 且返回值与入参类型一致
//
// Unary 一元、Binary 二元、Ternary 三元、Quaternary 四元
package function

// Runnable 没有入参，没有返回值
type Runnable interface {
	Run()
}

// Callable 有入参，有返回值
type Callable[T any] interface {
	Call() T
}

// Comparator 比较两个值的大小
type Comparator[T any] interface {
	Compare(a, b T) int
}

// Consumer 消费一个值
type Consumer[T any] interface {
	Accept(value T)
}

// BiConsumer [binary consumer] 消费两个值
type BiConsumer[T1 any, T2 any] interface {
	Accept(a T1, b T2)
}

type Function[T any, R any] interface {
	Apply(value T) R
}

type BiFunction[T1 any, T2 any, R any] interface {
	Apply(a T1, b T2) R
}

// Predicate 断言 判断一个值是否满足条件
type Predicate[T any] interface {
	Test(value T) bool
}

type Supplier[T any] interface {
	Get() T
}

type UnaryOperator[T any] interface {
	Apply(value T) T
}

type BinaryOperator[T any] interface {
	Apply(a, b T) T
}

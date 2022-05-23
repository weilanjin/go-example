package singleton

import "sync"

// 通过该接口可以避免 GetInstance 返回一个包私有类型的指针
type Singleton interface {
	foo()
}

type singleton struct{}

func (s singleton) foo() {}

var (
	instance *singleton
	once     sync.Once
)

func GetInstance() Singleton {
	once.Do(func() {
		instance = &singleton{}
	})

	return instance
}

package project

// 活动对象模式
//   也叫 并发对象 Concurrency Object、Actor 设计模式
//
// 解耦了方法的调用和执行,使它们在不同的线程(或者纤程、goroutine)之中.
// 引入了异步方法调用,允许应用程序并发地处理多个客户端请求,通过调度器调用并发方法的执行,提供了并发执行方法的能力
/*
	活动对象模式包含6个组件

	- proxy: 定义了客户端请求接口.当客户端调用它的方法时,方法调用被转换成方法请求,放入scheduler的activation queue之中
	- method request: 用来封装方法调用的上下文
	- activation queue: 待处理的方法请求队列
	- scheduler: 一个独立线程, 管理activation queue, 调度方法的执行
	- servant: 活动对象的方法执行的具体实现
	- future: 当客户端调用方法时, 一个future对象会立即返回, 允许客户端获取返回结果
*/

type MethodRequest int

const (
	Increment MethodRequest = iota
	Decrement
)

type Service struct {
	queue chan MethodRequest
	v     int
}

func NewService() *Service {
	s := &Service{
		queue: make(chan MethodRequest),
	}
	go s.schedule()
	return s
}

func (s *Service) Increment() {
	s.queue <- Increment
}

func (s *Service) Decrement() {
	s.queue <- Decrement
}

func (s *Service) schedule() {
	for request := range s.queue {
		if request == Increment {
			s.v++
		} else {
			s.v--
		}
	}
}
package channel

import "reflect"

// 任务编排 Or-Done 模式
// 如果有多个任务, 只要任意一个任务执行完成,就可以返回任务完成的型号

// Or 递归实现
func Or(channels ...<-chan any) <-chan any {
	if len(channels) == 0 {
		closeChan := make(chan any)
		close(closeChan)
		return closeChan
	}
	if len(channels) == 1 {
		return channels[0]
	}
	orDone := make(chan any)
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			m := len(channels) / 2 // 超过两个, 二分法递归处理
			select {
			case <-Or(channels[:m]...):
			case <-Or(channels[m:]...):
			}
		}
	}()
	return orDone
}

// OrV2 不实用递归实现, 使用了反射方法 reflect.Select, 避免了深层递归
func OrV2(channels ...<-chan any) <-chan any {
	if len(channels) == 0 {
		closeChan := make(chan any)
		close(closeChan)
		return closeChan
	}
	if len(channels) == 1 {
		return channels[0]
	}
	orDone := make(chan any)
	go func() {
		defer close(orDone)
		var cases []reflect.SelectCase
		for _, ch := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}
		reflect.Select(cases) // 一旦有一个通道收到数据, reflect.Select 就会返回, 此时关闭返回的通道 orDone
	}()
	return orDone
}
package channel

// 扇入

// FanIn 扇入, 将多个channel的数据合并到一起
func FanIn[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	for _, c := range channels {
		go func(ch <-chan T) {
			for v := range ch {
				out <- v
			}
		}(c)
	}
	return out
}

func FanOut[T any](ch <-chan T, out []chan T, async bool) {
	go func() {
		defer func() {
			for _, c := range out {
				close(c)
			}
		}()
		for v := range ch {
			for _, c := range out {
				if async {
					go func(c chan T) {
						c <- v
					}(c)
				} else {
					c <- v
				}
			}
		}
	}()
}
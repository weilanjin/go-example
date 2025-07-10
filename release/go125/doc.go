package go125

/*

	1. gc 标记和扫描小对象性能提升，预计减少 0-40% GC开销
	2. testing/synctest 提供测试并发代码的支持（伪造时钟和goroutine等待机制）
	3. 容器感知 gomaxprocs 根据 cgroup 限制调整，并动态更新
	4. go build -acan 内存泄露检测 在程序退出时执行内存泄露检测。
	5. 减少预购工具二进制文件，减少Go发行版的大小
	6. panic 信息输出变更
	7. 新实验性 encoding/json/v2 包
	8. WaitGroup.Go 方法
*/

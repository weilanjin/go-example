// 解决并发-等待问题
// 等所有goroutine执行完毕后才继续执行
// 例如:
// Linux - barrier (屏障)
// C++ - std::barrier
// Java - CyclicBarrier 和 CountDownLatch
// ---
// https://github.com/sourcegraph/conc
package waitgroup

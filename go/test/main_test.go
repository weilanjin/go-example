package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// 1. 设置阶段：在所有测试运行前执行
	setup()

	// 2. 运行测试：m.Run() 执行包内所有的测试函数
	// m.Run() 返回一个整数，作为程序的退出码
	// 会触发该包内所有 TestXxx、BenchmarkXxx 和 ExampleXxx 函数的执行
	exitCode := m.Run()

	// 3. 清理阶段：在所有测试运行后执行
	teardown()

	// 4. 退出程序
	os.Exit(exitCode)
}

func setup()    {}
func teardown() {}

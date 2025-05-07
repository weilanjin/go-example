// 1. 所有以_test.go 为后缀的源码文件都会被 go test 命令执行.
// 2. _test.go 源码文件, go build 命令不会将这些测试文件打包到可执行文件中
// test 文件有4类: Test开头的功能测试、Benchmark开头的性能测试、Example、 Fuzz模糊测试
package test

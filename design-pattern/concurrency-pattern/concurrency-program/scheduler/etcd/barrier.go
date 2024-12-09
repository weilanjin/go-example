package main

import (
	"bufio"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"os"
)

// 分布式屏障

/*
	// 创建分布式 Barrier
	func NewBarrier(client *v3.Client, key string) *Barrier

	func (b *Barrier) Hold() error    // 创建一个屏障,实际上会创建一个key.若有节点调用它的Wait方法,会被阻塞
	func (b *Barrier) Release() error // 打开屏障,实际上删除这个key.所有等待阻塞的节点都会放行
	func (b *Barrier) Wait() error    // 阻塞当前节点,等屏障被释放,如果屏障没有创建就会直接放行
*/

func Barrier(cli *clientv3.Client) {
	b := recipe.NewBarrier(cli, "barrier")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "hold":
			b.Hold()
			fmt.Println("barrier hold")
		case "release":
			b.Release()
			fmt.Println("barrier release")
		case "wait":
			b.Wait()
			fmt.Println("after wait")
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action")
		}
	}
}
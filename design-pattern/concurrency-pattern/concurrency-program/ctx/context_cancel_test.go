package ctx

import (
	"context"
	"crypto/sha256"
	"log"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestContextCancel(t *testing.T) {
	Bitcoin()
}

// 模拟挖矿
func Bitcoin() {
	targetBits, _ := strconv.Atoi(os.Args[1])
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	// 消耗算力的挖矿算法
	pow := func(ctx context.Context, targetBits int, ch chan string) {
		target := big.NewInt(1)
		target.Lsh(target, uint(256-targetBits)) // 除了前 targetBits 位，其他位都是 1

		var hashInt big.Int
		var hash [32]byte
		nonce := 0 // 随机数

		// 寻找一个满足当前难度的数
		for {
			select {
			case <-ctx.Done():
				log.Println("context is cancelled")
				ch <- ""
				return
			default:
				data := "hello world" + strconv.Itoa(nonce)
				hash = sha256.Sum256([]byte(data)) // 计算hash值
				hashInt.SetBytes(hash[:])          // 将hash值转换为big.Int
				if hashInt.Cmp(target) < 1 {       // hashInt <= target, 找到一个不大于目标值的数,也就是至少前targetBits位都为0
					ch <- data
					return
				} else { // 没有找到继续找
					nonce++
				}
			}
		}
	}

	// 生成一个可撤销的Context
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan string, 1)
	go pow(ctx, targetBits, ch) // 子 goroutine 去挖矿

	time.Sleep(time.Second) // 等待1s

	select {
	case res := <-ch:
		log.Println("find the hash: ", res)
	default:
		cancel() // 撤销pow计算
		log.Println("没有找到比目标值小的数:", ctx.Err())
	}
}

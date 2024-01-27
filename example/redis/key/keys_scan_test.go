package key

import (
	"log"
	"testing"
)

// 问题 1：新增的 key 可能没有遍历到
// 问题 2：遍历出了重复键
//
// keys 	--> scan
// hgetall 	--> hscan
// smembers --> sscan
// zrange 	--> zscan

func TestKeys(t *testing.T) {

	// keys pattern
	// * -- 比配任意字符
	// . -- 比配一个字符
	// e.g keys *
	//	   keys user:[0-9]* # 查找以 "user:" 开头且包含数字的所有键
	// keys命令会遍历所有键，所以它的时间复杂度是O（n）, 并且将匹配模式的键一次性地返回给客户端
	//「生产环境禁止使用」
	str := rdb.Keys(ctx, "user*").String()
	log.Println(str)
}

func TestScan(t *testing.T) {

	// SCAN cursor match count [type]
	// e.g SCAN 0 MATCH user:* COUNT 10 # 使用 SCAN 迭代查找以 "user:" 开头的所有键，并每次返回 10 个键
	//
	// scan命令使用游标（cursor）方式来逐步迭代匹配的键。它将匹配的键分批返回，以减轻服务的阻塞和性能压力。
	// SCAN命令的算法是基于近似随机采样的算法。它使用一个游标来记录当前迭代的位置，并在每次迭代时返回一批键。
	// 迭代过程中，服务器会根据游标位置和一定的采样数量来获取匹配的键。
	iter := rdb.ScanType(ctx, 0, "*", 0, "string").Iterator()
	for iter.Next(ctx) {
		log.Println("key:", iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	// 集群模式
	// err := rdb.ForEachMaster(ctx, func(ctx context.Context, rdb *redis.Client) error {
	//		iter := rdb.Scan(ctx, 0, "prefix:*", 0).Iterator()
	//		for iter.Next(ctx) {
	//			log.Println("key:", iter.Val())
	//		}
	//
	//		return iter.Err()
	// })
	// if err != nil {
	//	 panic(err)
	// }
}

func TestScanZset(t *testing.T) {
	iter := rdb.ZScan(ctx, "zsetkey128", 0, "*", 10).Iterator()
	for iter.Next(ctx) {
		log.Println("key:", iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}

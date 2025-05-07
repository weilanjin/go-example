package basis

import (
	"log"
	"testing"
)

// hash
/*
							hash类型命令时间复杂度
	+-------------------------------------+----------------------+
	| 命令                                 | 时间复杂度            |
	+-------------------------------------+----------------------+
	| hset key field value                | O(1)                 |
	| hget key field                      | O(1)                 |
	| hdel key field [field … ]           | O(k), k是field个数    |
	| hlen key                            | O(1)                 |
	| hgetall key                         | O(n), n是field总数    |
	| hmget field [field ...]             | O(k), k是field的个数  |
	| hmset field value [field value ...] | O(k), k是field的个数  |
	| hexists key field                   | O(1)                 |
	| hkeys key                           | O(n), n是field总数    |
	| hvals key                           | O(n), n是field总数    |
	| hsetnx key field value              | O(1)                 |
	| hincrby key field increment         | O(1)                 |
	| hincrbyfloat key field increment    | O(1)                 |
	| hstrlen key field                   | O(1)                 |
	+-------------------------------------+----------------------+
*/

// 两种内部编码
//
//			1.ziplist: 压缩列表
//				hash类型元素的个数  hash-max-ziplist-entries < 512 or 所有的值hash-max-ziplist-value < 65。
//				更加紧凑的结构实现多个元素连续存储。
//		     	过多读写效率会变差。
//			2.hashtable: 哈希列表
//	      		hashtable读写时间复杂度为 O(1)
func TestHashObject(t *testing.T) {
	rdb.HMSet(ctx, "hashkey", "f1", "v1", "f2", "v2")
	log.Println(rdb.ObjectEncoding(ctx, "hashkey")) // listpack

	// 当 value ＞ 64byte 时，内部编码 listpack -> hashtable
	rdb.HMSet(ctx, "hashkey", "f3", "hash类型元素的个数  hash-max-ziplist-entries < 512 and 所有的值hash-max-ziplist-value < 65。hash类型元素的个数  hash-max-ziplist-entries < 512 and 所有的值hash-max-ziplist-value < 65。")
	log.Println(rdb.ObjectEncoding(ctx, "hashkey")) // hashtable
}

func TestHash(t *testing.T) {

	// =================================================================
	//  hset key field value
	//  hget key field
	// 	hlen key  -- 计算 field 的个数
	// =================================================================

	result, err := rdb.HGet(ctx, "user:1", "name").Result()
	log.Println("result:", result, "err:", err) // result 默认值，err: redis.Nil

	var user struct {
		Name string `redis:"name"`
		Age  string `redis:"age"`
		City string `redis:"city"`
	}
	rdb.HMSet(ctx, "user:1", "name", "lance", "age", 18, "city", "贵州")

	// 如果 field 很多会阻塞redis
	// hscan 渐进式的遍历
	// rdb.HGetAll(ctx, "user:1").Scan(&user) // field 不是很多的情况
	_ = rdb.HMGet(ctx, "user:1", "name", "age", "city").Scan(&user)
	log.Println(user)
}

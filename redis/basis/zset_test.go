package basis

import (
	"log"
	"strconv"
	"testing"
)

/*
	----------------------------------------------------------------------------------------------------------------------------------
	命令                                                                                 时间复杂度
	----------------------------------------------------------------------------------------------------------------------------------
	zadd key score member [score member ...]                                           O(k*1og(n)), k是添加成员的个数, n是当前有序集合成员个数
	zcard key                                                                          O(1)
	zscore key member                                                                  O(1)
	zrank key member
	zrevrank key member                                               				   O(log(n)), n是当前有序集合成员个数
	zrem key member [member ...]                                                       O(k*log(n)), k是删除成员的个数, n是当前有序集成员个数
	zincrby key increment member                                                       O(log(n)), n是当前有序集合成员个数
	zrange key start end [withscores]
	zrevrange key start end [withscores] 									           O(log(n)+k), k是要获取的成员个数, n是当前有序集合成员个数
	zrangebyscore key min max [withscores]
	zrevrangebyscore key max min [withscores]   									   O(log(n)+k), k是要获取的成员个数, n是当前有序集合成员个数
	zcount                                                                             O(log(n)), n是当前有序集合成员个数
	zremrangebyrank key start end                                                      O(log(n)+k), k是要获取的成员个数, n是当前有序集合成员个数
	zremrangebyscore key min max                                                       O(log(n)+k), k是要获取的成员个数, n是当前有序集合成员个数
	zinterstore destination mumkeys key [key ...]                                      O(n)+O(m*log(m)), n是所有有序集合成员个数和, m是结果集中成员个数
	zunionstore destination numkeys key [key ... ]                                     O(n)+O(m*log(m)), n是所有有序集合成员个数和, m是结果集中成员个数
	----------------------------------------------------------------------------------------------------------------------------------

*/

// 【有序集合中的元素不能重复，但是score可以重复】
//  使用场景
// 		排行榜系统
// 		视频网站需要对用 户上传的视频做排行榜，榜单的维度可能是多个方面的
//    		按照时间、
//			按照播、
// 			放数量、
//			按照获得的赞数

// 三种内部编码
//
//	listpack、skiplist、hashtable
//		- skiplist 当元素个数超过128个
func TestZSetObject(t *testing.T) {
	rdb.ZAdd(ctx, "zsetkey",
		redis.Z{Member: "e1", Score: 50},
		redis.Z{Member: "e2", Score: 100},
		redis.Z{Member: "e3", Score: 150},
	)
	log.Println(rdb.ObjectEncoding(ctx, "zsetkey")) // listpack

	// 元素个数超过128个，使用 skiplist
	var rank []redis.Z
	for i := 0; i <= 129; i++ {
		rank = append(rank, redis.Z{Member: "e" + strconv.Itoa(i), Score: float64(i + 100)})
	}
	rdb.ZAdd(ctx, "zsetkey128", rank...)
	log.Println(rdb.ObjectEncoding(ctx, "zsetkey128")) // skiplist
}

func TestZSet(t *testing.T) {
	members := []redis.Z{
		{Member: "kris", Score: 1},
		{Member: "mike", Score: 91},
		{Member: "frank", Score: 200},
		{Member: "tom", Score: 220},
		{Member: "martin", Score: 251},
	}

	// 添加成员
	rdb.ZAdd(ctx, "user:ranking", members...)
	// 计算某个成员的排名
	log.Println(rdb.ZRank(ctx, "user:ranking", "tom"))    // 分数从低->高（从 0 开始）3
	log.Println(rdb.ZRevRank(ctx, "user:ranking", "tom")) // 分数从高->低（从 0 开始）1

	// 添加成员分数
	rdb.ZIncrBy(ctx, "user:ranking", 3, "tom")

	// 返回指定排名范围内的成员
	// withsocres - 带上成员分数
	// ============================================================
	// zrange key start end [withscores]
	// zrevrange key start end [withscores]
	// ============================================================

	// 返回指定分数范围的成员
	// min和max还支持开区间（小括号）和闭区间（中括号）
	// -inf和 +inf分别代表无限小和无限大
	//
	// limit offset count 起始位置和个数
	//
	// e.g zrangebyscore user:ranking (200 +inf withscores
	// ============================================================
	// zrangebyscore key min max [withscores] [limit offset count]
	// zrevrangebyscore key max min [withscores] [limit offset count]
	// ============================================================
}

func TestUserThumbNum(t *testing.T) {
	// 例子 用户点赞数
	// 用户 kris 上传了一个视频，并获得了 100 点的点赞数
	rdb.ZAdd(ctx, "user:thumb:202402", redis.Z{Member: "kris", Score: 100})
	// 之后再获得一个赞
	rdb.ZIncrBy(ctx, "user:thumb:202402", 1, "kris")
	// 点错了
	rdb.ZIncrBy(ctx, "user:thumb:202402", -1, "kris")
	// 用户作弊取消作品
	rdb.ZRem(ctx, "user:thumb:202402", "kris")
	// 展示👍最多的前十位
	rdb.ZRevRangeWithScores(ctx, "zsetkey128", 0, 9)
}
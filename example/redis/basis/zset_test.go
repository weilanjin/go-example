package basis

import (
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
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

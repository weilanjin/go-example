package basis

import (
	"log"
	"testing"
)

/*
	------------------------------------------------------------------------
	 命令                                     时间复杂度
	------------------------------------------------------------------------
	sadd key element [element...]           O(k), k是元素个数
	srem key element [element...]           O(k), k是元素个数
	scard key                               O(1)
	sismember key element                   O(1)
	srandmember key [count]                 O(count)
	spop key                                O(1)
	smembers key                            O(n), n是元素总数
	sinter key [key . . .] 或者sinterstore   O(m*k), 不是多个集合中元素最少的个数, m是键个数
	suinon key [key . . .] 或者suionstore    O(k), k是多个集合元素个数和
	sdiff key [key . . .] 或者sdiffstore     O(k), k是多个集合元素个数和
	------------------------------------------------------------------------
*/

// 使用场景
// 	1.用户标签
//  2.好友关系

// ·sadd=Tagging（标签）
// ·spop/srandmember=Random item（生成随机数，比如抽奖）
// ·sadd+sinter=Social Graph（社交需求）

// 两种内部编码
//
//		intset、listpack、hashtable
//			1.intset 整数集合
//	 		2.hashtable 哈希表 element > 512 or element != int
func TestSetObject(t *testing.T) {
	rdb.SAdd(ctx, "setIntObj", 1, 2)
	rdb.SAdd(ctx, "setHashObj", "has")

	log.Println(rdb.ObjectEncoding(ctx, "setIntObj"))  // intset
	log.Println(rdb.ObjectEncoding(ctx, "setHashObj")) // listpack
}

func TestSet(t *testing.T) {
	// ============================================================
	// sadd key element ...
	// ============================================================

	// 添加元素
	rdb.SAdd(ctx, "set2", "a", "c", "c")
	log.Println(rdb.SMembers(ctx, "set2"))

	// ============================================================
	//  srem key element ...
	// ============================================================
	// 删除元素
	rdb.SRem(ctx, "set2", "c")

	// 集合间操作
	// ============================================================
	//  srem key element ...
	// ============================================================
	rdb.SAdd(ctx, "user:1:follow", "it", "music", "his", "sports")
	rdb.SAdd(ctx, "user:2:follow", "it", "news", "ent", "sports")

	// 两个集合的交集（两个用户都感兴趣的标签）
	log.Println(rdb.SInter(ctx, "user:1:follow", "user:2:follow")) // it sports
	// 两个集合的并集
	log.Println(rdb.SUnion(ctx, "user:1:follow", "user:2:follow")) // it music his sports news ent
	// 两个集合的并集
	log.Println(rdb.SDiff(ctx, "user:1:follow", "user:2:follow")) // music his
	// 两个集合的交集结果保存
	log.Println(rdb.SInterStore(ctx, "user:1_2:follow", "user:1:follow", "user:2:follow")) // it sports

}

package basis

import (
	"log"
	"testing"
)

// list
// 列表的每个字符串称为元素 element
// 一个列表最多可以存储 2^-1个元素

/*
			列表的四种操作类型
	+----------+----------------------+
	| 操作类型  | 操作                  |
	+----------+----------------------+
	| 添加     | rpush lpush linsert  |
	| 杳       | Irange lindex llen   |
	| 删除     | Ipop rpop Irem Itrim |
	| 修改     | lset                 |
	| 阻塞操作  | blpop brpop          |
	+----------+----------------------+

	-------------------------------------------------------------------------------------
	操作类型   命令                                     时间复杂度
	-------------------------------------------------------------------------------------
	添加     rpush key value [value ...]            O(k)，k是元素个数
			lpush key value [value ...]            O(k)，k是元素个数
			linsert key beforelafter pivot value   O(n)，n是pivot 距离列表头或尾的距离
	查找     range key start end                    O(s+n), s是start 偏移量，n是start到end的范围
			lindex key index                       O(n)，n 是索引的偏移量
			llen key                               O(1)
	删除     lpop key                               O(1)
			rpop key                               O(1)
			lrem count value                       O(n), n是列表长度
			ltrim key start end                    O(n), n是要裁剪的元素总数
	修改     lset key index value                   O(n)，n 是索引的偏移量
	阻塞操作  blpop brpop                            O(1)
	-------------------------------------------------------------------------------------

*/

// 三种内部编码 ziplist、quicklist、linkedlist
//
//			1.ziplist(压缩列表)  -- 内存占用更少
//				element < list-max-ziplist-entries = 512 or value < list-max-ziplist-value = 64byte
//	     	2.quicklist
//		 	3.linkedlist(链表)
//
// quicklist 设计原理 https://matt.sh/redis-quicklist。
func TestObject(t *testing.T) {
	rdb.RPush(ctx, "listObj", "e1", "e2", "e3")
	log.Println(rdb.ObjectEncoding(ctx, "listObj")) // listpack
	var nums []any
	for i := 0; i < 513; i++ {
		nums = append(nums, i)
	}

	rdb.RPush(ctx, "listObj", nums...)
	log.Println(rdb.ObjectEncoding(ctx, "listObj")) // quicklist
}

// 使用场景
//
//		1.消息队列	lpush+brpop命令组合即可实现阻塞队列
//		2.文章列表
//		·lpush+lpop=Stack（栈）
//	 	·lpush+rpop=Queue（队列）
//	 	·lpsh+ltrim=Capped Collection（有限集合）
//	 	·lpush+brpop=Message Queue（消息队列）
func TestList(t *testing.T) {
	rdb.Del(ctx, "listkey1")

	//	var users = []User{{Name: "lance", Age: 18}, {Name: "lanjin", Age: 19}}
	user1 := User{Name: "lance", Age: 18}
	user2 := User{Name: "lanjin", Age: 19}

	// ============================================================
	// lpush key ...value 向左边添加元素
	// rpush key ...value 向右边添加元素
	// ============================================================

	// 从左插入元素
	log.Println(rdb.RPush(ctx, "listkey1", &user1, &user2).Err())

	// ============================================================
	// lrange key start stop
	// lindex key index      获取列表指定索引下的元素
	// ============================================================

	// 获取所有的元素
	var res []User
	err := rdb.LRange(ctx, "listkey1", 0, -1).ScanSlice(&res)
	log.Println(res, err)

	rdb.RPush(ctx, "listkey2", "a", "a", "a", "a", "a", "go", "b", "a", "a", "a")

	// ============================================================
	// linsert key before|after pivot value 向某个元素pivot前或后插入元素
	// ============================================================
	rdb.LInsertBefore(ctx, "listkey2", "a", "c")

	// ============================================================
	// lrem key count value
	// ============================================================
	// lrem命令会从列表中找到等于value的元素进行删除
	// 		count>0，从左到右，删除最多count个元素
	// 		count<0，从右到左，删除最多count绝对值个元素
	//      count=0，删除所有
	rdb.LRem(ctx, "listkey2", 3, "a")

	log.Println(rdb.LRange(ctx, "listkey2", 0, -1))
}

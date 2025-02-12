package basis

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"
)

// https://tableconvert.com/csv-to-ascii
/*
								字符串类型命令时间复杂度
	+--------------------------------+------------------------------------------------+
	| 命令                            | 时间复杂度                                      |
	+--------------------------------+------------------------------------------------+
	| set key value                  | O(1)                                           |
	| get key                        | O(1)                                           |
	| del key [key ...]              | O(k),不是键的个数                                |
	| mset key value [key value ...] | O(k),不是键的个数                                |
	| mget key [key ...]             | O(k),不是键的个数                                |
	| iner key                       | O(1)                                           |
	| decr key                       | O(1)                                           |
	| incrby key increment           | O(1)                                           |
	| decrby key decrement           | O(1)                                           |
	| incrbyfloat key increment      | O(1)                                           |
	| append key value               | O(1)                                           |
	| strlen key                     | O(1)                                           |
	| strange key offset value       | O(1)                                           |
	| getrange key start end         | O(n)，n是字符串长度，由于获取字符串非常快，所以如果字符串不是很长，可以视同为 O(1) |
	+--------------------------------+------------------------------------------------+

	例子
		1. 缓存用户信息
		2. 计数
		3. 共享session
		4. 限速 - 短信
*/

// 三种内部编码
//
//	1.int:    8个字节的长整型。
//	2.embstr: <=39个字节的字符串。
//	3.raw:    >39个字节的字符串。
func TestStringObject(t *testing.T) {
	// =================================================================
	// object encoding key
	// ================================================================
	rdb.Set(ctx, "key:int", 1994, 0)
	rdb.Set(ctx, "embstr39", "hello world", 0) // int
	rdb.Set(ctx, "raw39", "Println calls Output to print to the standard logger........", 0)

	log.Println(rdb.ObjectEncoding(ctx, "key:int"))  // int
	log.Println(rdb.ObjectEncoding(ctx, "embstr39")) // embster
	log.Println(rdb.ObjectEncoding(ctx, "raw39"))    // raw
}

func TestString(t *testing.T) {

	// =================================================================
	// mset mget mxxx 「如果 key 过多会造成 Redis 堵塞。或者网络拥塞」。
	// =================================================================

	us := generate(10)
	// mset k1 v1 k2 v2 一次性设置 kv
	log.Println(rdb.MSet(ctx, flotMapUser(us)...))

	// gset k1 k2 一次性获取 v (按照 key 的顺序返回)
	// 如没有值，对应的位置就会返回 nil
	// TODO 集群模式下是否要遍历？
	mGet := rdb.MGet(ctx, "user1", "user2", "user23")

	log.Println(mGet)
	users, err := mGet.Result()
	if err != nil {
		panic(err)
	}
	for _, v := range users {
		var u User
		user, ok := v.(string)
		if !ok { // v is nil
			continue
		}
		_ = u.UnmarshalBinary([]byte(user))
		log.Println(u)
	}

}

type ResultUser struct {
	User1 string `redis:"user1"`
	User2 string `redis:"user2"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func names(us []User) []string {
	var ns []string
	for _, n := range us {
		ns = append(ns, n.Name)
	}
	return ns
}

func flotMapUser(us []User) []any {
	var res []any
	for _, u := range us {
		res = append(res, u.Name, u)
	}
	return res
}

func generate(n int) (us []User) {
	for i := 1; i < n; i++ {
		us = append(us, User{
			ID:   int64(i * 1e4),
			Name: "user" + strconv.Itoa(i),
			Age:  18,
		})
	}
	return
}
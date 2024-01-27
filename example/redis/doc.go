// Package redis Remote Dictionary Service
//
// http://redis.io
//
// author blog http://antirez.com
// codebase https://github.com/redis/redis
// quicklist https://matt.sh/redis-quicklist

// “业务名:对象名:id：[属性]
//
// *「smembers和lrange、hgetall都属于比较重的命令，如果元素过多存在阻 塞Redis的可能性，这时候可以使用sscan来完成」
package redis

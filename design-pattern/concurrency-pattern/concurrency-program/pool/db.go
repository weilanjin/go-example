package pool

import (
	"time"
)

type DB struct {
}

func (db *DB) SetConnMaxIdleTime(d time.Duration) {}
func (db *DB) SetConnMaxLifetime(d time.Duration) {}
func (db *DB) SetMaxIdleConns(n int)              {} // 最大空闲的连接数 // 默认值是 2, 这个值对应与数据库相关的应用来说太小了.
func (db *DB) SetMaxOpenConns(n int)              {} // 最大的连接数

/*
func (db *DB) conn(ctx context.Context, strategy connReuseStrategy) (*driverConn, error) {
	db.mu.Lock()
	.....
	numFree := len(db.freeConn)
	if strategy == cachedOrNewConn && numFree > 0 { // 使用可重的策略,并且有可重用的连接
		conn := db.freeConn[0] // 使用第一个连接
		copy(db.freeConn, db.freeConn[1:]) // 把所选择的这个连接从freeConn中剔除
		db.freeConn = db.freeConn[:numFree-1]
		conn.inUse = true // 标记此连接使用中
		if conn.expired(lifetime) {// 如果此连接已经过期了
			db.maxLifetimeClosed++
			db.mu.Unlock()
			conn.Close()
			return nil, driver.ErrBadConn
		}
		db.mu.Unlock()
		....
		return conn, nil // 返回这个可重用的连接
	}
	....
	return dc, nil
}
*/

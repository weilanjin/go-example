package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var db *sql.DB

type Config struct {
	DSN         string        `yaml:"dsn"`
	Timeout     time.Duration `yaml:"timeout"`
	MaxOpenConn int           `yaml:"max_open_conn"`
	MaxIdleConn int           `yaml:"max_idle_conn"`
	MaxLifeTime time.Duration `yaml:"max_life_time"`
	MaxIdleTime time.Duration `yaml:"max_idle_time"`
}

func InitDB(conf Config) error {
	if conf.Timeout <= 0 {
		conf.Timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	var err error
	db, err = sql.Open("mysql", conf.DSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	// 影响最大并发数。
	// 过大可能导致数据库负载过高，过小会限制并发性能。
	//一般设置在 100~500，具体根据数据库负载情况调整。
	db.SetMaxOpenConns(conf.MaxOpenConn) // 设置最大连接数
	// 控制保持在池中的空闲连接数。
	// 过大会浪费资源，过小可能导致频繁创建连接，增加延迟。
	// 典型范围是 10~50。
	db.SetMaxIdleConns(conf.MaxIdleConn) // 设置闲置连接数
	// 控制连接存活的最大时间，避免连接长时间占用资源导致 MySQL 关闭连接。
	// 建议设置 30min ~ 1h，防止连接泄露。
	db.SetConnMaxLifetime(conf.MaxLifeTime) // 连接的最大可复用时间
	// 控制空闲连接的最长时间，防止长期空闲的连接占用资源。
	// 典型值 10min，根据业务需求调整。
	db.SetConnMaxIdleTime(conf.MaxIdleTime) // 空闲连接的最大生存时间

	return db.PingContext(ctx)
}

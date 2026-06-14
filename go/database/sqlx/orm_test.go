package sqlx

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

type User struct {
	ID     int
	Name   string
	Age    int
	Status int
}

func (u *User) Mapping() []*Mapping {
	return []*Mapping{
		{Column: "id", Result: &u.ID, Value: u.ID},
		{Column: "name", Result: &u.Name, Value: u.Name},
		{Column: "age", Result: &u.Age, Value: u.Age},
		{Column: "status", Result: &u.Status, Value: u.Status},
	}
}

var orm *DB

func TestMain(m *testing.M) {
	// 1. 设置阶段：在所有测试运行前执行

	timeout, _ := time.ParseDuration("5s")
	orm = NewDB(&DBConfig{
		DSN:         "root:@tcp(127.0.0.1:3306)/telegram_bot?charset=utf8mb4&parseTime=True&loc=Local",
		Timeout:     timeout,
		MaxOpenConn: 10,               //  数据库最大打开连接数
		MaxIdleConn: 5,                // 数据库最大空闲连接数
		MaxLifeTime: 30 * time.Minute, // 数据库连接最大生命周期
		MaxIdleTime: 10 * time.Minute, // 数据库连接最大空闲时间
	})

	slog.Info("db stats", "stats", orm.Stats())

	// 2. 运行测试：m.Run() 执行包内所有的测试函数
	// m.Run() 返回一个整数，作为程序的退出码
	// 会触发该包内所有 TestXxx、BenchmarkXxx 和 ExampleXxx 函数的执行
	exitCode := m.Run()

	// 3. 清理阶段：在所有测试运行后执行
	if err := orm.Close(); err != nil {
		slog.Error("Failed to close database connection", "error", err)
	}

	// 4. 退出程序
	os.Exit(exitCode)
}
func TestInsert(t *testing.T) {
	user := &User{ID: 1, Name: "Alice", Age: 30, Status: 1}
	rowsAffected, err := orm.Insert(context.Background(), "users", []string{"id"}, user)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	t.Logf("Rows affected: %d", rowsAffected)
}

func TestDelete(t *testing.T) {
	where := Wheres{
		{"id = ?", 123},
	}
	if "xx" == "xx" {
		where = append(where, Where{"status = ?", 1})
	}
	rowsAffected, err := orm.Delete(context.Background(), "your_table_name", where)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	t.Logf("Rows affected: %d", rowsAffected)
}

func TestQueryRow(t *testing.T) {
	query := Query[*User]{
		NewRow: func() *User {
			return &User{} // 替换为你的结构体类型
		},
		Where: Wheres{
			{"id = ?", 123},
		},
	}

	row, err := FindOne(context.Background(), orm, "your_table_name", query)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	t.Logf("Queried row: %+v", row)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 数据库连接配置
	db, err := sql.Open("mysql", "root:admin123@tcp(127.0.0.1:3306)/exp")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 生成测试数据
	generateTestData(db, 10000)
}

func generateTestData(db *sql.DB, count int) {
	// 准备插入语句
	stmt, err := db.Prepare(`
		INSERT INTO employees 
		(first_name, last_name, department, salary, hire_date) 
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal("准备插入语句失败:", err)
	}
	defer stmt.Close()

	// 部门列表
	departments := []string{"Engineering", "Marketing", "Sales", "HR", "Finance"}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("开始事务失败:", err)
	}

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成数据
	for i := 0; i < count; i++ {
		// 随机选择部门
		dept := departments[rand.Intn(len(departments))]

		// 基于部门的薪资范围
		var salary int
		switch dept {
		case "Engineering":
			salary = 50000 + rand.Intn(100000)
		case "Marketing":
			salary = 40000 + rand.Intn(80000)
		case "Sales":
			salary = 30000 + rand.Intn(120000)
		case "HR":
			salary = 35000 + rand.Intn(70000)
		default: // Finance
			salary = 45000 + rand.Intn(90000)
		}

		// 随机雇佣日期(最近5年)
		daysAgo := rand.Intn(1825) // 5年大约1825天
		hireDate := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

		// 随机名字
		firstName := fmt.Sprintf("First%d", rand.Intn(1000))
		lastName := fmt.Sprintf("Last%d", rand.Intn(1000))

		// 执行插入
		_, err = stmt.Exec(firstName, lastName, dept, salary, hireDate)
		if err != nil {
			tx.Rollback()
			log.Fatal("插入数据失败:", err)
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatal("提交事务失败:", err)
	}
	fmt.Printf("成功插入 %d 条测试数据\n", count)
}

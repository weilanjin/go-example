package sqlx

import (
	"database/sql"
	"fmt"
	"reflect"
)

func ScanStruct[T any](row *sql.Row) (T, error) {
	var data T

	v := reflect.ValueOf(&data).Elem()

	// 场景 A：如果 T 是结构体，我们需要利用反射将所有字段铺开
	// 注意：这里简单以 Kind() == reflect.Struct 判断，实际生产中可能需要排除 time.Time 或 sql.Nullxxx 等特殊结构体
	if v.Kind() == reflect.Struct {
		numCols := v.NumField()
		columns := make([]interface{}, numCols)

		// 提取结构体内部所有字段的内存地址
		for i := 0; i < numCols; i++ {
			columns[i] = v.Field(i).Addr().Interface()
		}

		// 将查询结果 Scan 进提取出的字段地址中
		err := row.Scan(columns...)
		return data, err
	}

	// 场景 B：如果 T 是基本类型 (如 int, string) 或者官方自带的值类型 (如 sql.NullString)
	// 完全不需要反射字段，直接取 data 的地址进行 Scan
	err := row.Scan(&data)
	return data, err
}

// https://go.dev/wiki/SQLInterface#getting-a-table
func ScanStructSlice[T any](rows *sql.Rows) (out []T) {
	var table []T
	var data T // 声明承接容器

	v := reflect.ValueOf(&data).Elem()

	// ==========================================
	// 场景 A：如果 T 是结构体类型
	// ==========================================
	if v.Kind() == reflect.Struct {
		numCols := v.NumField()
		columns := make([]any, numCols)

		// 循环外：提前获取结构体各字段的内存地址
		for i := 0; i < numCols; i++ {
			columns[i] = v.Field(i).Addr().Interface()
		}

		// 循环内：复用地址，零反射
		for rows.Next() {
			if err := rows.Scan(columns...); err != nil {
				fmt.Println("Case Read Error ", err)
				continue
			}
			table = append(table, data)
		}
		return table
	}

	// ==========================================
	// 场景 B：如果 T 是基础类型 (如 int, string, sql.NullString)
	// ==========================================
	for rows.Next() {
		// 基础类型直接取 data 的地址进行 Scan，完全不需要反射
		if err := rows.Scan(&data); err != nil {
			fmt.Println("Case Read Error ", err)
			continue
		}
		table = append(table, data)
	}

	return table
}

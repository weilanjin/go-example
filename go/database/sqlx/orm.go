package sqlx

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

type Row interface {
	TableName() string
	Mapping() []*Mapping
}

type Mapping struct {
	Column string
	Result any // query result (pointer)
	Value  any // insert, update value
}

type SQL struct {
	*sql.DB
	BatchSize        int  // 批量插入的大小，默认为一次性插入所有行
	SlowSQLThreshold int  // 单位毫秒，超过该值的 SQL 将被记录为慢 SQL
	Debug            bool // 是否启用调试模式，打印 SQL 语句和参数
}

// Insert inserts rows while ignoring specified columns
// 批量插入数据时忽略指定的列
func (db *SQL) Insert(omits []string, rows ...Row) error {
	if len(rows) == 0 {
		return fmt.Errorf("no rows to insert")
	}

	batchSize := db.BatchSize
	if batchSize <= 0 {
		batchSize = len(rows)
	}

	row := rows[0]
	table, mappings := row.TableName(), row.Mapping()

	columns := make([]string, 0, len(mappings))
	for _, m := range mappings {
		if !slices.Contains(omits, m.Column) {
			columns = append(columns, m.Column)
		}
	}

	// Process rows in batches
	// 分批处理行数据
	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}

		var placeholders []string
		var values []any

		for _, r := range rows[i:end] {
			rowPlaceholders := make([]string, len(columns))
			for j := range rowPlaceholders {
				rowPlaceholders[j] = "?"
			}
			placeholders = append(placeholders, "("+strings.Join(rowPlaceholders, ", ")+")")

			for _, m := range r.Mapping() {
				if !slices.Contains(omits, m.Column) {
					values = append(values, m.Value)
				}
			}
		}

		query := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES %s",
			table,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
		)

		if _, err := db.Exec(query, values...); err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes rows based on a column and its value
// 根据列和值删除行
func (db *SQL) Delete(table string, where OrderedMap[string, any]) error {
	var (
		whereClauses []string
		values       []any
	)

	for k, v := range where.All() {
		if v == nil {
			whereClauses = append(whereClauses, k)
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", k))
			values = append(values, v)
		}
	}

	var whereSQL string
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	_, err := db.Exec(fmt.Sprintf("DELETE FROM %s %s", table, whereSQL), values...)
	return err
}

func (db *SQL) Update(table string, id int, m map[string]any) error {
	if len(m) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setClauses := make([]string, 0, len(m))
	values := make([]any, 0, len(m))
	var idValue any = id

	for column, value := range m {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
		values = append(values, value)
	}

	if idValue == nil {
		return fmt.Errorf("id value is required for update")
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = ?",
		table,
		strings.Join(setClauses, ", "),
	)

	values = append(values, idValue) // id value goes last for where clause
	_, err := db.Exec(query, values...)
	return err
}

func List[T Row](where string, args []any, newRow func() T) ([]T, error) {
	if where != "" {
		where = "WHERE " + where
	}

	prototype := newRow()
	mappings := prototype.Mapping()
	columns := make([]string, len(mappings))
	for i, m := range mappings {
		columns[i] = m.Column
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(columns, ", "),
		prototype.TableName(),
		where,
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]T, 0)
	for rows.Next() {
		item := newRow()
		itemMappings := item.Mapping()
		scanArgs := make([]any, len(itemMappings))
		for i, m := range itemMappings {
			scanArgs[i] = m.Result
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

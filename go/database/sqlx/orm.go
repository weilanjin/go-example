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

type Where []struct {
	Clause string
	Value  any
}

type SQL struct {
	*sql.DB
	BatchSize        int  // 批量插入的大小，默认为一次性插入所有行
	SlowSQLThreshold int  // 单位毫秒，超过该值的 SQL 将被记录为慢 SQL
	Debug            bool // 是否启用调试模式，打印 SQL 语句和参数
}

func NewSQL(db *sql.DB) *SQL {
	return &SQL{DB: db}
}

// Insert inserts rows while ignoring specified columns
// 批量插入数据时忽略指定的列
func (db *SQL) Insert(omits []string, rows ...Row) (int64, error) {
	if len(rows) == 0 {
		return 0, fmt.Errorf("no rows to insert")
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

	var totalRowsAffected int64

	// Process rows in batches
	// 分批处理行数据
	for i := 0; i < len(rows); i += batchSize {
		end := min(i+batchSize, len(rows))

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

		res, err := db.Exec(query, values...)
		if err != nil {
			return 0, err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		totalRowsAffected += rowsAffected
	}

	return totalRowsAffected, nil
}

// Delete deletes rows based on a column and its value
// 根据列和值删除行
func (db *SQL) Delete(table string, where Where) (int64, error) {
	whereClauses := make([]string, 0, len(where))
	values := make([]any, 0, len(where))

	for _, w := range where {
		whereClauses = append(whereClauses, w.Clause)
		if w.Value != nil {
			values = append(values, w.Value)
		}
	}

	var whereSQL string
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " ")
	}

	res, err := db.Exec(fmt.Sprintf("DELETE FROM %s %s", table, whereSQL), values...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Update updates rows based on a column and its value
// 根据列和值更新行
func (db *SQL) Update(table string, where Where, m map[string]any) (int64, error) {
	if len(m) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	setClauses := make([]string, 0, len(m))
	values := make([]any, 0, len(m))

	for column, value := range m {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
		values = append(values, value)
	}

	whereClauses := make([]string, 0, len(where))
	whereValues := make([]any, 0, len(where))
	for _, w := range where {
		whereClauses = append(whereClauses, w.Clause)
		if w.Value != nil {
			whereValues = append(whereValues, w.Value)
		}
	}

	var whereSQL string
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " ")
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s %s",
		table,
		strings.Join(setClauses, ", "),
		whereSQL,
	)

	values = append(values, whereValues...) // where values go last for where clause
	res, err := db.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func List[T Row](where Where, args []any, newRow func() T) ([]T, error) {
	var whereSQL string
	if len(where) > 0 {
		whereClauses := make([]string, 0, len(where))
		for _, w := range where {
			whereClauses = append(whereClauses, w.Clause)
			if w.Value != nil {
				args = append(args, w.Value)
			}
		}
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
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
		whereSQL,
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

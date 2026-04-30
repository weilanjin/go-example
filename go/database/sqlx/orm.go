package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"
)

type Row interface {
	Mapping() []*Mapping
}

type Mapping struct {
	Column string
	Result any // query result (pointer)
	Value  any // insert, update value
}

type Where struct {
	Clause string
	Value  any
}

type Wheres []Where

func (w *Wheres) And(clause string, value any) {
	*w = append(*w, Where{Clause: "AND " + clause, Value: value})
}

func (w *Wheres) Or(clause string, value any) {
	*w = append(*w, Where{Clause: "OR " + clause, Value: value})
}

func (w *Wheres) build() (string, []any) {
	whereSQL, args := "", []any{}
	if len(*w) == 0 {
		return whereSQL, args
	}

	var whereClauses []string
	for i, w := range *w {
		clause := w.Clause
		if i == 0 {
			clause = strings.TrimPrefix(strings.TrimPrefix(clause, "AND "), "OR ")
		}
		whereClauses = append(whereClauses, clause)
		if w.Value != nil {
			args = append(args, w.Value)
		}
	}
	if len(whereClauses) == 0 {
		return whereSQL, args
	}

	whereSQL = "WHERE " + strings.Join(whereClauses, " ")
	return whereSQL, args
}

type DB struct {
	*sql.DB
	BatchSize        int // 批量插入的大小，默认为一次性插入所有行
	SlowSQLThreshold int // 单位毫秒，超过该值的 SQL 将被记录为慢 SQL
}

func NewSQL(db *sql.DB) *DB {
	return &DB{DB: db}
}

type printCtx struct{}

var printCtxKey = printCtx{}

func WithPrintCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, printCtxKey, struct{}{})
}

func isPrintCtx(ctx context.Context) bool {
	v := ctx.Value(printCtxKey)
	if v == nil {
		return false
	}
	_, ok := v.(struct{})
	return ok
}

// execContext wraps ExecContext with debug logging and slow SQL detection
// 封装 ExecContext，支持调试日志和慢 SQL 检测
func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if isPrintCtx(ctx) {
		log.Printf("[SQL] exec: %s args: %v", query, args)
	}
	start := time.Now()
	res, err := db.DB.ExecContext(ctx, query, args...)
	if elapsed := time.Since(start); db.SlowSQLThreshold > 0 && elapsed.Milliseconds() >= int64(db.SlowSQLThreshold) {
		log.Printf("[SlowSQL] %dms exec: %s args: %v", elapsed.Milliseconds(), query, args)
	}
	return res, err
}

// queryContext wraps QueryContext with debug logging and slow SQL detection
// 封装 QueryContext，支持调试日志和慢 SQL 检测
func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if isPrintCtx(ctx) {
		log.Printf("[SQL] query: %s args: %v", query, args)
	}
	start := time.Now()
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if elapsed := time.Since(start); db.SlowSQLThreshold > 0 && elapsed.Milliseconds() >= int64(db.SlowSQLThreshold) {
		log.Printf("[SlowSQL] %dms query: %s args: %v", elapsed.Milliseconds(), query, args)
	}
	return rows, err
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if isPrintCtx(ctx) {
		log.Printf("[SQL] query row: %s args: %v", query, args)
	}
	start := time.Now()
	row := db.DB.QueryRowContext(ctx, query, args...)
	if elapsed := time.Since(start); db.SlowSQLThreshold > 0 && elapsed.Milliseconds() >= int64(db.SlowSQLThreshold) {
		log.Printf("[SlowSQL] %dms query row: %s args: %v", elapsed.Milliseconds(), query, args)
	}
	return row
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.queryContext(ctx, query, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.execContext(ctx, query, args...)
}

// Insert inserts rows while ignoring specified columns
// 批量插入数据时忽略指定的列
func (db *DB) Insert(ctx context.Context, table string, omits []string, rows ...Row) (int64, error) {
	if len(rows) == 0 {
		return 0, fmt.Errorf("no rows to insert")
	}

	batchSize := db.BatchSize
	if batchSize <= 0 {
		batchSize = len(rows)
	}

	row := rows[0]
	mappings := row.Mapping()

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

		res, err := db.execContext(ctx, query, values...)
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
func (db *DB) Delete(ctx context.Context, table string, where Wheres) (int64, error) {
	whereSQL, args := where.build()
	res, err := db.execContext(ctx, fmt.Sprintf("DELETE FROM %s %s", table, whereSQL), args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Update updates rows based on a column and its value
// 根据列和值更新行
func (db *DB) Update(ctx context.Context, table string, where Wheres, update map[string]any) (int64, error) {
	if len(update) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	setClauses := make([]string, 0, len(update))
	values := make([]any, 0, len(update))
	for column, value := range update {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
		values = append(values, value)
	}

	whereSQL, args := where.build()
	query := fmt.Sprintf(
		"UPDATE %s SET %s %s",
		table,
		strings.Join(setClauses, ", "),
		whereSQL,
	)

	values = append(values, args...) // where values go last for where clause
	res, err := db.execContext(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (db *DB) FindOne(ctx context.Context, table string, query Query) (Row, error) {
	query.Limit = 1

	rows, err := db.Find(ctx, table, query)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found")
	}
	return rows[0], nil
}

type Query struct {
	NewRow  func() Row // 如何创建对象
	Where   Wheres     // 过滤条件
	OrderBy string     // 排序条件，例如 "created_at DESC"
	Limit   int        // 返回的最大行数，0表示不限制
	Offset  int        // 分页偏移量，0表示不偏移
}

// Find queries rows with support for filtering, sorting, and pagination
// 支持过滤、排序和分页的查询方法
func (db *DB) Find(ctx context.Context, table string, query Query) ([]Row, error) {
	if query.NewRow == nil {
		return nil, fmt.Errorf("NewRow function is required")
	}

	// Build WHERE clause / 构建WHERE子句
	var whereSQL, args = query.Where.build()

	prototype := query.NewRow()
	mappings := prototype.Mapping()
	columns := make([]string, len(mappings))
	for i, m := range mappings {
		columns[i] = m.Column
	}

	// Build query with ORDER BY and LIMIT / 构建带排序和分页的查询语句
	sqlText := fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(columns, ", "),
		table,
		whereSQL,
	)

	if query.OrderBy != "" {
		sqlText += " ORDER BY " + query.OrderBy
	}

	if query.Limit > 0 {
		sqlText += fmt.Sprintf(" LIMIT %d", query.Limit)
		if query.Offset > 0 {
			sqlText += fmt.Sprintf(" OFFSET %d", query.Offset)
		}
	}

	rows, err := db.queryContext(ctx, sqlText, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]Row, 0)
	for rows.Next() {
		item := query.NewRow()
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

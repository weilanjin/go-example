package std

import (
	"database/sql"
	"fmt"
	"strings"
)

// ------------ model -------------------
type User struct {
	ID           int
	Name         string
	Email        string
	PhoneNo      string
	RegisterTime int64
}

func (*User) TableName() string {
	return "users"
}

func (m *User) Mapping() map[string]func(*User) any {
	return map[string]func(*User) any{
		"mame":         func(u *User) any { return u.Name },
		"email":        func(u *User) any { return u.Email },
		"phoneNo":      func(u *User) any { return u.PhoneNo },
		"registerTime": func(u *User) any { return u.RegisterTime },
	}
}

func (m *User) MappingSelect() map[string]func(*User) any {
	return map[string]func(*User) any{
		"id":           func(u *User) any { return &u.ID },
		"mame":         func(u *User) any { return &u.Name },
		"email":        func(u *User) any { return &u.Email },
		"phoneNo":      func(u *User) any { return &u.PhoneNo },
		"registerTime": func(u *User) any { return &u.RegisterTime },
	}
}

// ------------ repository -------------------
type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// 查询单个
func (repo *UserRepo) FindOne(id int) (*User, error) {
	res := new(User)
	mapping := res.MappingSelect()
	var (
		colunms = make([]string, 0, len(mapping))
		values  = make([]any, 0, len(mapping))
	)
	for k, fn := range mapping {
		colunms = append(colunms, k)
		values = append(values, fn(res))
	}

	sqlText := "SELECT " + strings.Join(colunms, ",") + " FROM " + res.TableName() + " WHERE id = ?"
	row := repo.db.QueryRow(sqlText, id)
	err := row.Scan(values...)
	return res, err
}

// 查询多个
func (repo *UserRepo) Find(registerTime int64) ([]*User, error) {
	model := new(User)
	mapping := model.MappingSelect()

	columns := make([]string, 0, len(mapping))
	for col := range mapping {
		columns = append(columns, col)
	}

	cond := []string{
		"1 = 1", // Always true condition to simplify appending
	}
	if registerTime > 0 {
		cond = append(cond, "AND registerTime > ?")
	}

	sqlText := "SELECT " + strings.Join(columns, ",") + " FROM " + model.TableName() + " WHERE " + strings.Join(cond, " ") + " ORDER BY id DESC"
	rows, err := repo.db.Query(sqlText, registerTime)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	var list []*User
	for rows.Next() {
		record := new(User)
		values := make([]any, 0, len(columns))
		for _, col := range columns {
			values = append(values, mapping[col](record))
		}
		if err = rows.Scan(values...); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		list = append(list, record)
	}

	return list, nil
}

// 分页查询
func (repo *UserRepo) FindPage(page, pageSize int, name string) ([]*User, int64, error) {
	if pageSize <= 0 {
		pageSize = 10 // 默认每页10条
	}
	var total int64

	sqlCount := "SELECT COUNT(*) FROM " + new(User).TableName()
	if err := repo.db.QueryRow(sqlCount).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("db.QueryRow: %w", err)
	}
	if total == 0 {
		return []*User{}, 0, nil
	}

	model := new(User)
	mapping := model.MappingSelect()

	columns := make([]string, 0, len(mapping))
	for col := range mapping {
		columns = append(columns, col)
	}

	cond := []string{
		"1 = 1", // Always true condition to simplify appending
	}
	if name != "" {
		cond = append(cond, "AND name LIKE ?")
	}

	offset := (page - 1) * pageSize

	sqlText := "SELECT " + strings.Join(columns, ",") + " FROM " + model.TableName() + " WHERE " + strings.Join(cond, " ") + " ORDER BY id DESC LIMIT ? OFFSET ?"
	rows, err := repo.db.Query(sqlText, pageSize, offset, name)
	if err != nil {
		return nil, 0, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	var list []*User
	for rows.Next() {
		record := new(User)
		values := make([]any, 0, len(columns))
		for _, col := range columns {
			values = append(values, mapping[col](record))
		}
		if err = rows.Scan(values...); err != nil {
			return nil, 0, fmt.Errorf("rows.Scan: %w", err)
		}
		list = append(list, record)
	}
	return list, total, nil
}

// 创建
func (repo *UserRepo) Create(user *User) error {
	mapping := user.Mapping()

	var (
		colunms = make([]string, 0, len(mapping))
		pos     = make([]string, 0, len(mapping))
		values  = make([]any, 0, len(mapping))
	)
	for k, fn := range mapping {
		colunms = append(colunms, k)
		pos = append(pos, "?")
		values = append(values, fn(user))
	}

	sqlText := "INSERT INTO " + user.TableName() + " (" + strings.Join(colunms, ",") + ") VALUES (" + strings.Join(pos, ",") + ")"
	_, err := repo.db.Exec(sqlText, values...)
	return err
}

// 批量创建
func (repo *UserRepo) CreateBatchMax(users []*User, maxBatchSize int) error {
	if len(users) == 0 {
		return nil
	}
	if maxBatchSize <= 0 {
		maxBatchSize = 1000 // 默认最大值
	}
	// 字段映射
	mapping := users[0].Mapping()

	columns := make([]string, 0, len(mapping))
	placeholders := make([]string, 0, len(mapping))
	for col := range mapping {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
	}
	posStr := "(" + strings.Join(placeholders, ",") + ")"
	tableName := users[0].TableName()

	// 分批处理
	for i := 0; i < len(users); i += maxBatchSize {
		end := min(i+maxBatchSize, len(users))
		batch := users[i:end]

		// 构建批次SQL
		var (
			valuePlaceholders []string
			allValues         []any
		)
		for _, user := range batch {
			valuePlaceholders = append(valuePlaceholders, posStr)
			for _, col := range columns {
				allValues = append(allValues, mapping[col](user))
			}
		}

		sqlText := "INSERT INTO " + tableName + " (" + strings.Join(columns, ",") + ") VALUES " + strings.Join(valuePlaceholders, ",")
		if _, err := repo.db.Exec(sqlText, allValues...); err != nil {
			return err
		}
	}
	return nil
}

// 更新
func (repo *UserRepo) Update(user *User) error {
	mapping := user.Mapping()

	var (
		setClauses = make([]string, 0, len(mapping))
		values     = make([]any, 0, len(mapping)+1)
	)
	for k, fn := range mapping {
		setClauses = append(setClauses, k+" = ?")
		values = append(values, fn(user))
	}
	values = append(values, user.ID)

	sqlText := "UPDATE " + user.TableName() + " SET " + strings.Join(setClauses, ",") + " WHERE id = ?"
	_, err := repo.db.Exec(sqlText, values...)
	return err
}

// 删除
func (repo *UserRepo) Delete(id int) error {
	sqlText := "DELETE FROM " + new(User).TableName() + " WHERE id = ?"
	_, err := repo.db.Exec(sqlText, id)
	return err
}

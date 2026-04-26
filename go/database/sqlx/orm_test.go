package sqlx

import (
	"context"
	"testing"
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

func TestInsert(t *testing.T) {
	db := NewSQL(nil)
	user := &User{ID: 1, Name: "Alice", Age: 30, Status: 1}
	rowsAffected, err := db.Insert(context.Background(), "users", []string{"id"}, user)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	t.Logf("Rows affected: %d", rowsAffected)
}

func TestDelete(t *testing.T) {
	db := NewSQL(nil)
	where := Wheres{
		{"id = ?", 123},
	}
	if "xx" == "xx" {
		where = append(where, Where{"status = ?", 1})
	}
	rowsAffected, err := db.Delete(context.Background(), "your_table_name", where)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	t.Logf("Rows affected: %d", rowsAffected)
}

func TestQueryRow(t *testing.T) {
	db := NewSQL(nil)
	query := Query{
		NewRow: func() Row {
			return &User{} // 替换为你的结构体类型
		},
		Where: Wheres{
			{"id = ?", 123},
		},
	}

	row, err := db.FindOne(context.Background(), "your_table_name", query)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	t.Logf("Queried row: %+v", row)
}

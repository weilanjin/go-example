package sqlx

import "testing"

func TestDelete(t *testing.T) {
	db := NewSQL(nil)
	where := Where{
		{"id = ?", 123},
	}
	// if "xx" == "xx" {
	// 	where = append(where, Where{"status = ?", 1})
	// }
	rowsAffected, err := db.Delete("your_table_name", where)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	t.Logf("Rows affected: %d", rowsAffected)
}

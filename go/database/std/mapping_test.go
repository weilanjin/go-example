package std

//
//import (
//	"fmt"
//	"testing"
//)
//
//type ALevel struct {
//	LevelName string `json:"level_name"`
//	Rate      int    `db:"rate"`
//}
//
//type BInfo struct {
//	Flag  bool
//	Score int
//}
//
//type Row struct {
//	ID     int64 `json:"id"`
//	Name   string
//	ALevel       // 匿名
//	Info   BInfo // 具名
//}
//
//
//
//type ALevel struct {
//	LevelName string `json:"level_name"`
//	Rate      int    `db:"rate"`
//}
//
//type Info struct {
//	Flag  bool
//	Score int
//}
//
//type Row struct {
//	ID     int64  `db:"id"`
//	Name   string `db:"name"`
//	ALevel        // 匿名
//	Detail Info   // 具名
//}
//
//row := &Row{
//ID:     1,
//Name:   "Tom",
//ALevel: ALevel{LevelName: "gold", Rate: 10},
//Detail: Info{Flag: true, Score: 99},
//}
//
//m := StructToMappings(row)
//
//
//func TestMapping(t *testing.T) {
//	row := &Row{}
//	cols, ptrs := StructToSQLColumnsAndPtrs(row)
//
//	fmt.Println("cols =", cols)
//	fmt.Println("ptrs =", ptrs)
//}
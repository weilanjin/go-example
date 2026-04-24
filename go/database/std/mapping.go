package std

//
//import (
//	"reflect"
//	"strings"
//	"unicode"
//)
//
//func camelToSnake(s string) string {
//	var sb strings.Builder
//	for i, r := range s {
//		if unicode.IsUpper(r) {
//			if i > 0 {
//				sb.WriteByte('_')
//			}
//			sb.WriteRune(unicode.ToLower(r))
//		} else {
//			sb.WriteRune(r)
//		}
//	}
//	return sb.String()
//}
//
//func StructToMappings(v any) []Mapping {
//	rv := reflect.ValueOf(v)
//	if rv.Kind() != reflect.Pointer {
//		panic("StructToMappings requires a pointer to struct")
//	}
//	return parseMappings(rv, "")
//}
//
//func parseMappings(rv reflect.Value, prefix string) []Mapping {
//	rt := rv.Type()
//
//	// pointer → elem
//	if rt.Kind() == reflect.Pointer {
//		rv = rv.Elem()
//		rt = rt.Elem()
//	}
//
//	mappings := make([]Mapping, 0)
//
//	for i := 0; i < rt.NumField(); i++ {
//		f := rt.Field(i)
//		fv := rv.Field(i)
//
//		// ---- 1. 匿名组合 struct ----
//		if f.Anonymous {
//			mappings = append(mappings, parseMappings(fv.Addr(), prefix)...)
//			continue
//		}
//
//		// ---- 2. 具名组合 struct：加前缀展开 ----
//		if f.Type.Kind() == reflect.Struct && f.Type.PkgPath() != "" {
//			newPrefix := prefix + camelToSnake(f.Name) + "_"
//			mappings = append(mappings, parseMappings(fv.Addr(), newPrefix)...)
//			continue
//		}
//
//		// ---- 3. 普通字段 ----
//		var col string
//		if db := f.Tag.Get("db"); db != "" && db != "-" {
//			col = db
//		} else if json := f.Tag.Get("json"); json != "" && json != "-" {
//			col = json
//		} else {
//			col = camelToSnake(f.Name)
//		}
//
//		m := Mapping{
//			Column: prefix + col,
//			Args:   fv.Interface(),
//			Dest:   fv.Addr().Interface(),
//		}
//
//		mappings = append(mappings, m)
//	}
//
//	return mappings
//}
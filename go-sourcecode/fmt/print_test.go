package fmt_test

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

/*
--------------------------------------------------------------------------------
| data type   | syntax                                                      |
| type        | %T(integer) = int %T(&integer) = *int                       |
| int         | %d = 23, %b = 10111, %o = 27, %x = 17                       |
| bool        | %t(true) = true                                             |
| float       | %.4f(3.141592xxx) = 3.1416                                  |
| string      | %s = foo "bar", %q = "foo \"bar\"", %#q = `foo "bar"`       |
| bytes       | %s([]byte("a⌘")) = a⌘                                       |
| map         | %v = map[peanut:true]                                       |
| slice       | %q = ["Kitano" "Kobayashi"]                                 |
| struct      | %+v = {Name:Kim Age:22}                                     |
--------------------------------------------------------------------------------
*/
func TestPrintf(t *testing.T) {

	integer := 23
	str := `foo "bar"`
	m := map[string]bool{"peanut": true}
	person := struct {
		Name string
		Age  int
	}{"Kim", 22}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("| %-12s| %-60s|\n", "data type", "syntax")
	fmt.Printf("| %-12s| %-60s|\n", "type", fmt.Sprintf("%%T(integer) = %T %%T(&integer) = %T", integer, &integer))
	// hex      %h 十六进制
	// decimal  %d 十进制
	// octal    %o 八进制
	// binary   %b 二进制
	fmt.Printf("| %-12s| %-60s|\n", "int", fmt.Sprintf("%%d = %d, %%b = %b, %%o = %o, %%x = %x", integer, integer, integer, integer))
	fmt.Printf("| %-12s| %-60s|\n", "bool", fmt.Sprintf("%%t(true) = %t", true))
	// %.4f 保留4位小数
	fmt.Printf("| %-12s| %-60s|\n", "float", fmt.Sprintf("%%.4f(3.141592xxx) = %.4f", math.Pi))
	fmt.Printf("| %-12s| %-60s|\n", "string", fmt.Sprintf("%%s = %s, %%q = %q, %%#q = %#q", str, str, str))
	fmt.Printf("| %-12s| %-60s|\n", "bytes", fmt.Sprintf(`%%s([]byte("a⌘")) = %s`, []byte("a⌘")))
	fmt.Printf("| %-12s| %-60s|\n", "map", fmt.Sprintf("%%v = %v", m))
	// slice %q 只能对string 有用
	fmt.Printf("| %-12s| %-60s|\n", "slice", fmt.Sprintf("%%q = %q", []string{"Kitano", "Kobayashi"}))
	fmt.Printf("| %-12s| %-60s|\n", "struct", fmt.Sprintf("%%+v = %+v", person))
	fmt.Println(strings.Repeat("-", 80))
}
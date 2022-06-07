package leetcode

import (
	"log"
	"testing"
	"unicode"
)

func Test6(t *testing.T) {
	x, y, z := move("R2(LF)", 0, 0, top)
	log.Println(x, y, z)
}

const (
	left = iota
	right
	top
	bottom
)

// R 向左
// L 向右
// F 向前
// B 后退
// R2(LF) = RLFLF
func move(cmd string, x0, y0, z0 int) (x, y, z int) {
	x, y, z = x0, y0, z0
	repeat := 0
	repeatCmd := ""
	for _, v := range cmd {
		switch {
		case unicode.IsNumber(v):
			repeat = repeat*10 + (int(v) - '0')
		case v == ')':
			for i := 0; i < repeat; i++ {
				move(repeatCmd, x, y, z)
			}
			repeat = 0
			repeatCmd = ""
		case repeat > 0 && v != '(' && v != ')':
			repeatCmd = repeatCmd + string(v)
		case v == 'L':
			z = (z + 1) % 4
		case v == 'R':
			z = (z - 1 + 4) % 4
		case v == 'F':
			switch {
			case z == left || z == right:
				x = x - z + 1
			case z == top || z == bottom:
				y = y - z + 2
			}
		case v == 'B':
			switch {
			case z == left || z == right:
				x = x + z - 1
			case z == top || z == bottom:
				y = y + z - 2

			}
		}
	}
	return
}

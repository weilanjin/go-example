package go126

import "fmt"

func fmtErrorf() {
	err := fmt.Errorf("go version") // errors.New("go version")
	_ = err
}

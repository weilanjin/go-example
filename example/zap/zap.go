package main

import "fmt"

func main() {
	//logger := zap.Must(zap.NewDevelopment()).Named("go")
	fmt.Printf("\x1b[%dm%s\x1b[0m", 31, "weilanjin")
}
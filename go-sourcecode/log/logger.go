package log

import (
	"io"
	"log"
	"os"
)

const prefix = "[Server] "

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	// 可选
	file := openFile("server.log")
	mw := io.MultiWriter(os.Stderr, file)
	log.SetOutput(mw)
}

func Debug(msg string) {
	log.SetPrefix(prefix + "DEBUG ")
	log.Output(2, msg)
}

func Info(msg string) {
	log.SetPrefix(prefix + "INFO  ")
	log.Output(2, msg)
}

func Error(msg string) {
	log.SetPrefix(prefix + "ERROR ")
	log.Output(2, msg)
}

func openFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		Error("open file failed")
		panic(err)
	}
	return f
}
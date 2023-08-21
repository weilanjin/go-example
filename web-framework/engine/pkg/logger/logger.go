package logger

import (
	"fmt"
	"io"
	"log"
	"lovec.wlj/web-framework/engine/pkg/file"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Settings struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	TimeFormat string `yaml:"time-format"`
}

type logLevel int

const (
	debugLevel logLevel = iota
	infoLevel
	errorLevel
	fatalLevel
)

func (l logLevel) Name() string {
	levels := []string{"DEBUG", "INFO", "ERROR", "FATAL"}
	if l < 0 || int(l) >= len(levels) {
		return "INFO"
	}
	return levels[l]
}

var (
	mu                 sync.Mutex
	logger             *log.Logger
	defaultPrefix      = ""
	defaultCallerDepth = 2
)

const flags = log.LstdFlags

func init() {
	logger = log.New(os.Stdout, defaultPrefix, flags)
}

func Setup(settings *Settings) {
	dir := settings.Path
	filename := fmt.Sprintf("%s-%s.%s", settings.Name, time.Now().Format("20060102150405"), settings.Ext)
	logFile, err := file.MustOpen(filename, dir)
	if err != nil {
		log.Fatalf("logging.Setup: %v", err)
	}
	nw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(nw, defaultPrefix, flags)
}

func setPrefix(level logLevel) {
	var logPrefix string
	_, f, line, ok := runtime.Caller(defaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s] %s:%d ", level.Name(), filepath.Base(f), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", level.Name())
	}
	logger.SetPrefix(logPrefix)
}

func Debug(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(debugLevel)
	logger.Println(v...)
}

func Info(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(infoLevel)
	logger.Println(v...)
}

func Error(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(errorLevel)
	logger.Println(v...)
}

func Fatal(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(fatalLevel)
	logger.Println(v...)
}
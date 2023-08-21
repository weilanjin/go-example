package config

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var ServerConfig *Server

type Server struct {
	Bind       string `conf:"bind"`
	Port       int    `conf:"port"`
	MaxClients int    `conf:"maxClients"`
}

func SetupConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	ServerConfig = parse(file)
}

func DefServer() {
	ServerConfig = &Server{
		Bind:       "0.0.0.0",
		Port:       8080,
		MaxClients: 3000,
	}
}

func parse(src io.Reader) *Server {
	rawMap := make(map[string]string)

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' { // 跳过注释
			continue
		}
		pivot := strings.IndexAny(line, " ")
		if pivot > 0 && pivot < len(line)-1 {
			key := line[:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[strings.ToLower(key)] = value
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	server := new(Server)
	t := reflect.TypeOf(server)
	v := reflect.ValueOf(server)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fValue := v.Elem().Field(i)
		key, ok := field.Tag.Lookup("conf")
		if !ok {
			key = field.Name
		}
		value, ok := rawMap[strings.ToLower(key)]
		if ok {
			switch field.Type.Kind() {
			case reflect.String:
				fValue.SetString(value)
			case reflect.Int:
				if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
					fValue.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := "true" == value
				fValue.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					split := strings.Split(value, ",")
					fValue.Set(reflect.ValueOf(split))
				}
			}

		}
	}
	return server
}
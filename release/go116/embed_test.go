package go116

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

// 只读的

//go:embed doc.txt
var doc string

//go:embed person.json
var personFile []byte

func TestOneFile(t *testing.T) {
	fmt.Println(doc) // this is a doc

	// Test json file
	type Person struct {
		Name string
		Age  int
	}

	var person Person
	json.Unmarshal(personFile, &person)
	fmt.Println(person)
}

//go:embed static
var fsEmbed embed.FS

func TestFs(t *testing.T) {
	// fs := http.FileServer(http.Dir("static")) // 创建一个文件服务器, 文件更新会实时更新
	fs := http.FileServerFS(fsEmbed) // 文件更新不会实时更新
	http.Handle("/", fs)

	log.Println("Listening on : 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
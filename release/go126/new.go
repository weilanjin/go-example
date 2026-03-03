package go126

import (
	"encoding/json"
	"time"
)

type Person struct {
	Name string `json:"name"`
	Age  *int   `json:"age"` // age if known, otherwise null
}

func personJson(name string, born time.Time) ([]byte, error) {
	return json.Marshal(Person{
		Name: name,
		Age:  new(yearsSince(born)),
	})
}

func yearsSince(t time.Time) int {
	return int(time.Since(t).Hours() / (24 * 365.25)) // approximate years
}

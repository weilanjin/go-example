package scan

import (
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/redis/go-redis/v9"
	"lovec.wlj/example/redis/initialize"
	"testing"
)

var rdb redis.UniversalClient

func init() {
	rdb = initialize.Redis()
}

type Model struct {
	Str1    string   `redis:"str1"`
	Str2    string   `redis:"str2"`
	Int     int      `redis:"int"`
	Bool    bool     `redis:"bool"`
	Ignored struct{} `redis:"-"`
}

func TestScan1(t *testing.T) {
	ctx := context.Background()
	// Set some fields.
	model := Model{
		Str1: "hello",
		Str2: "world",
		Int:  123,
		Bool: true, // 	rdb.HSet(ctx, "user1", "bool", 1)
	}

	rdb.HSet(ctx, "user1", &model)

	var model1, model2 Model
	// Scan all fields into the model.
	if err := rdb.HGetAll(ctx, "user1").Scan(&model1); err != nil {
		panic(err)
	}

	// Or scan a subset of the fields.
	if err := rdb.HMGet(ctx, "user1", "str1", "int").Scan(&model2); err != nil {
		panic(err)
	}

	spew.Dump(model1)
	spew.Dump(model2)
}

type UserInfo struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Age      int      `json:"age"`
	IsSingle bool     `json:"isSingle"`
	Ignored  struct{} `json:"ignored"`
}

func (u *UserInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func TestScan2(t *testing.T) {
	ctx := context.Background()
	userInfo := UserInfo{
		Username: "lance",
		Password: "xxx1234",
		Age:      29,
		IsSingle: true,
	}

	if err := rdb.Set(ctx, "userInfo", &userInfo, 0).Err(); err != nil {
		panic(err)
	}
	var user1 UserInfo
	if err := rdb.Get(ctx, "userInfo").Scan(&user1); err != nil {
		panic(err)
	}
	spew.Dump(user1)
}

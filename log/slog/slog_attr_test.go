package gslog

import (
	"errors"
	"log/slog"
	"testing"
)

// slog.Attr 可以用来统一规定 kv key 的定义

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func UserID(id int) slog.Attr {
	return slog.Int("user_id", id)
}

func TestAtt(t *testing.T) {
	slog.Info("test slog attr", UserID(12), Error(errors.New("this is a error")))
}

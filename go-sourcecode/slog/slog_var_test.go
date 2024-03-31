package gslog

import (
	"log"
	"log/slog"
	"strings"
	"testing"
)

type User struct {
	Username string
	Password string
	Email    string
}

// 特殊字段处理
func (u User) LogValue() slog.Value {
	if u.Password != "" {
		u.Password = "*******"
	}
	if u.Email != "" {
		idx := strings.Index(u.Email, "@")
		if idx != -1 {
			ename, domain := u.Email[:idx], u.Email[idx:]
			if len(ename) == 1 { //  a@qq.com -> *@qq.com
				ename = "*"
			} else if len(ename) == 2 { // aa@qq.com -> a*@qq.com
				ename = ename[:1] + "*"
			} else if len(ename) > 2 { // aaac@qq.com -> a**c@qq.com
				ename = ename[:1] + strings.Repeat("*", len(ename)-2) + ename[len(ename)-1:]
			}
			u.Email = ename + domain
		}
	}
	return slog.GroupValue(
		slog.String("username", u.Username),
		slog.String("email", u.Email),
		slog.String("password", u.Password),
	)
}

// 特殊字段处理
func TestVar(t *testing.T) {
	user := User{
		Username: "lance",
		Password: "xxxx",
		Email:    "12345547@qq.com",
	}
	slog.Info("user info", slog.Any("user", user))
	log.Println(user)
}

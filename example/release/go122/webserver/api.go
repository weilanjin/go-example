package main

import (
	"encoding/json"
	"net/http"
)

type UserServer struct{}

func NewUserServer() *UserServer {
	return &UserServer{}
}

func (s *UserServer) UserInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	data := map[string]any{
		"code": 0,
		"data": map[string]string{"userID": userID},
		"msg":  "ok",
	}
	bytes, _ := json.Marshal(data)
	w.Write(bytes)
}

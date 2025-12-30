package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/go-telegram/bot"
)

const (
	token = "8441906451:AAGMpRGiyFi3HRe-06cfchlqKf8pmlS-OdA" // @lancewei_bot
)

var (
	ctx = context.Background()
)

// httpProxyOption 返回通过代理访问 Telegram API 的选项
func httpProxyOption() bot.Option {
	proxyURL, _ := url.Parse("http://127.0.0.1:7890")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return bot.WithHTTPClient(5*time.Second, &http.Client{Transport: transport})
}

// TestGetMe 测试获取机器人信息
func TestGetMe(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}
	me, err := botClient.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}
	data, _ := json.MarshalIndent(me, "", "  ")
	fmt.Printf("Bot Me:\n %s\n", data)
}
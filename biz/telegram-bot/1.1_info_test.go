package main

import (
	"testing"

	"github.com/go-telegram/bot"
)

func TestSetMyName(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}
	_, err = botClient.SetMyName(ctx, &bot.SetMyNameParams{
		Name: "本地调试🤖",
	})
	if err != nil {
		t.Fatalf("SetMyName: %v", err)
	}
}

func TestSetMyShortDescription(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}
	_, err = botClient.SetMyShortDescription(ctx, &bot.SetMyShortDescriptionParams{
		ShortDescription: "这是一个用于本地调试的机器人。",
	})
	if err != nil {
		t.Fatalf("SetMyShortDescription: %v", err)
	}
}

func TestSetMyDescription(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}
	_, err = botClient.SetMyDescription(ctx, &bot.SetMyDescriptionParams{
		Description: "欢迎使用本地调试机器人！这个机器人可以帮助你测试和调试Telegram Bot的各种功能。",
	})
	if err != nil {
		t.Fatalf("SetMyDescription: %v", err)
	}
}

func TestGetMyName(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption(), bot.WithDebug(), bot.WithSkipGetMe())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}
	name, err := botClient.GetMyName(ctx, &bot.GetMyNameParams{})
	if err != nil {
		t.Fatalf("GetMyName: %v", err)
	}
	t.Logf("Bot Name: %s", name.Name)
}
package main

import (
	"testing"
)

func TestSetMyName(t *testing.T) {
	if err := b.SetMyName("本地调试🤖", ""); err != nil {
		t.Fatalf("SetMyName: %v", err)
	}
}

func TestSetMyShortDescription(t *testing.T) {
	if err := b.SetMyShortDescription("这是一个用于本地调试的机器人。^^", ""); err != nil {
		t.Fatalf("SetMyShortDescription: %v", err)
	}
}

func TestSetMyDescription(t *testing.T) {
	if err := b.SetMyDescription("oo欢迎使用本地调试机器人！这个机器人可以帮助你测试和调试Telegram Bot的各种功能。", ""); err != nil {
		t.Fatalf("SetMyDescription: %v", err)
	}
}

func TestGetMyName(t *testing.T) {
	info1, _ := b.MyName("")
	info2, _ := b.MyShortDescription("")
	info3, _ := b.MyDescription("")
	t.Logf("Bot Name: %+v", info1.Name)
	t.Logf("Bot Short Description: %+v", info2.ShortDescription)
	t.Logf("Bot Description: %+v", info3.Description)
}

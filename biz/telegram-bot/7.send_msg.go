package main

import (
	"testing"

	"gopkg.in/telebot.v4"
)

func TestSendRichMsg(t *testing.T) {
	chatId := int64(123456789) // 替换为实际的聊天ID
	what := "这是一个测试消息，包含**加粗文本**、_斜体文本_和[链接](https://example.com)。"
	opts := &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdownV2,
	}

	b.Send(&telebot.Chat{ID: chatId}, what, opts)
}

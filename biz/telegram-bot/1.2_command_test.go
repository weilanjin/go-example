package main

import (
	"context"
	"log"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// https://core.telegram.org/bots/api#botcommand
func TestSetMyCommands(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}

	// 设置命令列表
	_, err = botClient.SetMyCommands(context.Background(), &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Start the bot"},
			{Command: "help", Description: "Show help info"},
			{Command: "profile", Description: "Show your profile"},
			{Command: "ping", Description: "Check bot status"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("bot commands set successful")
}

func TestDeleteMyCommands(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}

	// 删除命令列表
	_, err = botClient.DeleteMyCommands(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("bot commands deleted successful")
}

func TestGetMyCommands(t *testing.T) {
	botClient, err := bot.New(token, httpProxyOption())
	if err != nil {
		t.Fatalf("bot.New: %v", err)
	}

	commands, err := botClient.GetMyCommands(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, cmd := range commands {
		log.Printf("Command: /%s - %s\n", cmd.Command, cmd.Description)
	}
}
package main

import (
	"log"
	"testing"

	tele "gopkg.in/telebot.v4"
)

// https://core.telegram.org/bots/api#botcommand
// TestSetMyCommands 设置 bot 命令列表
// TestSetMyCommands sets the bot command list
func TestSetMyCommands(t *testing.T) {
	// 设置命令列表
	// Set command list
	err := b.SetCommands([]tele.Command{
		{Text: "start", Description: "Start the bot"},
		{Text: "help", Description: "Show help info"},
		{Text: "profile", Description: "Show your profile"},
		{Text: "ping", Description: "Check bot status"},
	})
	if err != nil {
		t.Fatalf("SetCommands: %v", err)
	}

	log.Println("bot commands set successful")
}

// TestDeleteMyCommands 删除 bot 命令列表
// TestDeleteMyCommands deletes the bot command list
func TestDeleteMyCommands(t *testing.T) {
	// 删除命令列表
	// Delete command list
	if err := b.DeleteCommands(); err != nil {
		t.Fatalf("DeleteCommands: %v", err)
	}

	log.Println("bot commands deleted successful")
}

// TestGetMyCommands 获取 bot 命令列表
// TestGetMyCommands gets the bot command list
func TestGetMyCommands(t *testing.T) {
	// 获取命令列表
	// Get command list
	commands, err := b.Commands()
	if err != nil {
		t.Fatalf("Commands: %v", err)
	}

	for _, cmd := range commands {
		t.Logf("Command: /%s - %s\n", cmd.Text, cmd.Description)
	}
}

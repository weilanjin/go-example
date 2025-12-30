package main

import (
	"encoding/json"
	"fmt"
	"testing"

	tele "gopkg.in/telebot.v4"
)

// TestSetMenuButtonWebApp 设置菜单按钮为 WebApp
// TestSetMenuButtonWebApp sets the menu button to a WebApp
func TestSetMenuButtonWebApp(t *testing.T) {
	// 设置菜单按钮
	// Set menu button
	err := b.SetMenuButton(nil, &tele.MenuButton{
		Type:   tele.MenuButtonWebApp,
		Text:   "Open",
		WebApp: &tele.WebApp{URL: "https://bear777.win/uat/"},
	})
	if err != nil {
		t.Fatalf("SetMenuButton: %v", err)
	}
}

// TestGetMenuButton 获取菜单按钮
// TestGetMenuButton gets the menu button
func TestGetMenuButton(t *testing.T) {
	// 获取菜单按钮
	// Get menu button
	// menuButton, err := b.MenuButton(&tele.User{}) // bug: telebot v4 中该方法有问题，无法获取 WebApp URL
	menuButton, err := MenuButton(nil)
	if err != nil {
		t.Fatalf("MenuButton: %v", err)
	}
	t.Logf("Menu Button: %+v", menuButton.WebApp.URL)
}

// MenuButton returns the current value of the bot's menu button in a private chat,
// or the default menu button.
func MenuButton(chat *tele.User) (*tele.MenuButton, error) {
	params := map[string]interface{}{}

	// chat_id is optional
	if chat != nil {
		params["chat_id"] = chat.Recipient()
	}

	data, err := b.Raw("getChatMenuButton", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result *tele.MenuButton
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("telebot: %w", err)
	}
	return resp.Result, nil
}

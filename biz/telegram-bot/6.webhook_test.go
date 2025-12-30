package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/weilanjin/go-example/pkg/uid"
	tele "gopkg.in/telebot.v4"
)

// TestGetWebhookInfo 获取 Webhook 信息
// TestGetWebhookInfo gets the webhook information
func TestGetWebhookInfo(t *testing.T) {
	// 获取 Webhook 信息
	// Get webhook info
	webhookInfo, err := b.Webhook()
	if err != nil {
		t.Fatalf("Webhook: %v", err)
	}
	data, _ := json.MarshalIndent(webhookInfo, "", "  ")
	fmt.Printf("Webhook Info:\n %s\n", data)
}

// TestSetWebhook 设置 Webhook
// TestSetWebhook sets the webhook
func TestSetWebhook(t *testing.T) {
	// 替换为你的 Webhook URL
	// Replace with your webhook URL
	webhookURL := "https://www.baidu.com/webhook"

	// 设置 Webhook
	// Set webhook
	err := b.SetWebhook(&tele.Webhook{
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: webhookURL,
		},
		SecretToken: uid.UUID(),
	})
	if err != nil {
		t.Fatalf("SetWebhook: %v", err)
	}
	fmt.Println("Webhook set successfully")
}

// TestDeleteWebhook 删除 Webhook
// TestDeleteWebhook removes the webhook
func TestDeleteWebhook(t *testing.T) {
	// 删除 Webhook
	// Remove webhook
	err := b.RemoveWebhook(true)
	if err != nil {
		t.Fatalf("RemoveWebhook: %v", err)
	}
	fmt.Println("Webhook removed successfully")
}

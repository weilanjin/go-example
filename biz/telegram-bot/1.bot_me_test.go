package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetMe(t *testing.T) {
	me := b.Me

	data, _ := json.MarshalIndent(me, "", "  ")
	fmt.Printf("Bot Me:\n %s\n", data)

	// Output:
	/*
		{
		  "id": 8441906451,
		  "first_name": "本地调试🤖",
		  "last_name": "",
		  "username": "lancewei_bot",
		  "language_code": "",
		  "is_bot": true,
		  "is_premium": false,
		  "added_to_attachment_menu": false,
		  "active_usernames": null,
		  "emoji_status_custom_emoji_id": "",
		  "can_join_groups": true,
		  "can_read_all_group_messages": true,
		  "supports_inline_queries": false,
		  "can_connect_to_business": false,
		  "has_main_web_app": false
		}
	*/
}
package main

import (
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

func main() {
	// 构建内联键盘（确认/取消）
	// Build inline keyboard (Confirm/Cancel)
	menu := &tele.ReplyMarkup{}
	btnConfirm := menu.Data("确认", "action:login:confirm")
	btnCancel := menu.Data("取消", "action:login:cancel")
	menu.Inline(menu.Row(btnConfirm, btnCancel))

	// 处理 /start 命令，发送登录授权提示
	// Handle /start command and send login authorization prompt
	b.Handle("/start", func(c tele.Context) error {
		if err := c.Send("登录授权", menu); err != nil {
			slog.Error("Send failed", "err", err)
		}
		return nil
	})

	// 处理确认按钮回调
	// Handle confirm button callback
	b.Handle(&btnConfirm, func(c tele.Context) error {
		_ = c.Delete()
		if err := c.Send("登录已确认"); err != nil {
			slog.Error("Send failed", "err", err)
		}
		return nil
	})

	// 处理取消按钮回调
	// Handle cancel button callback
	b.Handle(&btnCancel, func(c tele.Context) error {
		if err := c.Send("登录已取消"); err != nil {
			slog.Error("Send failed", "err", err)
		}
		return nil
	})

	// 启动 Bot（阻塞）
	// Start the bot (blocking)
	b.Start()
}

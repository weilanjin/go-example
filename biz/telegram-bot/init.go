package main

import (
	"log"
	"log/slog"
	"net/http"
	"net/url"

	tele "gopkg.in/telebot.v4"
)

const token = "8441906451:AAGMpRGiyFi3HRe-06cfchlqKf8pmlS-OdA" // @lancewei_bot

var (
	b *tele.Bot
)

func init() {
	intBot()
}

func intBot() {
	proxyURL, _ := url.Parse("http://127.0.0.1:7890")

	pref := tele.Settings{
		Token:   token,
		Offline: false,
		OnError: func(err error, c tele.Context) {
			slog.Error("telegram bot error", "err", err, "bot", c.Bot())
		},
		Client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		},
	}

	var err error
	b, err = tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %v", err)
	}
}

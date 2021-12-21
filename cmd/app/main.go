package main

import (
	"log"

	"github.com/Elren44/elren_bot/configs"
	"github.com/Elren44/elren_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := configs.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	tbot := telegram.NewBot(bot, cfg)

	tbot.Start()

}

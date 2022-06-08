package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Elren44/elren_bot/configs"
	"github.com/Elren44/elren_bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	go func() {
		herokuUp()
	}()
	go func() {
		heroku()
	}()

	cfg, err := configs.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true

	tbot := telegram.NewBot(bot, cfg)

	tbot.Start()

}

func herokuUp() {
	ticker := time.NewTicker(5 * time.Minute)
	req, err := http.NewRequest("GET", "https://ya.ru", nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.DefaultClient

	for {
		select {
		case _ = <-ticker.C:
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("ticker -", resp.StatusCode)
		}
	}
}

func heroku() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port if not specified
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

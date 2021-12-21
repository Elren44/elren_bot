package telegram

import (
	"github.com/Elren44/elren_bot/pkg/grabbing"
	"log"

	"github.com/Elren44/elren_bot/configs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	config          *configs.Config
	searchParameter *grabbing.SearchParameter
}

func NewBot(bot *tgbotapi.BotAPI, cfg *configs.Config) *Bot {
	s := grabbing.NewSearchParameter("", "")
	return &Bot{bot: bot, config: cfg, searchParameter: s}
}

func (b *Bot) initChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	return updates
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommands(update.Message); err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())

				b.bot.Send(msg)
			}
			continue
		}
		if update.Message.ReplyToMessage != nil {
			if err := b.handleTorrentSearch(update.Message); err != nil {
				if err.Error() != "Bad Request: message text is empty" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())

					b.bot.Send(msg)
				}
				testMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "выберите что искать")
				testMsg.ReplyMarkup = menuKeyboiard
				if _, err := b.bot.Send(testMsg); err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())

					b.bot.Send(msg)
				}
			}
			continue
		}
		if update.Message.Text != "" {
			if err := b.handleCallback(update); err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())

				b.bot.Send(msg)
			}
			continue
		}
	}
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates := b.initChannel()

	b.handleUpdates(updates)
}

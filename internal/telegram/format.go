package telegram

import (
	"github.com/Elren44/elren_bot/internal/grabbing"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) sendFilms(message *tgbotapi.Message, films []grabbing.Film) error {
	var str string
	var err error
	for _, film := range films {
		str, err = grabbing.SprintFilm(film)
		if err != nil {
			continue
		}
		msgFilms := tgbotapi.NewMessage(message.Chat.ID, str)
		if _, err := b.bot.Send(msgFilms); err != nil {
			return err
		}
	}
	return nil
}

package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/Elren44/elren_bot/pkg/cdn"
	"github.com/Elren44/elren_bot/pkg/grabbing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart  = "start"
	commandSearch = "s"
	commandHelp   = "help"
)

var menuKeyboiard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("иностр. мульт"),
		tgbotapi.NewKeyboardButton("русский мульт"),
		tgbotapi.NewKeyboardButton("русский фильм"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("иностр. фильм"),
		tgbotapi.NewKeyboardButton("сериал"),
		tgbotapi.NewKeyboardButton("все видео"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("онлайн"),
		tgbotapi.NewKeyboardButton("игры"),
		tgbotapi.NewKeyboardButton("поиск везде"),
	),
)

//handleCallback requesting for a title depending on the search type
func (b *Bot) handleCallback(update tgbotapi.Update) error {
	delmsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	b.bot.Send(delmsg)
	switch update.Message.Text {
	case "поиск везде":
		if err := b.sendReply(update, "поиск везде"); err != nil {
			return err
		}
		return nil
	case "игры":
		if err := b.sendReply(update, "игры"); err != nil {
			return err
		}
		return nil
	case "онлайн":
		if err := b.sendReply(update, "онлайн"); err != nil {
			return err
		}
		return nil
	case "все видео":
		if err := b.sendReply(update, "все видео"); err != nil {
			return err
		}
		return nil
	case "сериал":
		if err := b.sendReply(update, "сериал"); err != nil {
			return err
		}
		return nil
	case "иностр. фильм":
		if err := b.sendReply(update, "иностр. фильм"); err != nil {
			return err
		}
		return nil
	case "русский фильм":
		if err := b.sendReply(update, "русский фильм"); err != nil {
			return err
		}
		return nil
	case "русский мульт":
		if err := b.sendReply(update, "русский мульт"); err != nil {
			return err
		}
		return nil
	case "иностр. мульт":
		if err := b.sendReply(update, "иностр. мульт"); err != nil {
			return err
		}
		return nil
	}

	return nil
}

//handleTorrentSearch get search type and text and looking for the requested information
func (b *Bot) handleTorrentSearch(message *tgbotapi.Message) error {
	b.searchParameter.Value = message.Text
	switch message.ReplyToMessage.Text {
	case "иностр. мульт":
		b.searchParameter.Type = grabbing.ForeignCartoons
	case "русский мульт":
		b.searchParameter.Type = grabbing.RussianCartoons
	case "русский фильм":
		b.searchParameter.Type = grabbing.RussianMovie
	case "иностр. фильм":
		b.searchParameter.Type = grabbing.ForeignMovie
	case "сериал":
		b.searchParameter.Type = grabbing.Serials
	case "все видео":
		b.searchParameter.Type = grabbing.AllMovie
	case "игры":
		b.searchParameter.Type = grabbing.Game
	case "поиск везде":
		b.searchParameter.Type = grabbing.AllFiles
	case "онлайн":
		if err := b.handleSearchCommand(message); err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			if _, err := b.bot.Send(msg); err != nil {
				return err
			}
		}
		return nil
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "Подождите идет поиск")
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	fmt.Println(b.searchParameter.Value, b.searchParameter.Type, message.ReplyToMessage.Text)
	films, err := grabbing.Find(b.searchParameter.Value, b.searchParameter.Type)
	if err != nil {
		return err
	}
	str := ""
	for _, f := range films {

		str, err = grabbing.SprintFilm(f)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ничего не найдено, проверьте название, попробуйте искать во всех видео, англ языке или онлайн")
			if _, err := b.bot.Send(msg); err != nil {
				return err
			}
		}
		msgFilms := tgbotapi.NewMessage(message.Chat.ID, str)
		if _, err := b.bot.Send(msgFilms); err != nil {
			return err
		}
	}
	testMsg := tgbotapi.NewMessage(message.Chat.ID, "выберите категорию для поиска")
	testMsg.ReplyMarkup = menuKeyboiard
	if _, err := b.bot.Send(testMsg); err != nil {
		return err
	}
	return nil
}

//sendReply prompts the user to enter a title
func (b *Bot) sendReply(update tgbotapi.Update, searchType string) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, searchType)
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "введите название",
		Selective:             false,
	}
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleCommands(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandSearch:
		return b.handleSearchCommand(message)
	case commandHelp:
		return b.handleHelpCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	testMsg := tgbotapi.NewMessage(message.Chat.ID, "выберите категорию для поиска")
	testMsg.ReplyMarkup = menuKeyboiard
	if _, err := b.bot.Send(testMsg); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "введите команду /start чтобы показать меню поиска")
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

//handleSearchCommand search online movie on the Videocdn
func (b *Bot) handleSearchCommand(message *tgbotapi.Message) error {
	movie, err := cdn.Videocdn(message.Text, b.config)
	fmt.Println(message.Text, "онлайн")
	if err != nil {
		return err
	}
	if len(movie.Data) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ничего не найдено, проверьте название или попробуйте искать на англ языке")
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
	}
	for i, data := range movie.Data {

		date, err := time.Parse("2006-01-02", data.Year)
		var layout string = "2 Jan 2006"
		if err != nil {
			return err
		}
		d := date.Format(layout)

		link := strings.TrimLeft(data.Link, "//")
		m := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.Chat.ID,
				ReplyToMessageID: 0,
			},
			ParseMode: "html",
			Text:      "<a href=\"" + link + "\">" + data.Title + "</a>" + " - " + d,
		}
		if _, err := b.bot.Send(m); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		imgLink := cdn.GrabPoster(data.ImdbID)
		img := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(imgLink))
		if _, err := b.bot.Send(img); err != nil {
			return nil
		}

		if i != len(movie.Data)-1 {
			separator := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           message.Chat.ID,
					ReplyToMessageID: 0,
				},
				ParseMode: "html",
				Text:      "&#9660",
			}
			if _, err := b.bot.Send(separator); err != nil {
				return nil
			}
		}

	}
	testMsg := tgbotapi.NewMessage(message.Chat.ID, "выберите категорию для поиска")
	testMsg.ReplyMarkup = menuKeyboiard
	if _, err := b.bot.Send(testMsg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "неизвестная команда")

	if _, err := b.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

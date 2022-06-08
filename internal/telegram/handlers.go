package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/Elren44/elren_bot/internal/grabbing"
	"github.com/Elren44/elren_bot/internal/video_db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart  = "start"
	commandSearch = "s"
	commandHelp   = "help"

	//button const
	FOREIGH_CARTOON  = "иностр. мульт"
	RUS_CARTOON      = "русский мульт"
	RUS_MOVIE        = "русский фильм"
	FOREIGH_MOVIE    = "иностр. фильм"
	TV_SERIAL        = "сериал"
	ALL_VIDEO        = "все видео"
	ONLINE           = "онлайн"
	GAMES            = "игры"
	SEARCH_EVERYWERE = "поиск везде"
)

var menuKeyboiard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(FOREIGH_CARTOON),
		tgbotapi.NewKeyboardButton(RUS_CARTOON),
		tgbotapi.NewKeyboardButton(RUS_MOVIE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(FOREIGH_MOVIE),
		tgbotapi.NewKeyboardButton(TV_SERIAL),
		tgbotapi.NewKeyboardButton(ALL_VIDEO),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ONLINE),
		tgbotapi.NewKeyboardButton(GAMES),
		tgbotapi.NewKeyboardButton(SEARCH_EVERYWERE),
	),
)

//handleCallback requesting for a title depending on the search type
func (b *Bot) handleCallback(update tgbotapi.Update) error {
	delmsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	b.bot.Send(delmsg)
	switch update.Message.Text {
	case FOREIGH_CARTOON:
		if err := b.sendReply(update, FOREIGH_CARTOON); err != nil {
			return err
		}
		return nil
	case RUS_CARTOON:
		if err := b.sendReply(update, RUS_CARTOON); err != nil {
			return err
		}
		return nil
	case RUS_MOVIE:
		if err := b.sendReply(update, RUS_MOVIE); err != nil {
			return err
		}
		return nil
	case FOREIGH_MOVIE:
		if err := b.sendReply(update, FOREIGH_MOVIE); err != nil {
			return err
		}
		return nil
	case TV_SERIAL:
		if err := b.sendReply(update, TV_SERIAL); err != nil {
			return err
		}
		return nil
	case ALL_VIDEO:
		if err := b.sendReply(update, ALL_VIDEO); err != nil {
			return err
		}
		return nil
	case ONLINE:
		if err := b.sendReply(update, ONLINE); err != nil {
			return err
		}
		return nil
	case GAMES:
		if err := b.sendReply(update, GAMES); err != nil {
			return err
		}
		return nil
	case SEARCH_EVERYWERE:
		if err := b.sendReply(update, SEARCH_EVERYWERE); err != nil {
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
	case FOREIGH_CARTOON:
		b.searchParameter.Type = grabbing.ForeignCartoons
	case RUS_CARTOON:
		b.searchParameter.Type = grabbing.RussianCartoons
	case RUS_MOVIE:
		b.searchParameter.Type = grabbing.RussianMovie
	case FOREIGH_MOVIE:
		b.searchParameter.Type = grabbing.ForeignMovie
	case TV_SERIAL:
		b.searchParameter.Type = grabbing.Serials
	case ALL_VIDEO:
		b.searchParameter.Type = grabbing.AllMovie
	case GAMES:
		b.searchParameter.Type = grabbing.Game
	case SEARCH_EVERYWERE:
		b.searchParameter.Type = grabbing.AllFiles
	case ONLINE:
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
	movie, err := video_db.Videocdn(message.Text, b.config)
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
		imgLink := video_db.GrabPoster(data.ImdbID)
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

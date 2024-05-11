package service

import (
	"blog-app/internal/models"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type ITelegramService interface {
	Send(post *models.Post) error
}

type TelegramService struct {
	Bot    *gotgbot.Bot
	ChatID int64
	Alive  bool
}

func (t *TelegramService) Send(post *models.Post) error {
	text := fmt.Sprintf("<b>%s</b>\n\n%s", post.Title, post.Content)
	_, err := t.Bot.SendMessage(t.ChatID, text, &gotgbot.SendMessageOpts{ParseMode: "html"})
	if err != nil {
		return err
	}
	return nil
}

func NewTelegramService(bot *gotgbot.Bot, chatID int64) *TelegramService {
	return &TelegramService{
		Bot:    bot,
		ChatID: chatID,
	}
}

package service

import (
	"blog-app/internal/models"
	"bytes"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"io"
	"os"
)

type TelegramService struct {
	Bot    *gotgbot.Bot
	ChatID int64
	Alive  bool
}

func (t *TelegramService) Send(post *models.Post) (*gotgbot.Message, error) {
	text := fmt.Sprintf("<b>%s</b>\n\n%s", post.Title, post.Content)
	var msg *gotgbot.Message
	if post.Attachment != nil {
		var buf bytes.Buffer
		file, _ := os.Open(fmt.Sprintf("backend/attachments/%s", *post.Attachment))
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)
		_, err := io.Copy(&buf, file)
		if err != nil {
			return nil, err
		}
		_, err = t.Bot.SendPhoto(t.ChatID, buf.Bytes(), &gotgbot.SendPhotoOpts{
			Caption:   text,
			ParseMode: "html",
		})
		buf.Reset()
		if err != nil {
			return nil, err
		}
	} else {
		message, err := t.Bot.SendMessage(t.ChatID, text, &gotgbot.SendMessageOpts{ParseMode: "html"})
		if err != nil {
			return nil, err
		}
		msg = message
	}
	return msg, nil
}

func (t *TelegramService) Delete(id int64) error {
	ok, err := t.Bot.DeleteMessage(t.ChatID, id, &gotgbot.DeleteMessageOpts{})
	if err != nil || !ok {
		return err
	}
	return nil
}

func NewTelegramService(bot *gotgbot.Bot, chatID int64) *TelegramService {
	return &TelegramService{
		Bot:    bot,
		ChatID: chatID,
		Alive:  true,
	}
}

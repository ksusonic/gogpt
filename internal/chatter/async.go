package chatter

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/models"
)

func (c *Chatter) Worker(
	yaARTResponse <-chan models.YaARTResponse,
) {
	for artResponse := range yaARTResponse {
		c.processYaARTResponse(artResponse)
	}
}

func (c *Chatter) processYaARTResponse(response models.YaARTResponse) {
	log := c.log.With(
		zap.Int64("chat_id", response.ChatID),
		zap.Int("message_id", response.MessageID),
	)

	if response.Err != nil {
		switch {
		case errors.Is(response.Err, models.GenerationTimeoutErr):
			message := tgbotapi.NewMessage(response.ChatID, "Таймаут обработки запроса. Попробуйте еще раз 🌸")
			message.ReplyToMessageID = response.MessageID
			_, _ = c.bot.Send(message)
		default:
			c.somethingWrong(response.ChatID, response.MessageID, response.Err)
		}

		return
	}

	photo := tgbotapi.NewPhoto(response.ChatID, tgbotapi.FileBytes{
		Name:  "generation result",
		Bytes: response.Image,
	})
	photo.ReplyToMessageID = response.MessageID

	_, err := c.bot.Send(photo)
	if err != nil {
		log.Error("send generated photo message", zap.Error(err))
		c.somethingWrong(response.ChatID, response.MessageID, err)

		return
	}
}

func (c *Chatter) somethingWrong(chatID int64, messageID int, err error) {
	message := tgbotapi.NewMessage(chatID, fmt.Sprintf("☠️Извините, что-то пошло не так: \n\n%+v", err))
	message.ReplyToMessageID = messageID

	_, _ = c.bot.Send(message)
}

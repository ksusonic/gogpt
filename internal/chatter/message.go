package chatter

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/models"
)

func (c *Chatter) HandleMessage(message *tgbotapi.Message) error {
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			_, _ = c.bot.Send(tgbotapi.NewMessage(
				message.Chat.ID,
				"👋 Привет! Все, что ты далее пишешь мне отправится в модель YandexART, так что пожалуйста не присылай ничего лишнего 🙌",
			))
		case "art":
			prompt := message.CommandArguments()
			if prompt == "" {
				_, _ = c.bot.Send(tgbotapi.NewMessage(
					message.Chat.ID,
					"Напиши запрос после команды: /art синий кит",
				))
				return nil
			}

			return c.handleArtRequest(message, prompt)
		default:
			_, _ = c.bot.Send(tgbotapi.NewMessage(
				message.Chat.ID,
				"🤖 Неизвестная команда",
			))
		}

		return nil
	}

	if message.Text != "" {
		return c.handleArtRequest(message, message.Text)
	}

	return nil
}

func (c *Chatter) handleArtRequest(message *tgbotapi.Message, prompt string) error {
	c.log.Info(
		prompt,
		zap.String("username", message.From.UserName),
		zap.Int64("user_id", message.From.ID),
		zap.Int64("chat_id", message.Chat.ID),
	)

	var queuePrompt string
	if len(c.artGenerateChan) > 0 {
		queuePrompt = fmt.Sprintf("Вы %d в очереди", len(c.artGenerateChan)+1)
	}

	_, err := c.bot.Send(tgbotapi.NewMessage(
		message.Chat.ID,
		fmt.Sprintf("⭐️ Отправляю запрос в YandexART. %s", queuePrompt),
	))
	if err != nil {
		return err
	}

	c.artGenerateChan <- models.YaARTRequest{
		UserName:  message.From.UserName,
		Prompt:    prompt,
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
	}

	return nil
}

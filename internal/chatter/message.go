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
			if message.CommandArguments() == "" {
				_, _ = c.bot.Send(tgbotapi.NewMessage(
					message.Chat.ID,
					"Сгенерируй запрос после команды",
				))
			}
		default:
			_, _ = c.bot.Send(tgbotapi.NewMessage(
				message.Chat.ID,
				"🤖 Неизвестная команда",
			))
		}

		return nil
	}

	if message.Text == "" {
		return nil
	}

	c.log.Info(
		message.Text,
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
		Prompt:    message.Text,
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
	}

	return nil
}

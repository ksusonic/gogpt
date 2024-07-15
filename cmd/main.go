package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/models"
	"github.com/ksusonic/gogpt/internal/service/ya_art"
)

func main() {
	log := zap.Must(zap.NewDevelopment())

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic("init bot api", zap.Error(err))
	}

	bot.Debug = true

	log.Debug("authorized on account", zap.String("username", bot.Self.UserName))

	yandexART := ya_art.NewService(
		os.Getenv("CLOUD_CATALOG_ID"),
		os.Getenv("API_KEY"),
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			log.Info(
				update.Message.Text,
				zap.String("username", update.Message.From.UserName),
				zap.Int64("user_id", update.Message.From.ID),
				zap.Int64("chat_id", update.Message.Chat.ID),
			)

			_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⭐️ Отправляю запрос в YandexART"))
			if err != nil {
				log.Error("send message", zap.Error(err))
				continue
			}

			go func(update *tgbotapi.Update) {
				ctxLog := log.With(
					zap.String("username", update.Message.From.UserName),
					zap.Time("started_at", time.Now()),
				)

				response, err := yandexART.Generate(update.Message.Text)
				if err != nil {
					ctxLog.Error("Generate", zap.Error(err))
					_, err = bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						errorMessageTemplate(err),
					))
					return
				}

				timeout := time.NewTimer(time.Second * 30)
				ticker := time.NewTicker(time.Second * 3)

				for {
					select {
					case <-ticker.C:
						ctxLog.Debug("checking result")

						result, err := yandexART.CheckResult(response.Id, ctxLog)
						if err != nil {
							if !errors.Is(err, models.NotReadyErr) {
								ctxLog.Error("CheckResult", zap.Error(err))
							}

							continue
						}

						photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
							Name:  "generation result",
							Bytes: result,
						})
						photo.ReplyToMessageID = update.Message.MessageID
						_, err = bot.Send(photo)
						if err != nil {
							ctxLog.Error("send generated photo message", zap.Error(err))
							message := tgbotapi.NewMessage(update.Message.Chat.ID, errorMessageTemplate(err))
							message.ReplyToMessageID = update.Message.MessageID
							_, err = bot.Send(message)
						}

						return

					case <-timeout.C:
						ctxLog.Error("generation wait timeout")
						message := tgbotapi.NewMessage(update.Message.Chat.ID, "Таймаут обработки запроса. Попробуйте еще раз 🌸")
						message.ReplyToMessageID = update.Message.MessageID
						bot.Send(message)
						return
					}
				}
			}(&update)
		}
	}
}

func errorMessageTemplate(err error) string {
	return fmt.Sprintf("☠️Извините, что-то пошло не так: \n\n%+v", err)
}

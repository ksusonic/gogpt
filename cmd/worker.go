package main

import (
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/models"
	"github.com/ksusonic/gogpt/internal/service/ya_art"
)

const (
	timeoutDuration = time.Second * 40
	retryInterval   = time.Second * 5
)

func worker(
	number int,
	bot *tgbotapi.BotAPI,
	yandexART *ya_art.Service,
	c <-chan *tgbotapi.Update,
	log *zap.Logger,
) {
	for update := range c {
		startedAt := time.Now()

		ctxLog := log.With(
			zap.Int("num_worker", number),
			zap.String("username", update.Message.From.UserName),
			zap.String("prompt", update.Message.Text),
		)
		ctxLog.Debug("processing yandexART request", zap.String("prompt", update.Message.Text))

		response, err := yandexART.Generate(update.Message.Text)
		if err != nil {
			ctxLog.Error("Generate", zap.Error(err))
			_, err = bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				errorMessageTemplate(err),
			))
			continue
		}

		timeout := time.NewTimer(timeoutDuration)
		ticker := time.NewTicker(retryInterval)

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

				goto end

			case <-timeout.C:
				ctxLog.Error("generation wait timeout")
				message := tgbotapi.NewMessage(update.Message.Chat.ID, "Таймаут обработки запроса. Попробуйте еще раз 🌸")
				message.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(message)

				goto end
			}
		}
	end:
		ctxLog.Info("ended processing", zap.Duration("took", time.Since(startedAt)))
	}
}

func errorMessageTemplate(err error) string {
	return fmt.Sprintf("☠️Извините, что-то пошло не так: \n\n%+v", err)
}

package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/cmd/worker"
	"github.com/ksusonic/gogpt/internal/service/ya_art"
)

var (
	tgToken    = os.Getenv("TELEGRAM_TOKEN")
	catalogID  = os.Getenv("CLOUD_CATALOG_ID")
	apiKey     = os.Getenv("API_KEY")
	numWorkers = 3
)

func main() {
	log := zap.Must(zap.NewDevelopment())

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic("init bot api", zap.Error(err))
	}

	log.Debug("authorized on account", zap.String("username", bot.Self.UserName))

	yandexART := ya_art.NewService(catalogID, apiKey)

	// 10 is blocking limit for generation requests
	generateChan := make(chan *tgbotapi.Update, 10)
	defer close(generateChan)

	for w := 1; w <= numWorkers; w++ {
		go worker.YaART(w, bot, yandexART, generateChan, log)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				_, err = bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"👋 Привет! Все, что ты далее пишешь мне отправится в модель YandexART, так что пожалуйста не присылай ничего лишнего 🙌",
				))

				continue
			}

			log.Info(
				update.Message.Text,
				zap.String("username", update.Message.From.UserName),
				zap.Int64("user_id", update.Message.From.ID),
				zap.Int64("chat_id", update.Message.Chat.ID),
			)

			_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⭐️ Отправляю запрос в YandexART. Ожидайте очередь"))
			if err != nil {
				log.Error("send message", zap.Error(err))
				continue
			}

			generateChan <- &update
		}
	}
}

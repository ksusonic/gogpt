package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/chatter"
	"github.com/ksusonic/gogpt/internal/models"
	"github.com/ksusonic/gogpt/internal/service/ya_art"
)

var (
	tgToken   = os.Getenv("TELEGRAM_TOKEN")
	catalogID = os.Getenv("CLOUD_CATALOG_ID")
	apiKey    = os.Getenv("API_KEY")

	maxConcurrentRequests = 10
	numWorkers            = 3
)

func main() {
	log := zap.Must(zap.NewDevelopment())

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic("init bot api", zap.Error(err))
	}

	log.Debug("authorized on account", zap.String("username", bot.Self.UserName))

	var (
		yaArtInChan  = make(chan models.YaARTRequest, maxConcurrentRequests)
		yaArtOutChan = make(chan models.YaARTResponse, maxConcurrentRequests)
	)

	chatBot := chatter.NewChatter(bot, log, yaArtInChan)
	go chatBot.Worker(yaArtOutChan)

	yandexART := ya_art.NewService(catalogID, apiKey, log)
	for w := 1; w <= numWorkers; w++ {
		go yandexART.Worker(w, yaArtInChan, yaArtOutChan)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			err := chatBot.HandleMessage(update.Message)
			if err != nil {
				log.Error("handle message", zap.Error(err), zap.Any("message", update.Message))
			}
		}
	}
}

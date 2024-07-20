package chatter

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/models"
)

type Chatter struct {
	bot *tgbotapi.BotAPI
	log *zap.Logger

	artGenerateChan chan<- models.YaARTRequest
}

func NewChatter(
	bot *tgbotapi.BotAPI,
	log *zap.Logger,
	artGenerateChan chan<- models.YaARTRequest,
) *Chatter {
	return &Chatter{
		bot:             bot,
		log:             log,
		artGenerateChan: artGenerateChan,
	}
}

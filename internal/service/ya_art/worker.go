package ya_art

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/models"
)

const (
	timeoutDuration = time.Second * 60
	retryInterval   = time.Second * 5
)

func (s *Service) Worker(
	number int,
	in <-chan models.YaARTRequest,
	out chan<- models.YaARTResponse,
) {
	for msg := range in {
		out <- s.processRequest(number, &msg)
	}
}

func (s *Service) processRequest(number int, msg *models.YaARTRequest) models.YaARTResponse {
	log := s.log.With(
		zap.Int("num_worker", number),
		zap.String("username", msg.UserName),
		zap.String("prompt", msg.Prompt),
	)

	startedAt := time.Now()
	defer func() {
		log.Info("art process end", zap.Duration("took", time.Since(startedAt)))
	}()

	out := models.YaARTResponse{
		ChatID:    msg.ChatID,
		MessageID: msg.MessageID,
	}

	response, err := s.Generate(msg.Prompt)
	if err != nil {
		log.Error("Generate", zap.Error(err))

		out.Err = err
		return out
	}

	var (
		timeout = time.NewTimer(timeoutDuration)
		ticker  = time.NewTicker(retryInterval)
	)

	for {
		select {
		case <-ticker.C:
			log.Debug("checking result")

			result, err := s.CheckResult(response.Id, log)
			if err != nil {
				if !errors.Is(err, models.GenerationNotReadyErr) {
					log.Error("CheckResult", zap.Error(err))
				}

				continue
			}

			out.Image = result
			return out

		case <-timeout.C:
			log.Error("generation wait timeout")

			out.Err = models.GenerationTimeoutErr
			return out
		}
	}
}

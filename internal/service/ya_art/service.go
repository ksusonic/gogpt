package ya_art

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	client               *http.Client
	yandexCloudCatalogID string
	apiKey               string
	log                  *zap.Logger
}

func NewService(
	yandexCloudCatalogID, apiKey string,
	log *zap.Logger,
) *Service {
	return &Service{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
		yandexCloudCatalogID: yandexCloudCatalogID,
		apiKey:               apiKey,
		log:                  log,
	}
}

package ya_art

import (
	"net/http"
	"time"
)

type Service struct {
	client               *http.Client
	yandexCloudCatalogID string
	apiKey               string
}

func NewService(yandexCloudCatalogID, apiKey string) *Service {
	return &Service{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
		yandexCloudCatalogID: yandexCloudCatalogID,
		apiKey:               apiKey,
	}
}

package ya_art

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ksusonic/gogpt/internal/models"
)

func (s *Service) Generate(message string) (*models.ImageGenerationAsyncResponse, error) {
	request := models.ImageGenerationAsyncRequest{
		ModelUri: fmt.Sprintf(models.YaArtModelURLTemplate, s.yandexCloudCatalogID),
		Messages: []models.PromptMessage{{Text: message}},
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", models.ImageGenerationAsyncURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create ImageGenerationAsync request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", s.apiKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST ImageGenerationAsync: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var parsedResponse models.ImageGenerationAsyncResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w: %s", err, body)
	}

	return &parsedResponse, nil
}

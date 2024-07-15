package ya_art

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/ksusonic/gogpt/internal/models"
)

func (s *Service) CheckResult(id string, log *zap.Logger) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(models.ImageCheckURLTemplate, id), nil)
	if err != nil {
		return nil, fmt.Errorf("create CheckResult request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", s.apiKey))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET ImageCheckURL: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var parsedResponse models.ImageCheckResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, fmt.Errorf("unmarshal response body: %w", err)
	}

	if !parsedResponse.Done || parsedResponse.Response == nil || parsedResponse.Response.Image == "" {
		return nil, models.NotReadyErr
	}

	decodedImage, err := base64.StdEncoding.DecodeString(parsedResponse.Response.Image)
	if err != nil {
		return nil, fmt.Errorf("decode base64 image: %w", err)
	}

	return decodedImage, nil
}

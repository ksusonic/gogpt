package models

const (
	ImageGenerationAsyncURL = "https://llm.api.cloud.yandex.net/foundationModels/v1/imageGenerationAsync"
	ImageCheckURLTemplate   = "https://llm.api.cloud.yandex.net/operations/%s"
	YaArtModelURLTemplate   = "art://%s/yandex-art/latest"
)

type ImageGenerationAsyncRequest struct {
	ModelUri string          `json:"modelUri"`
	Messages []PromptMessage `json:"messages"`
}

type PromptMessage struct {
	Text string `json:"text"`
}

type ImageGenerationAsyncResponse struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	CreatedBy   string `json:"createdBy"`
	ModifiedAt  string `json:"modifiedAt"`
	Done        bool   `json:"done"`
	Metadata    string `json:"metadata"`
	Error       struct {
		Code    string   `json:"code"`
		Message string   `json:"message"`
		Details []string `json:"details"`
	} `json:"error"`
	Response string `json:"response"`
}

type ImageCheckResponse struct {
	Id       string `json:"id"`
	Done     bool   `json:"done"`
	Response *struct {
		Type         string `json:"@type"`
		Image        string `json:"image"`
		ModelVersion string `json:"modelVersion"`
	} `json:"response"`
}

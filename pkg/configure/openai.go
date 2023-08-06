package configure

import "github.com/sashabaranov/go-openai"

type OpenAIConfig struct {
	ApiKey    string                          `json:"api_key"`
	Model     string                          `json:"model"`
	MaxTokens *int64                          `json:"max_tokens"`
	Prompts   *[]openai.ChatCompletionMessage `json:"prompts"`
}

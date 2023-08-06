package configure

type OpenAIConfig struct {
	ApiKey    string `json:"api_key"`
	Model     string `json:"model"`
	MaxTokens *int64 `json:"max_tokens"`
}

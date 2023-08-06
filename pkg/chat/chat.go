package chat

import "github.com/hanchchch/hey/pkg/configure"

type FunctionCall struct {
	Name      string
	Arguments []string
}

type ChatIO interface {
	Chat(string, *func(content string) error) (*string, *FunctionCall, error)
}

func NewChatIO(config configure.ModelConfig) ChatIO {
	if config.OpenAI != nil {
		return &OpenAIChatIO{config: *config.OpenAI}
	}
	return nil
}

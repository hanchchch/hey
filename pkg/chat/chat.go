package chat

import "github.com/hanchchch/hey/pkg/configure"

type FunctionCall struct {
	Name      string
	Arguments []string
}

type ChatIO interface {
	Response() chan string
	Chat(string) (*FunctionCall, error)
}

func NewChatIO(config configure.ModelConfig) ChatIO {
	if config.OpenAI != nil {
		return &OpenAIChatIO{config: *config.OpenAI, stream: make(chan string)}
	}
	return nil
}

func ListenResponse(chatIo ChatIO, callback func(string)) {
	for response := range chatIo.Response() {
		callback(response)
	}
}

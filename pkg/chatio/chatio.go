package chatio

import (
	"strings"
	"time"

	"github.com/hanchchch/hey/pkg/configure"
	"github.com/hanchchch/hey/pkg/llm"
)

type ChatIO struct {
	llm            llm.LLM
	line           chan string
	message        chan string
	debounce       time.Duration
	CodingLanguage *string
}

func NewChatIO(config configure.ModelConfig, debounce time.Duration) *ChatIO {
	io := &ChatIO{
		line:           make(chan string),
		message:        make(chan string),
		debounce:       debounce,
		CodingLanguage: nil,
	}
	if config.OpenAI != nil {
		io.llm = &llm.OpenAILLM{Config: *config.OpenAI, Stream: make(chan string)}
	}
	if io.llm == nil {
		return nil
	}
	return io
}

func (c *ChatIO) WaitForMessage() string {
	message := ""
	timer := time.NewTimer(c.debounce)
	for {
		select {
		case newLine := <-c.line:
			message += newLine
			timer.Reset(c.debounce)
		case <-timer.C:
			if message != "" {
				return message
			}
		}
	}
}

func (c *ChatIO) ListenResponse(callback func(string)) {
	justStartedToCode := false
	for response := range c.llm.Response() {
		if justStartedToCode {
			justStartedToCode = false
			code := response
			c.CodingLanguage = &code
		}
		if strings.Contains(response, "``") {
			if c.CodingLanguage == nil {
				justStartedToCode = true
			} else {
				c.CodingLanguage = nil
			}
		}
		callback(response)
	}
}

func (c *ChatIO) Writeln(line string) {
	c.line <- line + "\n"
}

func (c *ChatIO) Chat(message string) (*llm.FunctionCall, error) {
	return c.llm.Chat(message)
}

package chat

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/hanchchch/hey/pkg/configure"
	"github.com/sashabaranov/go-openai"
)

type OpenAIChatIO struct {
	config configure.OpenAIConfig
	stream chan string
}

func (c *OpenAIChatIO) Response() chan string {
	return c.stream
}

func (c *OpenAIChatIO) Chat(message string) (*FunctionCall, error) {
	ctx := context.Background()
	client := openai.NewClient(c.config.ApiKey)
	request := openai.ChatCompletionRequest{
		Stream:   true,
		Model:    c.config.Model,
		Messages: []openai.ChatCompletionMessage{},
	}

	if c.config.MaxTokens != nil {
		request.MaxTokens = int(*c.config.MaxTokens)
	}

	if c.config.Prompts != nil {
		request.Messages = *c.config.Prompts
	}

	request.Messages = append(request.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	stream, err := client.CreateChatCompletionStream(ctx, request)

	if err != nil {
		return nil, err
	}
	defer stream.Close()

	funcName := ""
	funcArgumentsStr := ""
	for {
		response, err := stream.Recv()
		if err != nil {
			break
		}

		delta := response.Choices[0].Delta

		if delta.Content != "" {
			c.stream <- delta.Content
		}

		if delta.FunctionCall != nil {
			if delta.FunctionCall.Name != "" {
				funcName += delta.FunctionCall.Name
			}
			if delta.FunctionCall.Arguments != "" {
				funcArgumentsStr += delta.FunctionCall.Arguments
			}
		}
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	funcArguments := []string{}
	if funcArgumentsStr != "" {
		json.Unmarshal([]byte(funcArgumentsStr), &funcArguments)
	}

	return &FunctionCall{
		Name:      funcName,
		Arguments: funcArguments,
	}, nil
}

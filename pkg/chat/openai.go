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
}

func (c *OpenAIChatIO) Chat(message string, onContent *func(string) error) (*string, *FunctionCall, error) {
	ctx := context.Background()
	client := openai.NewClient(c.config.ApiKey)
	request := openai.ChatCompletionRequest{
		Stream: true,
		Model:  c.config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	}

	if c.config.MaxTokens != nil {
		request.MaxTokens = int(*c.config.MaxTokens)
	}

	stream, err := client.CreateChatCompletionStream(ctx, request)

	if err != nil {
		return nil, nil, err
	}
	defer stream.Close()

	content := ""
	funcName := ""
	funcArgumentsStr := ""
	for {
		response, err := stream.Recv()
		if err != nil {
			break
		}

		delta := response.Choices[0].Delta

		if delta.Content != "" {
			content += delta.Content
			if onContent != nil {
				(*onContent)(delta.Content)
			}
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
		return nil, nil, err
	}

	funcArguments := []string{}
	if funcArgumentsStr != "" {
		json.Unmarshal([]byte(funcArgumentsStr), &funcArguments)
	}

	return &content, &FunctionCall{
		Name:      funcName,
		Arguments: funcArguments,
	}, nil
}

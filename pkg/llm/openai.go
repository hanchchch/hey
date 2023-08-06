package llm

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/hanchchch/hey/pkg/configure"
	"github.com/sashabaranov/go-openai"
)

type OpenAILLM struct {
	Config  configure.OpenAIConfig
	Stream  chan string
	History []openai.ChatCompletionMessage
}

func (o *OpenAILLM) Response() chan string {
	return o.Stream
}

func (o *OpenAILLM) Chat(message string) (*FunctionCall, error) {
	ctx := context.Background()
	client := openai.NewClient(o.Config.ApiKey)
	request := openai.ChatCompletionRequest{
		Stream:   true,
		Model:    o.Config.Model,
		Messages: []openai.ChatCompletionMessage{},
	}

	if o.Config.MaxTokens != nil {
		request.MaxTokens = int(*o.Config.MaxTokens)
	}

	if o.Config.Prompts != nil {
		request.Messages = *o.Config.Prompts
	}

	o.History = append(o.History, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	request.Messages = append(request.Messages, o.History...)

	stream, err := client.CreateChatCompletionStream(ctx, request)

	if err != nil {
		return nil, err
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
			o.Stream <- delta.Content
			content += delta.Content
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

	if content != "" {
		o.History = append(o.History, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
	}

	// TODO function call result?

	return &FunctionCall{
		Name:      funcName,
		Arguments: funcArguments,
	}, nil
}

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hanchchch/hey/pkg/chat"
	"github.com/hanchchch/hey/pkg/configure"
)

func main() {
	config, err := configure.FromJSON("config.json")
	if err != nil {
		fmt.Printf("Failed to parse configuration: %s\n", err)
		os.Exit(1)
	}

	names := config.ModelNames()
	if len(names) == 0 {
		fmt.Printf("No models found. Please add a model with `hey add`.\n")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Printf("Usage: hey <model>, <ask anything>\n")
		fmt.Printf("Available models: %v\n", names)
		os.Exit(0)
	}

	name, message := "", ""
	if strings.HasPrefix(os.Args[1], ",") {
		name = os.Args[1][:len(os.Args[1])-1]
		message = strings.Join(os.Args[2:], " ")
	} else {
		name = names[0]
		message = strings.Join(os.Args[1:], " ")
	}
	modelConfig := config.ModelConfig(name)
	if modelConfig == nil {
		fmt.Printf("Model %s not found. Available models: %v\n", name, names)
		os.Exit(1)
	}

	chatIo := chat.NewChatIO(*modelConfig)
	onContent := func(content string) error {
		fmt.Printf("%s", content)
		return nil
	}
	chatIo.Chat(message, &onContent)

	fmt.Println()
	os.Exit(0)
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/hanchchch/hey/pkg/chatio"
	"github.com/hanchchch/hey/pkg/configure"
)

func scanlinesForever(io *chatio.ChatIO) {
	for {
		var message string
		fmt.Scanln(&message)
		io.Writeln(message)
	}
}

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

	name := names[0]
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	modelConfig := config.ModelConfig(name)
	if modelConfig == nil {
		fmt.Printf("Model %s not found. Available models: %v\n", name, names)
		os.Exit(1)
	}
	fmt.Printf("Starting a chat with %v\n", name)

	io := chatio.NewChatIO(*modelConfig, 500*time.Microsecond)
	if io == nil {
		fmt.Printf("Failed to create chat io with model: %v\n", name)
		os.Exit(1)
	}

	// TODO
	code := ""
	codingLanguage := ""
	go io.ListenResponse(func(response string) {
		if io.CodingLanguage == nil {
			if code != "" {
				quick.Highlight(os.Stdout, code, codingLanguage, "terminal256", "monokai")
			}
			fmt.Print(response)
			code = ""
		} else {
			codingLanguage = *io.CodingLanguage
			code += response
		}
	})

	go scanlinesForever(io)

	for {
		fmt.Print(">>> ")
		message := io.WaitForMessage()
		_, err := io.Chat(message)
		if err != nil {
			fmt.Printf("Failed to chat: %s\n", err)
			os.Exit(1)
		}
		time.Sleep(100 * time.Microsecond)
		fmt.Println()
	}
}

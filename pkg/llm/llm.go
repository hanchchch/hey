package llm

type FunctionCall struct {
	Name      string
	Arguments []string
}

type LLM interface {
	Response() chan string
	Chat(string) (*FunctionCall, error)
}

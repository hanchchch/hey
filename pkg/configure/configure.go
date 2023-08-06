package configure

import (
	"encoding/json"
	"os"
)

type ModelConfig struct {
	OpenAI *OpenAIConfig `json:"openai"`
}

type Config map[string]ModelConfig

func FromJSON(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) ModelNames() []string {
	var names []string
	for name := range *c {
		names = append(names, name)
	}
	return names
}

func (c *Config) ModelConfig(name string) *ModelConfig {
	mc, ok := (*c)[name]
	if !ok {
		return nil
	}
	return &mc
}

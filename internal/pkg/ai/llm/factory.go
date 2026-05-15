package llm

import (
	"fmt"

	"StreamCore/config"
)

// NewClient creates an LLM Client from a provider configuration.
// Currently only "openai" type is supported (including OpenAI-compatible providers).
func NewClient(p config.ProviderConfig) (Client, error) {
	switch p.Type {
	case "openai":
		return NewOpenAI(p.APIKey, p.BaseURL), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", p.Type)
	}
}

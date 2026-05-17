package agent

import "StreamCore/internal/pkg/ai/llm"

// buildContext creates the initial LLM message list from the triggering user message.
// Future: enrich with recent chat history for multi-turn context.
func buildContext(userMessage string) []llm.Message {
	return []llm.Message{
		{Role: "user", Content: userMessage},
	}
}

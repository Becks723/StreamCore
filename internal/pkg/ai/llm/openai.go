package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type openaiClient struct {
	client openai.Client
}

// NewOpenAI creates an LLM Client backed by the OpenAI API.
// apiKey and baseURL are passed directly to the OpenAI SDK.
// baseURL can be used to point at OpenAI-compatible providers (DeepSeek, Groq, etc.).
func NewOpenAI(apiKey, baseURL string) Client {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return &openaiClient{client: openai.NewClient(opts...)}
}

func (c *openaiClient) Chat(ctx context.Context, params *ChatParams) (*ChatResult, error) {
	messages := toOpenAIMessages(params.SystemPrompt, params.Messages)

	tools := toOpenAITools(params.Tools)

	req := openai.ChatCompletionNewParams{
		Model:    shared.ChatModel(params.ModelName),
		Messages: messages,
		Tools:    tools,
	}

	if params.MaxTokens > 0 {
		req.MaxTokens = openai.Int(params.MaxTokens)
	}
	if params.Temperature > 0 {
		req.Temperature = openai.Float(params.Temperature)
	}

	resp, err := c.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("openai chat: %w", err)
	}

	return fromOpenAIResponse(resp), nil
}

func toOpenAIMessages(systemPrompt string, messages []Message) []openai.ChatCompletionMessageParamUnion {
	var result []openai.ChatCompletionMessageParamUnion

	if systemPrompt != "" {
		result = append(result, openai.SystemMessage(systemPrompt))
	}

	for _, m := range messages {
		switch m.Role {
		case "user":
			result = append(result, openai.UserMessage(m.Content))
		case "assistant":
			result = append(result, openai.AssistantMessage(m.Content))
		case "tool_result":
			result = append(result, openai.ToolMessage(m.Content, m.ToolUseID))
		}
	}
	return result
}

func toOpenAITools(tools []ToolDef) []openai.ChatCompletionToolParam {
	result := make([]openai.ChatCompletionToolParam, 0, len(tools))
	for _, t := range tools {
		var params map[string]any
		if len(t.InputSchema) > 0 {
			_ = json.Unmarshal(t.InputSchema, &params)
		}
		if params == nil {
			// If no schema provided, define an empty object schema
			params = map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			}
		}
		result = append(result, openai.ChatCompletionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        t.Name,
				Description: openai.String(t.Description),
				Parameters:  shared.FunctionParameters(params),
			},
		})
	}
	return result
}

func fromOpenAIResponse(resp *openai.ChatCompletion) *ChatResult {
	result := &ChatResult{}

	if len(resp.Choices) == 0 {
		return result
	}

	choice := resp.Choices[0]
	result.Text = choice.Message.Content

	for _, tc := range choice.Message.ToolCalls {
		result.ToolCalls = append(result.ToolCalls, ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}

	result.StopUsage = resp.Usage.TotalTokens
	return result
}

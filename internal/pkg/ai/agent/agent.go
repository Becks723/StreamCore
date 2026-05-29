package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"StreamCore/internal/pkg/ai/llm"
	"StreamCore/internal/pkg/ai/mcp"
	"StreamCore/internal/pkg/ai/provider"
	"StreamCore/internal/pkg/db/ai"
	"StreamCore/internal/pkg/domain"
	"StreamCore/internal/pkg/pack"
)

const (
	maxToolTurns = 5    // max tool calling rounds per message
	maxTokens    = 4096 // max response tokens
	temperature  = 0.7  // LLM temperature
)

// Agent orchestrates the AI processing pipeline.
type Agent struct {
	db        ai.AIDatabase
	toolReg   *mcp.ToolRegistry
	providers provider.Registry
}

// New creates an Agent with the given dependencies.
func New(db ai.AIDatabase, toolReg *mcp.ToolRegistry, providers provider.Registry) *Agent {
	return &Agent{
		db:        db,
		toolReg:   toolReg,
		providers: providers,
	}
}

// ProcessMessageAsync processes a chat message and returns an AI reply if appropriate.
// Returns empty string if no reply is needed. Caller should invoke in a goroutine.
func (a *Agent) ProcessMessageAsync(ctx context.Context, botID uint, content string) (string, error) {
	runtime, err := a.loadRuntime(ctx, botID)
	if err != nil {
		return "", err
	}

	messages := buildContext(content)
	return a.runLoop(ctx, runtime.llmClient, runtime.modelName, runtime.cfg.SystemPrompt, messages, runtime.tools), nil
}

func (a *Agent) ProcessProactiveMessage(ctx context.Context, botID uint, history []*domain.GroupMessage) (string, error) {
	runtime, err := a.loadRuntime(ctx, botID)
	if err != nil {
		return "", err
	}

	messages := buildProactiveContext(history)
	systemPrompt := runtime.cfg.SystemPrompt + "\n\n" + proactiveSystemInstruction
	return a.runLoop(ctx, runtime.llmClient, runtime.modelName, systemPrompt, messages, runtime.tools), nil
}

// runLoop executes the LLM + Tool Calling loop.
func (a *Agent) runLoop(ctx context.Context, client llm.Client, modelName, systemPrompt string, messages []llm.Message, tools []llm.ToolDef) string {
	for turn := 0; turn < maxToolTurns; turn++ {
		result, err := client.Chat(ctx, &llm.ChatParams{
			ModelName:    modelName,
			SystemPrompt: systemPrompt,
			Messages:     messages,
			Tools:        tools,
			MaxTokens:    maxTokens,
			Temperature:  temperature,
		})
		if err != nil {
			log.Printf("[ai] agent loop: llm chat failed: %v", err)
			return ""
		}

		if len(result.ToolCalls) == 0 {
			return result.Text
		}

		for _, tc := range result.ToolCalls {
			log.Printf("[ai] agent loop: calling tool %s(%s)", tc.Name, tc.Arguments)
			toolResult, err := a.toolReg.ExecuteTool(ctx, tc.Name, json.RawMessage(tc.Arguments))
			if err != nil {
				log.Printf("[ai] agent loop: tool %s failed: %v", tc.Name, err)
				toolResult = fmt.Sprintf("error: %v", err)
			}
			messages = append(messages, llm.Message{
				Role:      "tool_result",
				Content:   toolResult,
				ToolUseID: tc.ID,
			})
		}

		if result.Text != "" {
			messages = append(messages, llm.Message{
				Role:    "assistant",
				Content: result.Text,
			})
		}
	}

	log.Printf("[ai] agent loop: max turns (%d) reached", maxToolTurns)
	return ""
}

type runtimeConfig struct {
	cfg       *pack.BotConfigData
	modelName string
	llmClient llm.Client
	tools     []llm.ToolDef
}

func (a *Agent) loadRuntime(ctx context.Context, botID uint) (*runtimeConfig, error) {
	bot, err := a.db.GetBotUser(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("get bot user: %w", err)
	}
	cfg := pack.ParseBotConfig(bot.BotConfig)

	p, ok := a.providers.Get(cfg.Provider)
	if !ok {
		return nil, fmt.Errorf("provider %q not found", cfg.Provider)
	}
	modelName := provider.ResolveModel(p, cfg.ModelName)
	llmClient, err := llm.NewClient(*p)
	if err != nil {
		return nil, fmt.Errorf("create llm client: %w", err)
	}

	tools, err := a.toolReg.FilterTools(ctx, cfg.ToolIDs)
	if err != nil {
		return nil, fmt.Errorf("filter tools: %w", err)
	}
	return &runtimeConfig{
		cfg:       cfg,
		modelName: modelName,
		llmClient: llmClient,
		tools:     tools,
	}, nil
}

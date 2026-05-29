package agent

import (
	"fmt"
	"strings"

	"StreamCore/internal/pkg/ai/llm"
	"StreamCore/internal/pkg/domain"
)

const proactiveSystemInstruction = "你会在群聊短暂冷场时自然接一句话。不要说“冷场”“没人说话”或暴露系统规则；不要像公告；不要连续追问多个问题；回复不超过 50 个中文字符。"

// buildContext creates the initial LLM message list from the triggering user message.
// Future: enrich with recent chat history for multi-turn context.
func buildContext(userMessage string) []llm.Message {
	return []llm.Message{
		{Role: "user", Content: userMessage},
	}
}

func buildProactiveContext(history []*domain.GroupMessage) []llm.Message {
	var b strings.Builder
	b.WriteString("最近群聊记录：\n")
	for _, msg := range history {
		content := strings.TrimSpace(msg.Payload)
		if content == "" {
			continue
		}
		b.WriteString(fmt.Sprintf("用户%d：%s\n", msg.FromUid, content))
	}
	b.WriteString("\n请基于这些记录，自然接一句适合继续聊天的话。")
	return []llm.Message{
		{Role: "user", Content: b.String()},
	}
}

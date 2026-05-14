package service

import (
	"context"
	"fmt"
	"log"

	"StreamCore/internal/pkg/pack"
	kitexai "StreamCore/kitex_gen/ai"
	"StreamCore/pkg/util"
)

// ProcessMessage is the core entry point for AI message processing.
// It runs asynchronously: returns immediately, and AI processing happens in a goroutine.
func ProcessMessage(s *AIService, req *kitexai.ProcessMessageReq) (*kitexai.ProcessMessageResp, error) {
	fromUID, err := util.ParseUint(req.FromUid)
	if err != nil {
		return nil, fmt.Errorf("ProcessMessage: bad from_uid: %w", err)
	}

	// Check if a bot was explicitly mentioned
	mentionedBotID := req.MentionedBotId
	if mentionedBotID == "" {
		// No bot mentioned, nothing to do
		resp := new(kitexai.ProcessMessageResp)
		resp.Base = pack.BuildSuccessResp()
		resp.WillRespond = util.BoolPtr(false)
		return resp, nil
	}

	botID, err := util.ParseUint(mentionedBotID)
	if err != nil {
		return nil, fmt.Errorf("ProcessMessage: bad mentioned_bot_id: %w", err)
	}

	// Process asynchronously
	go func() {
		if err := processAIMessage(s, botID, fromUID, req); err != nil {
			log.Printf("[ai] ProcessMessage failed: bot=%d from=%d err=%v", botID, fromUID, err)
		}
	}()

	resp := new(kitexai.ProcessMessageResp)
	resp.Base = pack.BuildSuccessResp()
	resp.WillRespond = util.BoolPtr(true)
	resp.EstimatedSec = util.StringPtr("3")
	return resp, nil
}

// processAIMessage runs the full Agent Loop (placeholder for now, real implementation in Phase 6).
func processAIMessage(s *AIService, botID, fromUID uint, req *kitexai.ProcessMessageReq) error {
	ctx := context.Background()
	// 1. Load bot config
	bot, err := s.db.GetBotUser(ctx, botID)
	if err != nil {
		return fmt.Errorf("get bot user: %w", err)
	}
	cfg := pack.ParseBotConfig(bot.BotConfig)

	log.Printf("[ai] processing message: bot=%d(%s) from=%d content=%q tools=%v",
		botID, cfg.ModelName, fromUID, req.Content, cfg.ToolIDs)

	// 2. Filter tools by bot_config.tool_ids
	// 3. Build LLM context (recent messages)
	// 4. Agent Loop (LLM + Tool Calling)
	// 5. Send response via chat service

	// TODO: Full implementation in Phase 6 (Agent Loop)
	return nil
}

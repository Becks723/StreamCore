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
	mentionedBotID := req.MentionedBotId
	if mentionedBotID == "" {
		resp := new(kitexai.ProcessMessageResp)
		resp.Base = pack.BuildSuccessResp()
		resp.WillRespond = util.BoolPtr(false)
		return resp, nil
	}

	botID, err := util.ParseUint(mentionedBotID)
	if err != nil {
		return nil, fmt.Errorf("ProcessMessage: bad mentioned_bot_id: %w", err)
	}

	// Process asynchronously with a clean background context
	go func() {
		ctx := context.Background()

		reply, err := s.agent.ProcessMessageAsync(ctx, botID, req.Content)
		if err != nil {
			log.Printf("[ai] ProcessMessage failed: bot=%d err=%v", botID, err)
			return
		}
		if reply == "" {
			return
		}

		groupID, parseErr := util.ParseUint(req.RoomId)
		if parseErr != nil {
			log.Printf("[ai] ProcessMessage: bad room_id: %v", parseErr)
			return
		}
		if err := sendGroupBotMessage(s, ctx, botID, groupID, reply); err != nil {
			log.Printf("[ai] ProcessMessage: send reply failed: bot=%d err=%v", botID, err)
		}
	}()

	resp := new(kitexai.ProcessMessageResp)
	resp.Base = pack.BuildSuccessResp()
	resp.WillRespond = util.BoolPtr(true)
	resp.EstimatedSec = util.StringPtr("3")
	return resp, nil
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/pack"
	kitexai "StreamCore/kitex_gen/ai"
	kitexchat "StreamCore/kitex_gen/chat"
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

		// Send reply via chat service (writes to DB, returns receiver list)
		groupID, parseErr := util.ParseUint(req.RoomId)
		if parseErr != nil {
			log.Printf("[ai] ProcessMessage: bad room_id: %v", parseErr)
			return
		}
		rpcCtx := rpccontext.WithLoginUid(ctx, botID)
		msgResp, err := s.chatClient.SendGroupMessage(rpcCtx, &kitexchat.GroupClientMsg{
			ToGroupId: util.Uint2String(groupID),
			Payload:   reply,
			Timestamp: time.Now().UnixMilli(),
		})
		if err != nil {
			log.Printf("[ai] ProcessMessage: send reply failed: bot=%d err=%v", botID, err)
			return
		}

		// Push to receivers via Redis Pub/Sub so WebSocket clients get the AI reply
		pushData, _ := json.Marshal(map[string]interface{}{
			"msg_id":    msgResp.MsgId,
			"group_id":  util.String2Uint(msgResp.ToGroupId),
			"from_uid":  util.String2Uint(msgResp.FromUid),
			"content":   msgResp.Payload,
			"timestamp": msgResp.Timestamp,
		})
		for _, receiverUidStr := range msgResp.GetReceiverUids() {
			receiverUid := util.String2Uint(receiverUidStr)
			if receiverUid == 0 {
				continue
			}
			pubMsg, _ := json.Marshal(map[string]interface{}{
				"uid":  receiverUid,
				"type": "group_message",
				"data": json.RawMessage(pushData),
			})
			if err := s.rdb.Publish(ctx, "ws:push", pubMsg).Err(); err != nil {
				log.Printf("[ai] ProcessMessage: push to uid=%d failed: %v", receiverUid, err)
			}
		}
	}()

	resp := new(kitexai.ProcessMessageResp)
	resp.Base = pack.BuildSuccessResp()
	resp.WillRespond = util.BoolPtr(true)
	resp.EstimatedSec = util.StringPtr("3")
	return resp, nil
}

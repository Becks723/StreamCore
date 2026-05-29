package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"StreamCore/internal/pkg/base/rpccontext"
	kitexchat "StreamCore/kitex_gen/chat"
	"StreamCore/pkg/util"
)

func sendGroupBotMessage(s *AIService, ctx context.Context, botID, groupID uint, content string) error {
	rpcCtx := rpccontext.WithLoginUid(ctx, botID)
	msgResp, err := s.chatClient.SendGroupMessage(rpcCtx, &kitexchat.GroupClientMsg{
		ToGroupId: util.Uint2String(groupID),
		Payload:   content,
		Timestamp: time.Now().UnixMilli(),
	})
	if err != nil {
		return fmt.Errorf("chat send group message: %w", err)
	}

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
			log.Printf("[ai] sendGroupBotMessage: push to uid=%d failed: %v", receiverUid, err)
		}
	}
	return nil
}

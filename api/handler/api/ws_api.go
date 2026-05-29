package api

import (
	"context"
	"log"
	"regexp"
	"strings"

	apichat "StreamCore/api/chat"
	"StreamCore/api/pack"
	"StreamCore/api/rpc"
	"StreamCore/internal/pkg/base/rpccontext"
	kitexai "StreamCore/kitex_gen/ai"
	kitexchat "StreamCore/kitex_gen/chat"
	kitexcommon "StreamCore/kitex_gen/common"
	"StreamCore/pkg/util"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
)

// ChatHandler .
// @router /chat [GET]
func ChatHandler(ctx context.Context, c *app.RequestContext) {
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		pack.RespParamError(c, err)
		return
	}

	upgrader := websocket.HertzUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(c *app.RequestContext) bool {
			return true
		},
	}

	service := apichat.GlobalWSService()
	err = upgrader.Upgrade(c, func(conn *websocket.Conn) {
		service.RegisterClient(uid, conn)
		defer func() {
			service.UnregisterClient(uid, conn)
			_ = conn.Close()
		}()
		rpcCtx := rpccontext.WithLoginUid(context.Background(), uid)

		for {
			msgType, raw, readErr := conn.ReadMessage()
			if readErr != nil {
				log.Printf("ws read failed uid=%d: %v", uid, readErr)
				break
			}
			if msgType != websocket.TextMessage {
				continue
			}

			var wrapper apichat.MsgWrapper
			if err = sonic.Unmarshal(raw, &wrapper); err != nil {
				log.Printf("ws unmarshal wrapper failed uid=%d: %v", uid, err)
				continue
			}

			switch wrapper.Type {
			case apichat.WSMsgType_WhisperMessage:
				handleWhisperMessage(rpcCtx, service, uid, wrapper.Data)
			case apichat.WSMsgType_GroupMessage:
				handleGroupMessage(rpcCtx, service, uid, wrapper.Data)
			default:
				log.Printf("unknown ws message type=%s uid=%d", wrapper.Type, uid)
			}
		}
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
}

func handleWhisperMessage(ctx context.Context, service *apichat.WSService, uid uint, raw []byte) {
	var msg apichat.WhisperClientMsg
	if err := sonic.Unmarshal(raw, &msg); err != nil {
		log.Printf("ws unmarshal whisper failed uid=%d: %v", uid, err)
		return
	}

	resp, err := rpc.SendWhisperMessageRPC(ctx, &kitexchat.WhisperClientMsg{
		ToUid:     util.Uint2String(msg.ToUid),
		Payload:   msg.Content,
		Timestamp: msg.Timestamp,
	})
	if err != nil {
		log.Printf("chat rpc SendWhisperMessage failed uid=%d: %v", uid, err)
		return
	}

	push := &apichat.WhisperServerMsg{
		MsgId:     resp.MsgId,
		FromUid:   util.String2Uint(resp.FromUid),
		ToUid:     util.String2Uint(resp.ToUid),
		Content:   resp.Payload,
		Timestamp: resp.Timestamp,
	}
	if err = service.PushToUser(push.ToUid, apichat.WSMsgType_WhisperMessage, push); err != nil {
		log.Printf("push whisper to receiver failed uid=%d: %v", push.ToUid, err)
	}
	if err = service.PushToUser(push.FromUid, apichat.WSMsgType_WhisperMessage, push); err != nil {
		log.Printf("push whisper ack to sender failed uid=%d: %v", push.FromUid, err)
	}

	tip := &apichat.NewMessageTip{
		ConversationType: apichat.ConversationTypeWhisper,
		FromUid:          push.FromUid,
		TargetId:         push.FromUid,
		Preview:          buildPreview(push.Content),
		Timestamp:        push.Timestamp,
	}
	if err = service.PushToUser(push.ToUid, apichat.WSMsgType_NewMessageTip, tip); err != nil {
		log.Printf("push whisper tip failed uid=%d: %v", push.ToUid, err)
	}
}

func handleGroupMessage(ctx context.Context, service *apichat.WSService, uid uint, raw []byte) {
	var msg apichat.GroupClientMsg
	if err := sonic.Unmarshal(raw, &msg); err != nil {
		log.Printf("ws unmarshal group failed uid=%d: %v", uid, err)
		return
	}

	resp, err := rpc.SendGroupMessageRPC(ctx, &kitexchat.GroupClientMsg{
		ToGroupId: util.Uint2String(msg.GroupId),
		Payload:   msg.Content,
		Timestamp: msg.Timestamp,
	})
	if err != nil {
		log.Printf("chat rpc SendGroupMessage failed uid=%d: %v", uid, err)
		return
	}

	push := &apichat.GroupServerMsg{
		MsgId:     resp.MsgId,
		GroupId:   util.String2Uint(resp.ToGroupId),
		FromUid:   util.String2Uint(resp.FromUid),
		Content:   resp.Payload,
		Timestamp: resp.Timestamp,
	}

	if err = service.PushToUser(push.FromUid, apichat.WSMsgType_GroupMessage, push); err != nil {
		log.Printf("push group ack failed uid=%d: %v", push.FromUid, err)
	}
	for _, receiverUidStr := range resp.GetReceiverUids() {
		receiverUid := util.String2Uint(receiverUidStr)
		if receiverUid == 0 {
			continue
		}
		if err = service.PushToUser(receiverUid, apichat.WSMsgType_GroupMessage, push); err != nil {
			log.Printf("push group msg failed uid=%d: %v", receiverUid, err)
		}
		tip := &apichat.NewMessageTip{
			ConversationType: apichat.ConversationTypeGroup,
			FromUid:          push.FromUid,
			TargetId:         push.GroupId,
			Preview:          buildPreview(push.Content),
			Timestamp:        push.Timestamp,
		}
		if err = service.PushToUser(receiverUid, apichat.WSMsgType_NewMessageTip, tip); err != nil {
			log.Printf("push group tip failed uid=%d: %v", receiverUid, err)
		}
	}

	go triggerAIBots(rpccontext.WithLoginUid(context.Background(), uid), push.GroupId, push.FromUid, push.Content, push.Timestamp, push.MsgId)
	go triggerProactiveCheck(rpccontext.WithLoginUid(context.Background(), uid), push.GroupId, push.FromUid, push.Content, push.Timestamp, push.MsgId)
}

func triggerAIBots(ctx context.Context, groupID, fromUID uint, content string, timestamp, msgID int64) {
	botsResp, err := rpc.ListGroupBotsRPC(ctx, &kitexai.ListGroupBotsReq{GroupId: util.Uint2String(groupID)})
	if err != nil {
		log.Printf("triggerAIBots: list group bots failed group=%d: %v", groupID, err)
		return
	}
	if botsResp == nil || botsResp.Bots == nil {
		return
	}
	for _, bot := range botsResp.Bots {
		mentionedBotID := findMentionedBot(content, bot)
		if mentionedBotID == "" {
			continue
		}
		_, err := rpc.ProcessMessageRPC(ctx, &kitexai.ProcessMessageReq{
			RoomId:         util.Uint2String(groupID),
			MsgId:          msgID,
			FromUid:        util.Uint2String(fromUID),
			Content:        content,
			Timestamp:      timestamp,
			MentionedBotId: mentionedBotID,
		})
		if err != nil {
			log.Printf("triggerAIBots: ProcessMessage failed bot=%s: %v", bot.BotId, err)
		}
	}
}

func triggerProactiveCheck(ctx context.Context, groupID, fromUID uint, content string, timestamp, msgID int64) {
	_, err := rpc.ProactiveCheckRPC(ctx, &kitexai.ProactiveCheckReq{
		GroupId:   util.Uint2String(groupID),
		MsgId:     msgID,
		FromUid:   util.Uint2String(fromUID),
		Content:   content,
		Timestamp: timestamp,
	})
	if err != nil {
		log.Printf("triggerProactiveCheck: ProactiveCheck failed group=%d msg=%d: %v", groupID, msgID, err)
	}
}

func findMentionedBot(content string, bot *kitexcommon.BotInfo) string {
	// Match "@bot_name" with word boundary after the name (space, punctuation, end of string)
	pattern := `@` + regexp.QuoteMeta(bot.BotName) + `\b`
	if matched, _ := regexp.MatchString(pattern, content); matched {
		return bot.BotId
	}
	return ""
}

func buildPreview(content string) string {
	trimmed := strings.TrimSpace(content)
	runes := []rune(trimmed)
	if len(runes) > 24 {
		return string(runes[:24]) + "..."
	}
	return trimmed
}

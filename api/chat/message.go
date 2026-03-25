package chat

import "encoding/json"

const (
	WSMsgType_WhisperMessage = "whisper_message"
	WSMsgType_GroupMessage   = "group_message"
	WSMsgType_NewMessageTip  = "new_message_tip"
)

const (
	ConversationTypeWhisper = "whisper"
	ConversationTypeGroup   = "group"
)

type MsgWrapper struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// 私聊上行消息（Client -> Gateway）
type WhisperClientMsg struct {
	ToUid     uint   `json:"to_uid"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// 私聊下行消息（Gateway -> Client）
type WhisperServerMsg struct {
	MsgId     int64  `json:"msg_id"`
	FromUid   uint   `json:"from_uid"`
	ToUid     uint   `json:"to_uid"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// 群聊上行消息（Client -> Gateway）
type GroupClientMsg struct {
	GroupId   uint   `json:"group_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// 群聊下行消息（Gateway -> Client）
type GroupServerMsg struct {
	MsgId     int64  `json:"msg_id"`
	GroupId   uint   `json:"group_id"`
	FromUid   uint   `json:"from_uid"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// 新消息顶部提醒
type NewMessageTip struct {
	ConversationType string `json:"conversation_type"` // whisper/group
	FromUid          uint   `json:"from_uid"`
	TargetId         uint   `json:"target_id"` // peer uid or group id
	Preview          string `json:"preview"`
	Timestamp        int64  `json:"timestamp"`
}

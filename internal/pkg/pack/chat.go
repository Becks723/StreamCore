package pack

import (
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/chat"
	"StreamCore/pkg/util"
)

func WhisperMsg(msg *domain.WhisperMessage) *chat.WhisperServerMsg {
	if msg == nil {
		return nil
	}
	return &chat.WhisperServerMsg{
		MsgId:     msg.MsgId,
		FromUid:   util.Uint2String(msg.FromUid),
		ToUid:     util.Uint2String(msg.ToUid),
		Payload:   msg.Payload,
		Timestamp: msg.Timestamp,
	}
}

func GroupMsg(msg *domain.GroupMessage) *chat.GroupServerMsg {
	if msg == nil {
		return nil
	}
	receiverUids := make([]string, 0, len(msg.ReceiverUids))
	for _, uid := range msg.ReceiverUids {
		receiverUids = append(receiverUids, util.Uint2String(uid))
	}
	return &chat.GroupServerMsg{
		MsgId:        msg.MsgId,
		FromUid:      util.Uint2String(msg.FromUid),
		ToGroupId:    util.Uint2String(msg.GroupId),
		Payload:      msg.Payload,
		Timestamp:    msg.Timestamp,
		ReceiverUids: receiverUids,
	}
}

func WhisperHistoryData(data *domain.WhisperHistory) *chat.WhisperHistoryData {
	if data == nil {
		return nil
	}
	resp := &chat.WhisperHistoryData{
		HasMore:         data.HasMore,
		NextCursorMsgId: data.NextCursorMsgId,
	}
	resp.Items = make([]*chat.WhisperServerMsg, 0, len(data.Items))
	for _, item := range data.Items {
		resp.Items = append(resp.Items, WhisperMsg(item))
	}
	return resp
}

func WhisperHistoryAllData(items []*domain.WhisperMessage) *chat.WhisperHistoryAllData {
	resp := &chat.WhisperHistoryAllData{}
	resp.Items = make([]*chat.WhisperServerMsg, 0, len(items))
	for _, item := range items {
		resp.Items = append(resp.Items, WhisperMsg(item))
	}
	return resp
}

func GroupHistoryData(data *domain.GroupHistory) *chat.GroupHistoryData {
	if data == nil {
		return nil
	}
	resp := &chat.GroupHistoryData{
		HasMore:         data.HasMore,
		NextCursorMsgId: data.NextCursorMsgId,
	}
	resp.Items = make([]*chat.GroupServerMsg, 0, len(data.Items))
	for _, item := range data.Items {
		resp.Items = append(resp.Items, GroupMsg(item))
	}
	return resp
}

func GroupHistoryAllData(items []*domain.GroupMessage) *chat.GroupHistoryAllData {
	resp := &chat.GroupHistoryAllData{}
	resp.Items = make([]*chat.GroupServerMsg, 0, len(items))
	for _, item := range items {
		resp.Items = append(resp.Items, GroupMsg(item))
	}
	return resp
}

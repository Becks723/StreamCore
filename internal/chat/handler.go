package chat

import (
	"context"
	"fmt"

	"StreamCore/internal/chat/service"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/domain"
	"StreamCore/internal/pkg/pack"
	kitexchat "StreamCore/kitex_gen/chat"
	"StreamCore/pkg/util"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct {
	infra *base.InfraSet
}

func NewChatHandler(infra *base.InfraSet) kitexchat.ChatService {
	return &ChatServiceImpl{infra: infra}
}

func (s *ChatServiceImpl) SendWhisperMessage(ctx context.Context, msg *kitexchat.WhisperClientMsg) (*kitexchat.WhisperServerMsg, error) {
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("ChatService.SendWhisperMessage: get login uid failed: %w", err)
	}

	toUid, err := util.ParseUint(msg.ToUid)
	if err != nil {
		return nil, fmt.Errorf("ChatService.SendWhisperMessage: bad to_uid format: %w", err)
	}

	data, err := service.NewChatService(ctx, s.infra).SendWhisperMessage(uid, &domain.WhisperMessageReq{
		ToUid:     toUid,
		Payload:   msg.Payload,
		Timestamp: msg.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	return pack.WhisperMsg(data), nil
}

func (s *ChatServiceImpl) SendGroupMessage(ctx context.Context, msg *kitexchat.GroupClientMsg) (*kitexchat.GroupServerMsg, error) {
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("ChatService.SendGroupMessage: get login uid failed: %w", err)
	}

	groupId, err := util.ParseUint(msg.ToGroupId)
	if err != nil {
		return nil, fmt.Errorf("ChatService.SendGroupMessage: bad to_group_id format: %w", err)
	}

	data, err := service.NewChatService(ctx, s.infra).SendGroupMessage(uid, &domain.GroupMessageReq{
		GroupId:   groupId,
		Payload:   msg.Payload,
		Timestamp: msg.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	return pack.GroupMsg(data), nil
}

func (s *ChatServiceImpl) ListWhisperMessages(ctx context.Context, query *kitexchat.ListWhisperMessagesQuery) (*kitexchat.ListWhisperMessagesResp, error) {
	resp := new(kitexchat.ListWhisperMessagesResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListWhisperMessages: get login uid failed: %w", err)
	}

	peerUid, err := util.ParseUint(query.PeerUid)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListWhisperMessages: bad peer_uid format: %w", err)
	}

	pageSize := 0
	if query.IsSetPageSize() {
		pageSize = int(query.GetPageSize())
	}
	cursorMsgId := int64(0)
	if query.IsSetCursorMsgId() {
		cursorMsgId = query.GetCursorMsgId()
	}

	data, err := service.NewChatService(ctx, s.infra).ListWhisperMessages(uid, &domain.WhisperHistoryQuery{
		PeerUid:     peerUid,
		PageSize:    pageSize,
		CursorMsgId: cursorMsgId,
	})
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = pack.WhisperHistoryData(data)
	}
	return resp, nil
}

func (s *ChatServiceImpl) ListWhisperMessagesAll(ctx context.Context, query *kitexchat.ListWhisperMessagesAllQuery) (*kitexchat.ListWhisperMessagesAllResp, error) {
	resp := new(kitexchat.ListWhisperMessagesAllResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListWhisperMessagesAll: get login uid failed: %w", err)
	}

	peerUid, err := util.ParseUint(query.PeerUid)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListWhisperMessagesAll: bad peer_uid format: %w", err)
	}

	items, err := service.NewChatService(ctx, s.infra).ListWhisperMessagesAll(uid, peerUid)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = pack.WhisperHistoryAllData(items)
	}
	return resp, nil
}

func (s *ChatServiceImpl) ListGroupMessages(ctx context.Context, query *kitexchat.ListGroupMessagesQuery) (*kitexchat.ListGroupMessagesResp, error) {
	resp := new(kitexchat.ListGroupMessagesResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListGroupMessages: get login uid failed: %w", err)
	}

	groupId, err := util.ParseUint(query.GroupId)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListGroupMessages: bad group_id format: %w", err)
	}

	pageSize := 0
	if query.IsSetPageSize() {
		pageSize = int(query.GetPageSize())
	}
	cursorMsgId := int64(0)
	if query.IsSetCursorMsgId() {
		cursorMsgId = query.GetCursorMsgId()
	}

	data, err := service.NewChatService(ctx, s.infra).ListGroupMessages(uid, &domain.GroupHistoryQuery{
		GroupId:     groupId,
		PageSize:    pageSize,
		CursorMsgId: cursorMsgId,
	})
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = pack.GroupHistoryData(data)
	}
	return resp, nil
}

func (s *ChatServiceImpl) ListGroupMessagesAll(ctx context.Context, query *kitexchat.ListGroupMessagesAllQuery) (*kitexchat.ListGroupMessagesAllResp, error) {
	resp := new(kitexchat.ListGroupMessagesAllResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListGroupMessagesAll: get login uid failed: %w", err)
	}

	groupId, err := util.ParseUint(query.GroupId)
	if err != nil {
		return nil, fmt.Errorf("ChatService.ListGroupMessagesAll: bad group_id format: %w", err)
	}

	items, err := service.NewChatService(ctx, s.infra).ListGroupMessagesAll(uid, groupId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = pack.GroupHistoryAllData(items)
	}
	return resp, nil
}

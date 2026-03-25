package service

import (
	"errors"
	"fmt"
	"net/http"

	"StreamCore/kitex_gen/group"
	"StreamCore/pkg/util"
)

func (s *ChatService) isGroupMember(groupId, uid uint) (bool, error) {
	if s.groupClient == nil {
		return false, errors.New("group client not initialized")
	}

	resp, err := s.groupClient.IsGroupMember(s.ctx, &group.IsGroupMemberReq{
		GroupId: util.Uint2String(groupId),
		UserId:  util.Uint2String(uid),
	})
	if err != nil {
		return false, fmt.Errorf("group client IsGroupMember call failed: %w", err)
	}
	if resp == nil || resp.Base == nil {
		return false, errors.New("group client IsGroupMember got empty response")
	}
	if resp.Base.Code != http.StatusOK {
		return false, fmt.Errorf("group client IsGroupMember rejected: %s", resp.Base.Msg)
	}
	if resp.Data == nil {
		return false, errors.New("group client IsGroupMember got empty data")
	}

	return resp.Data.IsMember, nil
}

func (s *ChatService) listGroupMemberIDs(groupId uint) ([]uint, error) {
	if s.groupClient == nil {
		return nil, errors.New("group client not initialized")
	}

	resp, err := s.groupClient.ListGroupMemberIds(s.ctx, &group.ListGroupMemberIdsReq{
		GroupId: util.Uint2String(groupId),
	})
	if err != nil {
		return nil, fmt.Errorf("group client ListGroupMemberIds call failed: %w", err)
	}
	if resp == nil || resp.Base == nil {
		return nil, errors.New("group client ListGroupMemberIds got empty response")
	}
	if resp.Base.Code != http.StatusOK {
		return nil, fmt.Errorf("group client ListGroupMemberIds rejected: %s", resp.Base.Msg)
	}
	if resp.Data == nil {
		return nil, errors.New("group client ListGroupMemberIds got empty data")
	}

	memberUids := make([]uint, 0, len(resp.Data.MemberUids))
	for _, raw := range resp.Data.MemberUids {
		memberUid, parseErr := util.ParseUint(raw)
		if parseErr != nil {
			return nil, fmt.Errorf("group client ListGroupMemberIds returned invalid uid %q: %w", raw, parseErr)
		}
		memberUids = append(memberUids, memberUid)
	}
	return memberUids, nil
}

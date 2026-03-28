package pack

import (
	"strconv"

	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/group"
	"StreamCore/pkg/util"
)

func CreateGroupData(data *domain.GroupCreateResp) *group.CreateGroupData {
	if data == nil {
		return nil
	}
	return &group.CreateGroupData{
		GroupId: util.Uint2String(data.GroupId),
	}
}

func ApplyJoinGroupData(data *domain.GroupApplyResp) *group.ApplyJoinGroupData {
	if data == nil {
		return nil
	}
	return &group.ApplyJoinGroupData{
		ApplyId: strconv.FormatInt(data.ApplyId, 10),
	}
}

func IsGroupMemberData(isMember bool) *group.IsGroupMemberData {
	return &group.IsGroupMemberData{
		IsMember: isMember,
	}
}

func ListGroupMemberIdsData(memberUids []uint) *group.ListGroupMemberIdsData {
	respUids := make([]string, 0, len(memberUids))
	for _, uid := range memberUids {
		respUids = append(respUids, util.Uint2String(uid))
	}
	return &group.ListGroupMemberIdsData{
		MemberUids: respUids,
	}
}

func RespondGroupApplyData(data *domain.GroupApplyRespondResp) *group.RespondGroupApplyData {
	if data == nil {
		return nil
	}
	return &group.RespondGroupApplyData{
		ApplyId: strconv.FormatInt(data.ApplyId, 10),
		Status:  int32(data.Status),
	}
}

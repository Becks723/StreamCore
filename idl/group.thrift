namespace go group

include "common.thrift"

struct CreateGroupReq {
    1: required string name // 群名称
}

struct CreateGroupData {
    1: required string group_id // 创建后的群ID
}

struct CreateGroupResp {
    1: required common.BaseResp base
    2: optional CreateGroupData data
}

struct ApplyJoinGroupReq {
    1: required string group_id // 申请加入的群ID
    2: optional string reason // 入群申请理由
}

struct ApplyJoinGroupData {
    1: required string apply_id // 入群申请ID
}

struct ApplyJoinGroupResp {
    1: required common.BaseResp base
    2: optional ApplyJoinGroupData data
}

struct IsGroupMemberReq {
    1: required string group_id // 目标群ID
    2: required string user_id // 待校验的用户ID
}

struct IsGroupMemberData {
    1: required bool is_member // 是否为该群成员
}

struct IsGroupMemberResp {
    1: required common.BaseResp base
    2: optional IsGroupMemberData data
}

struct ListGroupMemberIdsReq {
    1: required string group_id // 目标群ID
}

struct ListGroupMemberIdsData {
    1: required list<string> member_uids // 群成员用户ID列表
}

struct ListGroupMemberIdsResp {
    1: required common.BaseResp base
    2: optional ListGroupMemberIdsData data
}

struct RespondGroupApplyReq {
    1: required string apply_id // 申请记录ID
    2: required i32 action // 1-同意 2-拒绝
}

struct RespondGroupApplyData {
    1: required string apply_id // 申请记录ID
    2: required i32 status // 1-已同意 2-已拒绝
}

struct RespondGroupApplyResp {
    1: required common.BaseResp base
    2: optional RespondGroupApplyData data
}

service GroupService {
    CreateGroupResp CreateGroup(1: required CreateGroupReq req)
    ApplyJoinGroupResp ApplyJoinGroup(1: required ApplyJoinGroupReq req)
    IsGroupMemberResp IsGroupMember(1: required IsGroupMemberReq req)
    ListGroupMemberIdsResp ListGroupMemberIds(1: required ListGroupMemberIdsReq req)
    RespondGroupApplyResp RespondGroupApply(1: required RespondGroupApplyReq req)
}

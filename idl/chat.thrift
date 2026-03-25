namespace go chat

include "common.thrift"

struct WhisperClientMsg {
    1: required string to_uid // 私聊接收方用户ID
    2: required string payload // 私聊消息内容
    3: required i64 timestamp // 客户端发送时间戳
}

struct WhisperServerMsg {
    1: required i64 msg_id // 消息唯一ID
    2: required string from_uid // 发送方用户ID
    3: required string to_uid // 接收方用户ID
    4: required string payload // 消息内容
    5: required i64 timestamp // 客户端发送时间戳
}

struct GroupClientMsg {
    1: required string to_group_id // 目标群ID
    2: required string payload // 群消息内容
    3: required i64 timestamp // 客户端发送时间戳
}

struct GroupServerMsg {
    1: required i64 msg_id // 消息唯一ID
    2: required string from_uid // 发送方用户ID
    3: required string to_group_id // 目标群ID
    4: required string payload // 消息内容
    5: required i64 timestamp // 客户端发送时间戳
    6: optional list<string> receiver_uids
}

struct ListWhisperMessagesQuery {
    1: required string peer_uid // 会话对端用户ID
    2: optional i32 page_size // 每页容量
    3: optional i64 cursor_msg_id // 查询页游标消息ID，0或不指定表示查询最新一页
}

struct WhisperHistoryData {
    1: required list<WhisperServerMsg> items // 当前页消息列表
    2: required bool has_more // 是否还有更早消息
    3: required i64 next_cursor_msg_id // 下一页游标消息ID
}

struct ListWhisperMessagesResp {
    1: required common.BaseResp base
    2: optional WhisperHistoryData data
}

struct ListWhisperMessagesAllQuery {
    1: required string peer_uid // 会话对端用户ID
}

struct WhisperHistoryAllData {
    1: required list<WhisperServerMsg> items // 全量消息列表
}

struct ListWhisperMessagesAllResp {
    1: required common.BaseResp base
    2: optional WhisperHistoryAllData data
}

struct ListGroupMessagesQuery {
    1: required string group_id // 目标群ID
    2: optional i32 page_size // 每页容量
    3: optional i64 cursor_msg_id // 查询页游标消息ID，0或不指定表示查询最新一页
}

struct GroupHistoryData {
    1: required list<GroupServerMsg> items // 当前页群消息列表
    2: required bool has_more // 是否还有更早消息
    3: required i64 next_cursor_msg_id // 下一页游标消息ID
}

struct ListGroupMessagesResp {
    1: required common.BaseResp base
    2: optional GroupHistoryData data
}

struct ListGroupMessagesAllQuery {
    1: required string group_id // 目标群ID
}

struct GroupHistoryAllData {
    1: required list<GroupServerMsg> items // 全量群消息列表
}

struct ListGroupMessagesAllResp {
    1: required common.BaseResp base
    2: optional GroupHistoryAllData data
}

service ChatService {
    WhisperServerMsg SendWhisperMessage(1: required WhisperClientMsg msg)
    GroupServerMsg SendGroupMessage(1: required GroupClientMsg msg)

    ListWhisperMessagesResp ListWhisperMessages(1: required ListWhisperMessagesQuery query)
    ListWhisperMessagesAllResp ListWhisperMessagesAll(1: required ListWhisperMessagesAllQuery query)
    ListGroupMessagesResp ListGroupMessages(1: required ListGroupMessagesQuery query)
    ListGroupMessagesAllResp ListGroupMessagesAll(1: required ListGroupMessagesAllQuery query)
}

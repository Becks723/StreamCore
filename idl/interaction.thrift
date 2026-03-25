namespace go interaction

include "common.thrift"

struct PublishLikeReq {
    1: optional string video_id
    2: optional string comment_id  // video_id 和 comment_id 二选一
    3: required i32    action_type // 1-点赞 2-取消点赞
}

struct PublishLikeResp {
    1: required common.BaseResp base
}

struct ListLikeQuery {
    1: required string user_id
    2: required i32    page_size
    3: required i32    page_num
}

struct ListLikeRespData {
    1: required list<common.VideoInfo> items
}

struct ListLikeResp {
    1: required common.BaseResp base
    2: required ListLikeRespData data
}

struct PublishCommentReq {
    1: optional string video_id
    2: optional string comment_id  // video_id 和 comment_id 二选一
    3: required string content
}

struct PublishCommentResp {
    1: required common.BaseResp base
}

struct ListCommentQuery {
    1: optional string video_id
    2: optional string comment_id  // video_id 和 comment_id 二选一
    3: required i32    page_size
    4: required i32    page_num
}

struct ListCommentRespData {
    1: required list<common.CommentInfo> items
}

struct ListCommentResp {
    1: required common.BaseResp base
    2: required ListCommentRespData data
}

struct DeleteCommentReq {
    1: required string video_id
    2: required string comment_id
}

struct DeleteCommentResp {
    1: required common.BaseResp base
}

service InteractionService {
    PublishLikeResp PublishLike(1: required PublishLikeReq req)
    ListLikeResp   ListLike(1: required ListLikeQuery req)
    PublishCommentResp PublishComment(1: required PublishCommentReq req)
    ListCommentResp ListComment(1: required ListCommentQuery query)
    DeleteCommentResp DeleteComment(1: required DeleteCommentReq req)
}
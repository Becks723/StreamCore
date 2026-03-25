namespace go social

include "common.thrift"

struct SocialList {
    1: required list<common.SocialUserInfo> items
    2: required i32 total
}

struct FollowReq {
    1: required string to_user_id
    2: required i32    action_type  // 0-关注 1-取关
}

struct FollowResp {
    1: required common.BaseResp base
}

struct ListFollowsQuery {
    1: required string user_id
    2: optional i32    page_size
    3: optional i32    page_num
}

struct ListFollowsResp {
    1: required common.BaseResp base
    2: required SocialList data
}

struct ListFollowersQuery {
    1: required string user_id
    2: optional i32    page_size
    3: optional i32    page_num
}

struct ListFollowersResp {
    1: required common.BaseResp base
    2: required SocialList        data
}

struct ListFriendsQuery {
    1: optional i32 page_size
    2: optional i32 page_num
}

struct ListFriendsResp {
    1: required common.BaseResp base
    2: required SocialList        data
}

service SocialService {
    FollowResp        Follow(1: required FollowReq req)
    ListFollowsResp   ListFollows(1: required ListFollowsQuery req)
    ListFollowersResp ListFollowers(1: required ListFollowersQuery req)
    ListFriendsResp   ListFriends(1: required ListFriendsQuery req)
}

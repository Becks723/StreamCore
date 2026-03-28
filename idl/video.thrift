namespace go video

include "common.thrift"

struct FeedQuery {
    1: optional string latest_time
}

struct FeedRespData {
    1: required list<common.VideoInfo> items
}

struct FeedResp {
    1: required common.BaseResp base
    2: optional FeedRespData data
}

struct PublishReq {
    1: optional binary data // required
    2: optional string title
    3: optional string description
    4: optional binary cover_data
}

struct PublishResp {
    1: required common.BaseResp base
}

struct ListQuery {
    1: required string user_id
    2: required i32 page_num
    3: required i32 page_size
}

struct ListRespData {
    1: required list<common.VideoInfo> items
    2: required i32 total
}

struct ListResp {
    1: required common.BaseResp base
    2: optional ListRespData data
}

struct PopularQuery {
    2: optional i32 page_num
    3: optional i32 page_size
}

struct PopularRespData {
    1: required list<common.VideoInfo> items
}

struct PopularResp {
    1: required common.BaseResp base
    2: optional PopularRespData data
}

struct SearchReq {
    1: required string keywords
    2: required i32    page_num
    3: required i32    page_size
    4: optional string from_date
    5: optional string to_date
    6: optional string username
}

struct SearchRespData {
    1: required list<common.VideoInfo> items
    2: required i32 total
}

struct SearchResp {
    1: required common.BaseResp base
    2: optional SearchRespData data
}

struct VisitQuery {
    1: required string video_id
}

struct VisitResp {
    1: required common.BaseResp base
    2: optional common.VideoInfo data
}

service VideoService {
    FeedResp    Feed(1: required FeedQuery req)
    PublishResp Publish(1: required PublishReq req)
    ListResp    List(1: required ListQuery req)
    PopularResp Popular(1: required PopularQuery req)
    SearchResp  Search(1: required SearchReq req)
    VisitResp   Visit(1: required VisitQuery req)
}

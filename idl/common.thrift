namespace go common

struct BaseResp {
    1: required i32 code
    2: required string msg
}

struct UserInfo {
    1: required string created_at
    2: required string updated_at
    3: required string deleted_at
    4: required string id
    5: required string username
    6: required string avatar_url
}

struct TokenInfo {
    1: required string access_token
    2: required string refresh_token
}

struct VideoInfo {
    1: required string created_at
    2: required string updated_at
    3: required string deleted_at
    4: required string id
    5: required string user_id
    6: required string video_url
    7: required string cover_url
    8: required string title
    9: required string description
    10: required i32   visit_count
    11: required i32   like_count
    12: required i32   comment_count
}

struct CommentInfo {
    1: required string created_at
    2: required string updated_at
    3: required string deleted_at
    4: required string id
    5: required string user_id
    6: required string video_id
    7: required string parent_id
    8: required i32   like_count
    9: required i32   child_count
    10: required string content
}

struct SocialUserInfo {
    1: required string id
    2: required string username
    3: required string avatar_url
}

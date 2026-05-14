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

struct BotInfo {
    1: required string bot_id           // Bot 对应的 user 表 ID
    2: required string bot_name         // 展示名称
    3: required string avatar_url       // 头像
    4: required string description      // 一句话描述
    5: required string system_prompt    // 人格/行为 Prompt
    6: required string model_name       // 模型标识
    7: required i32 trigger_mode        // 0: 仅 @提及触发，1: AI 自主判断（后期）
    8: required list<string> tool_ids   // 启用的 Tool ID 列表（全局 mcp_tools 的子集）
    9: required string created_at
    10: required string updated_at
}

struct MCPToolInfo {
    1: required string tool_id
    2: required string tool_name        // 如 "query_course_schedule"
    3: required string description      // 如 "查询指定学期的课程表"
    4: required string input_schema     // JSON Schema 字符串，定义参数
    5: required string server_name      // 所属 MCP Server 标识
    6: required i64 server_id           // 所属 MCP Server ID
}

struct MCPServerInfo {
    1: required i64 server_id
    2: required string server_name
    3: required string server_url
    4: required i32 sync_interval_sec   // 同步间隔（秒）
    5: required string last_synced_at
    6: required i32 status              // 1: active, 0: disabled
    7: required string created_at
    8: required string updated_at
}

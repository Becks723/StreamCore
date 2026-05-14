namespace go ai

include "common.thrift"

// ========================== 消息处理（核心） ==========================

struct ProcessMessageReq {
    1: required string room_id          // 群 ID
    2: required i64 msg_id              // 触发消息的 ID
    3: required string from_uid         // 发消息的用户
    4: required string content          // 消息文本
    5: required i64 timestamp           // 消息时间戳（ms）
    6: required string mentioned_bot_id // @了哪个 Bot，空串表示没有
}

struct ProcessMessageResp {
    1: required common.BaseResp base
    2: optional bool will_respond       // AI 是否会回复（false = 不响应此消息）
    3: optional string estimated_sec    // 预计回复等待时间
}

// ========================== Bot CRUD ==========================

struct CreateBotReq {
    1: required string bot_name
    2: required string system_prompt
    3: required string model_name
    4: required i32 trigger_mode        // 0: mention（第一期只支持 0）
    5: optional string description
    6: optional string avatar_url
    7: optional list<string> tool_ids
}

struct CreateBotResp {
    1: required common.BaseResp base
    2: optional common.BotInfo bot
}

struct UpdateBotReq {
    1: required string bot_id
    2: optional string bot_name
    3: optional string system_prompt
    4: optional string model_name
    5: optional i32 trigger_mode
    6: optional string description
    7: optional string avatar_url
    8: optional list<string> tool_ids
}

struct UpdateBotResp {
    1: required common.BaseResp base
    2: optional common.BotInfo bot
}

struct ListBotsReq {
    1: optional i32 page_size
    2: optional i32 page
}

struct ListBotsResp {
    1: required common.BaseResp base
    2: optional list<common.BotInfo> bots
    3: optional i32 total
}

struct GetBotReq {
    1: required string bot_id
}

struct GetBotResp {
    1: required common.BaseResp base
    2: optional common.BotInfo bot
}

struct DeleteBotReq {
    1: required string bot_id
}

struct DeleteBotResp {
    1: required common.BaseResp base
}

// ========================== Bot-群绑定 ==========================

struct AddBotToGroupReq {
    1: required string bot_id
    2: required string group_id
}

struct AddBotToGroupResp {
    1: required common.BaseResp base
}

struct RemoveBotFromGroupReq {
    1: required string bot_id
    2: required string group_id
}

struct RemoveBotFromGroupResp {
    1: required common.BaseResp base
}

struct ListGroupBotsReq {
    1: required string group_id
}

struct ListGroupBotsResp {
    1: required common.BaseResp base
    2: optional list<common.BotInfo> bots
}

struct ListBotGroupsReq {
    1: required string bot_id
}

struct ListBotGroupsResp {
    1: required common.BaseResp base
    2: optional list<string> group_ids
}

// ========================== MCP Server / Tool 管理 ==========================

struct RegisterMCPServerReq {
    1: required string server_name
    2: required string server_url       // HTTP/SSE endpoint
    3: optional string auth_token       // Server 认证 Token（如有）
    4: optional i32 sync_interval_sec   // 同步间隔，默认 300
}

struct RegisterMCPServerResp {
    1: required common.BaseResp base
    2: optional common.MCPServerInfo server
    3: optional list<common.MCPToolInfo> tools // 从该 Server 自动发现的 Tool 列表
}

struct RefreshMCPServerReq {
    1: required i64 server_id
}

struct RefreshMCPServerResp {
    1: required common.BaseResp base
    2: optional list<common.MCPToolInfo> tools // 刷新后的 Tool 列表
}

struct ListMCPServersReq {}

struct MCPServerListData {
    1: required list<common.MCPServerInfo> servers
}

struct ListMCPServersResp {
    1: required common.BaseResp base
    2: optional MCPServerListData data
}

struct DeleteMCPServerReq {
    1: required i64 server_id
}

struct DeleteMCPServerResp {
    1: required common.BaseResp base
}

struct ListToolsReq {}

struct ToolListData {
    1: required list<common.MCPToolInfo> tools
}

struct ListToolsResp {
    1: required common.BaseResp base
    2: optional ToolListData data
}

// ========================== 用户凭证 ==========================

struct CredentialInfo {
    1: required string credential_id
    2: required string service_name     // "jwch" / "library" 等
    3: required string username         // 学号等（脱敏展示）
    4: required string saved_at
}

struct SaveCredentialReq {
    1: required string service_name     // 如 "jwch"
    2: required string username         // 学号
    3: required string password         // 密码（TLS 传输，落库 AES 加密）
}

struct SaveCredentialResp {
    1: required common.BaseResp base
    2: optional CredentialInfo credential
}

struct DeleteCredentialReq {
    1: required string credential_id
}

struct DeleteCredentialResp {
    1: required common.BaseResp base
}

struct ListCredentialsReq {
    1: optional string service_name     // 按服务筛选，不填 = 全部
}

struct CredentialListData {
    1: required list<CredentialInfo> credentials
}

struct ListCredentialsResp {
    1: required common.BaseResp base
    2: optional CredentialListData data
}

// ========================== Service ==========================

service AIService {
    // 核心：接收聊天消息，AI 异步处理后以 Bot 身份回写群聊
    ProcessMessageResp ProcessMessage(1: required ProcessMessageReq req)

    // Bot CRUD
    CreateBotResp CreateBot(1: required CreateBotReq req)
    UpdateBotResp UpdateBot(1: required UpdateBotReq req)
    ListBotsResp  ListBots(1: required ListBotsReq req)
    GetBotResp    GetBot(1: required GetBotReq req)
    DeleteBotResp DeleteBot(1: required DeleteBotReq req)

    // Bot-群绑定
    AddBotToGroupResp    AddBotToGroup(1: required AddBotToGroupReq req)
    RemoveBotFromGroupResp RemoveBotFromGroup(1: required RemoveBotFromGroupReq req)
    ListGroupBotsResp    ListGroupBots(1: required ListGroupBotsReq req)
    ListBotGroupsResp    ListBotGroups(1: required ListBotGroupsReq req)

    // MCP Server / Tool 管理
    RegisterMCPServerResp  RegisterMCPServer(1: required RegisterMCPServerReq req)
    RefreshMCPServerResp   RefreshMCPServer(1: required RefreshMCPServerReq req)
    ListMCPServersResp     ListMCPServers(1: required ListMCPServersReq req)
    DeleteMCPServerResp    DeleteMCPServer(1: required DeleteMCPServerReq req)
    ListToolsResp          ListTools(1: required ListToolsReq req)

    // 用户凭证
    SaveCredentialResp    SaveCredential(1: required SaveCredentialReq req)
    DeleteCredentialResp  DeleteCredential(1: required DeleteCredentialReq req)
    ListCredentialsResp   ListCredentials(1: required ListCredentialsReq req)
}

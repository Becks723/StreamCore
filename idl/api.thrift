namespace go api

include "user.thrift"
include "video.thrift"
include "interaction.thrift"
include "social.thrift"
include "chat.thrift"
include "group.thrift"
include "ai.thrift"

service UserApi {
    user.RegisterResp  Register(1: required user.RegisterReq req) (api.post="/user/register")
    user.LoginResp     Login(1: required user.LoginReq req) (api.post="/user/login")
    user.InfoResp      GetInfo(1: required user.InfoQuery req) (api.get="/user/info")
    user.AvatarResp    UploadAvatar(1: required user.AvatarReq req) (api.put="/user/avatar/upload")
    user.RefreshTokenResp RefreshToken(1: required user.RefreshTokenReq req) (api.post="/auth/refresh")
    user.MFAQrcodeResp MFAQrcode(1: required user.MFAQrcodeReq req) (api.get="/auth/mfa/qrcode")
    user.MFABindResp   MFABind(1: required user.MFABindReq req) (api.post="/auth/mfa/bind")
    user.MFAVerifyResp MFAVerify(1: required user.MFAVerifyReq req) (api.post="/auth/mfa/verify")
}

service VideoApi {
    video.FeedResp    Feed(1: required video.FeedQuery req) (api.get="/video/feed")
    video.PublishResp Publish(1: required video.PublishReq req) (api.post="/video/publish")
    video.ListResp    List(1: required video.ListQuery req) (api.get="/video/list")
    video.PopularResp Popular(1: required video.PopularQuery req) (api.get="/video/popular")
    video.SearchResp  Search(1: required video.SearchReq req) (api.post="/video/search")
    video.VisitResp   Visit(1: required video.VisitQuery req) (api.get="/video/:vid")
}

service InteractionApi {
    interaction.PublishLikeResp PublishLike(1: required interaction.PublishLikeReq req) (api.post="/like/action")
    interaction.ListLikeResp   ListLike(1: required interaction.ListLikeQuery req) (api.get="/like/list")
    interaction.PublishCommentResp PublishComment(1: required interaction.PublishCommentReq req) (api.post="/comment/publish")
    interaction.ListCommentResp ListComment(1: required interaction.ListCommentQuery query) (api.get="/comment/list")
    interaction.DeleteCommentResp DeleteComment(1: required interaction.DeleteCommentReq req) (api.delete="/comment/delete")
}

service SocialApi {
    social.FollowResp        Follow(1: required social.FollowReq req) (api.post="/relation/action")
    social.ListFollowsResp   ListFollows(1: required social.ListFollowsQuery req) (api.get="/following/list")
    social.ListFollowersResp ListFollowers(1: required social.ListFollowersQuery req) (api.get="/follower/list")
    social.ListFriendsResp   ListFriends(1: required social.ListFriendsQuery req) (api.get="/friends/list")
}

service ChatApi {
    chat.ListWhisperMessagesAllResp ListWhisperMessagesAll(1: required chat.ListWhisperMessagesAllQuery req) (api.get="/chat/whisper/history/all")
    chat.ListWhisperMessagesResp ListWhisperMessages(1: required chat.ListWhisperMessagesQuery req) (api.get="/chat/whisper/history/page")
    chat.ListGroupMessagesAllResp ListGroupMessagesAll(1: required chat.ListGroupMessagesAllQuery req) (api.get="/chat/group/history/all")
    chat.ListGroupMessagesResp ListGroupMessages(1: required chat.ListGroupMessagesQuery req) (api.get="/chat/group/history/page")
}

service GroupApi {
    group.CreateGroupResp CreateGroup(1: required group.CreateGroupReq req) (api.post="/chat/group/create")
    group.ApplyJoinGroupResp ApplyJoinGroup(1: required group.ApplyJoinGroupReq req) (api.post="/chat/group/apply")
    group.RespondGroupApplyResp RespondGroupApply(1: required group.RespondGroupApplyReq req) (api.post="/chat/group/apply/respond")
}

struct PlaceholderResp {}

service WsApi {
    PlaceholderResp ChatHandler() (api.get="/chat")
}

service AIApi {
    ai.CreateBotResp CreateBot(1: required ai.CreateBotReq req) (api.post="/ai/bot")
    ai.UpdateBotResp UpdateBot(1: required ai.UpdateBotReq req) (api.put="/ai/bot")
    ai.ListBotsResp  ListBots(1: required ai.ListBotsReq req) (api.get="/ai/bots")
    ai.GetBotResp    GetBot(1: required ai.GetBotReq req) (api.get="/ai/bot/:bot_id")
    ai.DeleteBotResp DeleteBot(1: required ai.DeleteBotReq req) (api.delete="/ai/bot/:bot_id")

    ai.AddBotToGroupResp    AddBotToGroup(1: required ai.AddBotToGroupReq req) (api.post="/ai/bot/group")
    ai.RemoveBotFromGroupResp RemoveBotFromGroup(1: required ai.RemoveBotFromGroupReq req) (api.delete="/ai/bot/group")
    ai.ListGroupBotsResp    ListGroupBots(1: required ai.ListGroupBotsReq req) (api.get="/ai/group/:group_id/bots")

    ai.RegisterMCPServerResp  RegisterMCPServer(1: required ai.RegisterMCPServerReq req) (api.post="/ai/mcp/server")
    ai.RefreshMCPServerResp   RefreshMCPServer(1: required ai.RefreshMCPServerReq req) (api.post="/ai/mcp/server/refresh")
    ai.ListMCPServersResp     ListMCPServers(1: required ai.ListMCPServersReq req) (api.get="/ai/mcp/servers")
    ai.DeleteMCPServerResp    DeleteMCPServer(1: required ai.DeleteMCPServerReq req) (api.delete="/ai/mcp/server/:server_id")
    ai.ListToolsResp          ListTools(1: required ai.ListToolsReq req) (api.get="/ai/tools")

    ai.SaveCredentialResp   SaveCredential(1: required ai.SaveCredentialReq req) (api.post="/ai/credential")
    ai.DeleteCredentialResp DeleteCredential(1: required ai.DeleteCredentialReq req) (api.delete="/ai/credential/:credential_id")
    ai.ListCredentialsResp  ListCredentials(1: required ai.ListCredentialsReq req) (api.get="/ai/credentials")
}

package constants

import "time"

const (
	ApiServiceName         = "api"
	UserServiceName        = "user"
	VideoServiceName       = "video"
	InteractionServiceName = "interaction"
	SocialServiceName      = "social"
	ChatServiceName        = "chat"
	GroupServiceName       = "group"

	MaxConnections = 1000
	MaxQPS         = 100

	JWT_AccessSecret           = "access_token_secret"
	JWT_RefreshSecret          = "refresh_token_secret"
	JWT_AccessTokenExpiration  = 12 * time.Hour
	JWT_RefreshTokenExpiration = 7 * 24 * time.Hour

	MFAQrcodeWidth   = 256
	MFAQrcodeHeight  = 256
	MFATokenExpiry   = 10 * time.Minute // mfa token ttl
	TOTPSecretExpiry = 10 * time.Minute // totp密钥ttl
	TOTPInterval     = 30               // totp验证码刷新间隔，秒
	TOTPFailureLimit = 10               // 规定时间内最多允许的绑定失败次数，防爆破用
	TOTPFailureReset = 5 * time.Minute  // 防爆破的刷新时间

	MuxConnection = 1

	LikeTarType_Video   = 1
	LikeTarType_Comment = 2
	LikeAction_Like     = 1
	LikeAction_Unlike   = 2
	UserLikesCacheLimit = 200

	VideoVisitQueueSize                   = 5000
	VideoVisitBatchSize                   = 1000
	VideoVisitFlushInterval time.Duration = 10 * time.Second

	FollowAction_Follow   = 0
	FollowAction_Unfollow = 1
	SocialCacheExpiration = 30 * time.Minute

	ChatMsgTypeWhisper = 1
	ChatMsgTypeGroup   = 2

	ChatDeliverStatusPending = 0
	ChatDeliverStatusDone    = 1

	ChatGroupRoleOwner  = 1
	ChatGroupRoleMember = 2

	ChatGroupMemberStatusActive = 0
	ChatGroupMemberStatusLeft   = 1

	ChatGroupApplyStatusPending  = 0
	ChatGroupApplyStatusApproved = 1
	ChatGroupApplyStatusRejected = 2

	ChatGroupApplyActionApprove = 1
	ChatGroupApplyActionReject  = 2
)

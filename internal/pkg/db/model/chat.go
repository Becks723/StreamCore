package model

type ChatMsgModel struct {
	Id            int64  `gorm:"primaryKey"`
	SenderUid     uint   `gorm:"index:idx_chat_sender_receiver_msg,priority:1"`
	ReceiverUid   uint   `gorm:"index:idx_chat_sender_receiver_msg,priority:2"`
	GroupId       uint   `gorm:"index:idx_chat_group_msg,priority:1"`
	Payload       string `gorm:"type:text"`
	MsgType       int    `gorm:"index"` // 1-whisper, 2-group
	Timestamp     int64  `gorm:"index"`
	DeliverStatus int
}

type ChatGroupModel struct {
	Id       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:128"`
	OwnerUid uint   `gorm:"index"`
}

type ChatGroupMemberModel struct {
	Id      uint `gorm:"primaryKey"`
	GroupId uint `gorm:"index:idx_chat_group_user,unique;index:idx_chat_group_member_status,priority:1"` // Group ID that this membership belongs to
	UserId  uint `gorm:"index:idx_chat_group_user,unique;index:idx_chat_group_member_status,priority:2"` // User ID that joins the group
	Role    int  // Role of the member in the group, such as owner or normal member
	Status  int  // 0-active, 1-left
}

// ChatGroupApplyModel a model when a user requests for joining a group
type ChatGroupApplyModel struct {
	Id           int64  `gorm:"primaryKey"`
	GroupId      uint   `gorm:"index:idx_chat_group_apply_status,priority:1"`
	ApplicantUid uint   `gorm:"index:idx_chat_group_apply_status,priority:2"`
	Reason       string `gorm:"type:text"`                                    // Optional application reason message
	Status       int    `gorm:"index:idx_chat_group_apply_status,priority:3"` // 0-pending, 1-approved, 2-rejected
	CreatedAtTs  int64  `gorm:"index"`                                        // Creation timestamp in milliseconds
	UpdatedAtTs  int64  // Last update timestamp in milliseconds.
}

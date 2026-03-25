package domain

type WhisperMessageReq struct {
	ToUid     uint
	Payload   string
	Timestamp int64
}

type WhisperMessage struct {
	MsgId     int64
	FromUid   uint
	ToUid     uint
	Payload   string
	Timestamp int64
	MsgType   int
}

type GroupMessageReq struct {
	GroupId   uint
	Payload   string
	Timestamp int64
}

type GroupMessage struct {
	MsgId        int64
	FromUid      uint
	GroupId      uint
	Payload      string
	Timestamp    int64
	ReceiverUids []uint
}

type GroupCreateReq struct {
	Name string
}

type GroupCreateResp struct {
	GroupId uint
}

type GroupApplyReq struct {
	GroupId uint
	Reason  string
}

type GroupApplyResp struct {
	ApplyId int64
}

type GroupApplyRespondReq struct {
	ApplyId int64
	Action  int
}

type GroupApplyRespondResp struct {
	ApplyId int64
	Status  int
}

type WhisperHistoryQuery struct {
	PeerUid     uint
	PageSize    int
	CursorMsgId int64
}

type GroupHistoryQuery struct {
	GroupId     uint
	PageSize    int
	CursorMsgId int64
}

type WhisperHistory struct {
	Items           []*WhisperMessage
	HasMore         bool
	NextCursorMsgId int64
}

type GroupHistory struct {
	Items           []*GroupMessage
	HasMore         bool
	NextCursorMsgId int64
}

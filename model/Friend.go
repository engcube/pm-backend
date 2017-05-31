package PrivateMessageModel

// Friend 联系人
type Friend struct {
	FriendID     int
	FriendUserID int
	Email        string
	Nickname     string
	InsertTime   int64
	UpdateTime   int64
	IsDeleted    bool
	SentMsgs     []Message // 由自己发送的消息
	RecieveMsgs  []Message // 由对方发送的消息
	UnreadCount  int       // 由对方发送的未读消息
	TotalCount   int       // 双方所有消息数
  LastMessage  Message
}


package PrivateMessageModel

import (
	"fmt"
	"pm-backend/public"
	"strconv"
	"time"
)

// Message 私信
type Message struct {
	MessageID     int
	RecieverEmail string
	Sender        int
	Reciever      int
	Content       string
	IsViewed      bool
	InsertTime    int64
	UpdateTime    int64
	IsDeleted     bool
}

// New 增加Message
func (m *Message) New() error {
	res, err := PrivateMessageBackendPublic.Insert(SQL_ADD_MESSAGE, m.Sender, m.Reciever, m.Content, time.Now().Unix())
	if err != nil {
		return err
	}
	m.MessageID = int(res)
	m.InsertTime = time.Now().Unix()
	m.IsDeleted = false
	m.IsViewed = false
	return nil
}

// Read 阅读Message
func (m *Message) Read() error {
	if m.IsViewed {
		return fmt.Errorf("Message alread read")
	}
	cnt, err := PrivateMessageBackendPublic.Update(SQL_READ_MESSAGE, time.Now().Unix(), m.MessageID)
	if err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

// Delete 删除Message
func (m *Message) Delete() error {
	cnt, err := PrivateMessageBackendPublic.Update(SQL_DELETE_MESSAGE, time.Now().Unix(), m.MessageID)
	if err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

// Get 获取Message信息
func (m *Message) Get() error {
	if m.MessageID == 0 {
		return fmt.Errorf("MessageID not provided")
	}
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_MESSAGE, m.MessageID)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return fmt.Errorf("No message fetched")
	}
	uid, _ := strconv.ParseInt(rows[0][1], 10, 64)
	touid, _ := strconv.ParseInt(rows[0][2], 10, 64)
	m.Sender = int(uid)
	m.Reciever = int(touid)
	m.Content = string(rows[0][3])
	viewd, _ := strconv.ParseInt(rows[0][4], 10, 32)
	if int(viewd) > 0 {
		m.IsViewed = true
	} else {
		m.IsViewed = false
	}
	insert, _ := strconv.ParseInt(rows[0][5], 10, 64)
	update, _ := strconv.ParseInt(rows[0][6], 10, 64)
	m.InsertTime = int64(insert)
	m.UpdateTime = int64(update)
	return nil
}

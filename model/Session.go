package PrivateMessageModel

import (
	"fmt"
	"pm-backend/public"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
)

const (
	SESSION_EXPIRATION_DURATION = 60 * 60 * 24 * 15 //会话15天过期
)

// Session 会话信息
type Session struct {
	SessionID  string
	UserID     int
	UpdateTime int64
}

// New 创建新的会话
func (s *Session) New() error {
	if s.UserID == 0 {
		return fmt.Errorf("No UserID provided")
	}
	if s.SessionID == "" {
		s.SessionID = uuid.NewV4().String()
	}
	if s.UpdateTime == 0 {
		s.UpdateTime = time.Now().Unix()
	}
	_, err := PrivateMessageBackendPublic.Insert(SQL_NEW_SESSION, s.SessionID, s.UserID, s.UpdateTime, s.UpdateTime)
	if err != nil {
		return err
	}
	return nil
}

// Get 获取会话信息
func (s *Session) Get() error {
	if s.SessionID == "" {
		return fmt.Errorf("No SessionID provided")
	}
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_SESSION, s.SessionID)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return fmt.Errorf("Invalid Session")

	}
	if len(rows) > 1 {
		return fmt.Errorf("Get Session Error with multiple rows ")
	}
	userID, _ := strconv.ParseInt(string(rows[0][1]), 10, 64)
	updatetime, _ := strconv.ParseInt(string(rows[0][2]), 10, 64)
	if s.UserID != 0 {
		if s.UserID != int(userID) {
			return fmt.Errorf("invalid session")
		}
	}
	s.UserID = int(userID)
	s.UpdateTime = int64(updatetime)
	if !s.Valid() {
		return fmt.Errorf("session timeout")
	}
	return nil
}

// Delete 销毁会话
func (s *Session) Delete() error {
	if s.SessionID == "" {
		return fmt.Errorf("No SessionID provided")
	}
	_, err := PrivateMessageBackendPublic.Update(SQL_DELETE_SESSION, time.Now().Unix(), s.SessionID)
	return err
}

// Update 更新会话
func (s *Session) Update() error {
	s.UpdateTime = time.Now().Unix()
	cnt, err := PrivateMessageBackendPublic.Update(SQL_UPDATE_SESSION, s.UpdateTime, s.SessionID)
	if cnt == 0 {
		return fmt.Errorf("no row updated")
	}
	return err
}

// Valid 会话是否超时
func (s *Session) Valid() bool {
	now := time.Now().Unix()
	if now-s.UpdateTime >= SESSION_EXPIRATION_DURATION {
		return false
	}
	return true
}

package PrivateMessageModel

import (
	"fmt"
	"testing"
)

var (
	SessionID string
	UserID    = 10
)

func Test_Login(t *testing.T) {
	s := Session{UserID: UserID}
	err := s.New()
	if err != nil {
		t.Error(err)
	}
	SessionID = s.SessionID
	fmt.Println(SessionID)
}

func Test_Valid(t *testing.T) {
	s := Session{SessionID: SessionID, UserID: UserID}
	err := s.Valid()
	if err != nil {
		t.Error(err)
	}

	s2 := Session{SessionID: SessionID, UserID: 2}
	err = s2.Valid()
	if err == nil {
		t.Error("invalid")
	}
	fmt.Println(err.Error())
}

func Test_Logout(t *testing.T) {
	s := Session{SessionID: SessionID}
	err := s.Delete()
	if err != nil {
		t.Error(err)
	}
}

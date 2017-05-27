package PrivateMessageAPIV1

import (
	"net/http"
	"pm-backend/model"
	"pm-backend/public"

	"fmt"

	"github.com/ant0ine/go-json-rest/rest"
)

// ParseSession 解析session
func ParseSession(sessionID string) (int, error) {
	if sessionID == "" {
		return 0, fmt.Errorf("No sessionid in header")
	}
	session := PrivateMessageModel.Session{SessionID: sessionID}
	err := session.Get()
	if err != nil {
		return 0, err
	}
	return session.UserID, nil
}

// GetSession Get /#version/session, 获取会话信息
func GetSession(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	if sessionID == "" {
		rest.Error(w, "no SessionID in header", http.StatusInternalServerError)
		return
	}
	session := PrivateMessageModel.Session{SessionID: sessionID}
	err := session.Get()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	w.WriteJson(session)
}

// PostSession Post /#version/session, 创建新的会话（登入）
func PostSession(w rest.ResponseWriter, r *rest.Request) {
	session := PrivateMessageModel.Session{}
	user := PrivateMessageModel.User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bValid, err := user.Validate()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_USER_LOGIN)
		return
	}
	if !bValid {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_USER_PASSWORD)
		return
	}
	session.UserID = user.UserID
	err = session.New()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	user.SessionID = session.SessionID
	w.WriteJson(user)
}

// PutSession Put /#version/session, 更新会话信息
func PutSession(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	if sessionID == "" {
		rest.Error(w, "no SessionID in header", http.StatusInternalServerError)
		return
	}
	session := PrivateMessageModel.Session{SessionID: sessionID}
	err := session.Get()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	err = session.Update()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	w.WriteJson(session)
}

// DeleteSession Delete /#version/session, 销毁当前会话（登出）
func DeleteSession(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	if sessionID == "" {
		rest.Error(w, "no SessionID in header", http.StatusInternalServerError)
		return
	}
	session := PrivateMessageModel.Session{SessionID: sessionID}
	err := session.Get()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	err = session.Delete()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	w.WriteJson(session)
}

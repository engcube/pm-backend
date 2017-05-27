package PrivateMessageAPIV1

import (
	"net/http"
	"pm-backend/model"
	"pm-backend/public"

	"fmt"

	"github.com/ant0ine/go-json-rest/rest"
)

// GetUserInfo GET /api/#version/user/:id； 获取id用户的信息
func GetUserInfo(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	err = user.Get()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_USER_FETCH)
		return
	}
	w.WriteJson(user)
}

// Register POST /api/#version/user；创建新的用户（注册）
func Register(w rest.ResponseWriter, r *rest.Request) {
	user := PrivateMessageModel.User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Email == "" {
		rest.Error(w, "user email required", PrivateMessageBackendPublic.ERR_FIELD_MISSED)
		return
	}
	if user.Password == "" {
		rest.Error(w, "user password required", PrivateMessageBackendPublic.ERR_FIELD_MISSED)
		return
	}
	if user.Username == "" {
		rest.Error(w, "username reuqired", PrivateMessageBackendPublic.ERR_FIELD_MISSED)
		return
	}

	err = user.Register()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_USER_REGISTER)
		return
	}
	w.WriteJson(user)
}

// ModifyUsername PUT /api/#version/user；更新用户的信息
func ModifyUsername(w rest.ResponseWriter, r *rest.Request) {
	user, err := ValidSession(r)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	if user.Username == "" {
		rest.Error(w, "empty username", PrivateMessageBackendPublic.ERR_FIELD_MISSED)
		return
	}
	err = user.Update()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_USER_UPDATE)
		return
	}
	w.WriteJson(user)
}

// DeleteUser DELETE /api/#version/user；删除用户（注销）
func DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	user, err := ValidSession(r)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	err = user.Delete()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_USER_DELETE)
		return
	}
	session := PrivateMessageModel.Session{SessionID: user.SessionID}
	err = session.Delete()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_ERROR)
		return
	}
	w.WriteJson(user)
}

// ValidSession 验证session
func ValidSession(r *rest.Request) (*PrivateMessageModel.User, error) {
	sessionID := r.Header.Get("SessionID")
	userid, err := ParseSession(sessionID)
	if err != nil {
		return nil, err
	}
	user := PrivateMessageModel.User{SessionID: sessionID}
	err = r.DecodeJsonPayload(&user)
	if err != nil {
		return nil, err
	}
	if userid != user.UserID {
		return nil, fmt.Errorf("permission denied")
	}
	return &user, nil
}

package PrivateMessageAPIV1

import (
	"net/http"
	"pm-backend/model"
	"pm-backend/public"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// GetAllFriends GET /api/#version/user/:id； 获取id用户的信息
func GetAllFriends(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	friends, err := user.GetAllFriends()
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_FRIEND_GET)
		return
	}
	w.WriteJson(friends)
}

// GetFriend GET /api/#version/friend/:id；获取联系人信息
func GetFriend(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	fid, _ := strconv.ParseInt(r.PathParam("id"), 10, 64)
	friends, err := user.GetFriend([]int{int(fid)})
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_FRIEND_GET)
		return
	}
	w.WriteJson(friends)
}

// AddFriend POST /api/#version/friend；创建新联系人
func AddFriend(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	friend := PrivateMessageModel.Friend{}
	err = r.DecodeJsonPayload(&friend)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	err = user.AddFriend(&friend)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_FRIEND_ADD)
		return
	}
	w.WriteJson(friend)
}

// DeleteFriend DELETE /api/#version/friend；删除指定联系人
func DeleteFriend(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("SessionID")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	friend := PrivateMessageModel.Friend{}
	err = r.DecodeJsonPayload(&friend)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	err = user.DeleteFriend(&friend)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_FRIEND_ADD)
		return
	}
	w.WriteJson(friend)
}

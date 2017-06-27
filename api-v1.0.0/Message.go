package PrivateMessageAPIV1

import (
	"net/http"
	"pm-backend/model"
	"pm-backend/public"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// GetAllMessageCount GET /api/#version/message/amount；获取私信数目
func GetAllMessageCount(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("Authorization")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	friends, err := user.GetMessages([]int{})
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_MESSAGE_GET)
		return
	}
	tmps := make([]PrivateMessageModel.Friend, 0)
	for _, friend := range friends {
		friend.RecieveMsgs = make([]PrivateMessageModel.Message, 0)
		friend.SentMsgs = make([]PrivateMessageModel.Message, 0)
		tmps = append(tmps, friend)
	}
	w.WriteJson(tmps)
}

// GetMessageCount /api/#version/message/amount/:id；获取z指定用户的私信数目
func GetMessageCount(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("Authorization")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	fid, _ := strconv.ParseInt(r.PathParam("id"), 10, 64)
	friends, err := user.GetMessages([]int{int(fid)})
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_MESSAGE_GET)
		return
	}
	tmps := make([]PrivateMessageModel.Friend, 0)
	for _, friend := range friends {
		friend.RecieveMsgs = make([]PrivateMessageModel.Message, 0)
		friend.SentMsgs = make([]PrivateMessageModel.Message, 0)
		tmps = append(tmps, friend)
	}
	w.WriteJson(tmps)
}

// GetMessages GET /api/#version/message；获取所有私信信息
func GetMessages(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("Authorization")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	friends, err := user.GetMessages([]int{})
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_MESSAGE_GET)
		return
	}
	w.WriteJson(friends)
}

// GetMessage GET /api/#version/message/:id；获取指定用户的私信
func GetMessage(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("Authorization")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	fid, _ := strconv.ParseInt(r.PathParam("id"), 10, 64)
	friends, err := user.GetMessages([]int{int(fid)})
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_MESSAGE_GET)
		return
	}
  for _, k := range friends {
    for _, m := range k.RecieveMsgs {
      if !m.IsViewed {
          user.ReadMessage(&m)
      }
    }
  }
	w.WriteJson(friends)
}

// SendMessage POST /api/#version/message；发送私信
func SendMessage(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("Authorization")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	message := PrivateMessageModel.Message{}
	err = r.DecodeJsonPayload(&message)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = user.SendMessage(&message)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_MESSAGE_SEND)
		return
	}
	w.WriteJson(message)
}

// DeleteMessage DELETE /api/#version/message；删除指定私信
func DeleteMessage(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("Authorization")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	message := PrivateMessageModel.Message{}
	err = r.DecodeJsonPayload(&message)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = user.DeleteMessage(&message)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_MESSAGE_DELETE)
		return
	}
	w.WriteJson(message)
}

// ReadMessage PUT /api/#version/message；阅读发送给自己的指定私信
func ReadMessage(w rest.ResponseWriter, r *rest.Request) {
	sessionID := r.Header.Get("Authorization")
	userid, err := ParseSession(sessionID)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_SESSION_PARSE)
		return
	}
	user := PrivateMessageModel.User{UserID: userid}
	message := PrivateMessageModel.Message{}
	err = r.DecodeJsonPayload(&message)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = user.ReadMessage(&message)
	if err != nil {
		rest.Error(w, err.Error(), PrivateMessageBackendPublic.ERR_MESSAGE_READ)
		return
	}
	w.WriteJson(message)
}

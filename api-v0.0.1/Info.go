package PrivateMessageAPIV0

import (
	"pm-backend/model"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/coreos/go-semver/semver"
)

// Info 获取api信息
func Info(w rest.ResponseWriter, r *rest.Request) {
	info := PrivateMessageModel.AppInfo{Owner: "zpcpromac@gmail.com", Version: (r.Env["VERSION"].(*semver.Version)).String(), Usage: "PrivateMessage Backend Service"}
	w.WriteJson(info)
}

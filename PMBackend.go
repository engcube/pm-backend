package main

import (
	"log"
	"flag"
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/coreos/go-semver/semver"
	"pm-backend/api-v0.0.1"
	"pm-backend/api-v1.0.0"
)

// SemVerMiddleware 版本控制
type SemVerMiddleware struct {
	MinVersion string
	MaxVersion string
}

// MiddlewareFunc 版本控制
func (mw *SemVerMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	minVersion, err := semver.NewVersion(mw.MinVersion)
	if err != nil {
		panic(err)
	}
	maxVersion, err := semver.NewVersion(mw.MaxVersion)
	if err != nil {
		panic(err)
	}
	return func(writer rest.ResponseWriter, request *rest.Request) {
		version, err := semver.NewVersion(request.PathParam("version"))
		if err != nil {
			rest.Error(
				writer,
				"Invalid version: "+err.Error(),
				http.StatusBadRequest,
			)
			return
		}
		if version.LessThan(*minVersion) {
			rest.Error(
				writer,
				"Min supported version is "+minVersion.String(),
				http.StatusBadRequest,
			)
			return
		}
		if maxVersion.LessThan(*version) {
			rest.Error(
				writer,
				"Max supported version is "+maxVersion.String(),
				http.StatusBadRequest,
			)
			return
		}
		request.Env["VERSION"] = version
		handler(writer, request)
	}
}

var (
    server    = flag.String("s", "localhost:9090", "listen server address")
)

func main() {
	flag.Parse()

	// 版本控制
	svmw := SemVerMiddleware{
		MinVersion: "0.0.1",
		MaxVersion: "2.0.0"}

	api := rest.NewApi()
	// 状态接口
	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)

	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/status", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(statusMw.GetStatus())
		}),
		rest.Get("/#version/info", svmw.MiddlewareFunc(
			func(w rest.ResponseWriter, req *rest.Request) {
				version := req.Env["VERSION"].(*semver.Version)
				if version.Major >= 1 {
					PrivateMessageAPIV1.Info(w, req)
				} else {
					PrivateMessageAPIV0.Info(w, req)
				}
			},
		)),
		// Session管理
		rest.Get("/#version/session", PrivateMessageAPIV1.GetSession),
		rest.Post("/#version/session", PrivateMessageAPIV1.PostSession),
		rest.Put("/#version/session", PrivateMessageAPIV1.PutSession),
		rest.Delete("/#version/session", PrivateMessageAPIV1.DeleteSession),

		// 用户管理
		rest.Get("/#version/user", PrivateMessageAPIV1.GetUserInfo),
		rest.Post("/#version/user", PrivateMessageAPIV1.Register),
		rest.Put("/#version/user", PrivateMessageAPIV1.ModifyUsername),
		rest.Delete("/#version/user", PrivateMessageAPIV1.DeleteUser),
		//rest.PUT("/#version/user/password", PrivateMessageAPIV1.ModifyPassword),

		// 联系人管理
		rest.Get("/#version/friend", PrivateMessageAPIV1.GetAllFriends),
		rest.Get("/#version/friend/:id", PrivateMessageAPIV1.GetFriend),
		rest.Post("/#version/friend", PrivateMessageAPIV1.AddFriend),
		//rest.Put("/#version/friend", PrivateMessageAPIV1.ModifyFriendNickname),
		rest.Delete("/#version/friend", PrivateMessageAPIV1.DeleteFriend),

		// 消息管理
		rest.Get("/#version/message/amount", PrivateMessageAPIV1.GetAllMessageCount),
		rest.Get("/#version/message/amount/:id", PrivateMessageAPIV1.GetMessageCount),
		rest.Get("/#version/message", PrivateMessageAPIV1.GetMessages),
		rest.Get("/#version/message/:id", PrivateMessageAPIV1.GetMessage),
		rest.Post("/#version/message", PrivateMessageAPIV1.SendMessage),
		rest.Delete("/#version/message", PrivateMessageAPIV1.DeleteMessage),
		rest.Put("/#version/message", PrivateMessageAPIV1.ReadMessage),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(*server, nil))
}

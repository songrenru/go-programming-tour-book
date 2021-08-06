package server

import (
	"encoding/json"
	"fmt"
	"github.com/songrenru/chatroom/global"
	"github.com/songrenru/chatroom/logic"
	"html/template"
	"net/http"
)

func homeHandleFunc(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Fprint(writer, "模板解析错误")
		return
	}

	err = tmpl.Execute(writer, nil)
	if err != nil {
		fmt.Fprint(writer, "模板执行错误")
		return
	}
}

func userListHandleFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	userList := logic.Broadcaster.GetUserList()
	b, err := json.Marshal(userList)

	if err != nil {
		fmt.Fprint(w, `[]`)
	} else {
		fmt.Fprint(w, string(b))
	}
}

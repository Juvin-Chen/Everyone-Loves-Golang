package handler

import (
	"fmt"
	"net/http"
)

// UserHandler
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// 解析 GET 参数
	query := r.URL.Query()
	name := query.Get("name")

	// 解析 POST 表单参数
	r.ParseForm()
	password := r.PostForm.Get("password")
	// 收到的 GET 参数，POST 参数
	fmt.Fprintf(w, "Received GET parameters: %s, POST parameters: %s", name, password)
}

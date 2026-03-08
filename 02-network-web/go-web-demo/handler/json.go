package handler

import (
	"encoding/json"
	"go-web-demo/model"
	"net/http"
)

func JSONResponseHandler(w http.ResponseWriter, r *http.Request) {
	user := model.User{ID: 1, Name: "Gopher", Age: 18}

	// 设置json响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 序列化并返回json
	json.NewEncoder(w).Encode(user)
}

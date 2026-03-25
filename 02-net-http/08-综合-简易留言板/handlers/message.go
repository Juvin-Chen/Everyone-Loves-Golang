/*
编写处理器（handlers）
包含三个处理器：
	1.IndexHandler：显示留言列表
	2.NewMessageFormHandler：显示发布表单
	3.CreateMessageHandler：处理表单提交（POST）
*/

package handlers

import (
	"html/template"
	"log"
	"message-board/middleware"
	"message-board/store"
	"net/http"
	"strconv"
	"strings"
)

// IndexHandler 显示所有留言
func IndexHandler(messageStore *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages := messageStore.GetAll()
		// 这是一个临时的数据盒子，专门用来把留言数据传给 HTML 模板
		data := struct {
			Messages []store.Message
		}{
			Messages: messages,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// NewMessageFormHandler 显示发布表单
func NewMessageFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateMessageHandler 处理 POST 提交
func CreateMessageHandler(store *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// 解析表单
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		nickname := strings.TrimSpace(r.FormValue("nickname"))
		content := strings.TrimSpace(r.FormValue("content"))
		if nickname == "" || content == "" {
			http.Error(w, "昵称和内容不能为空", http.StatusBadRequest)
			return
		}
		// 保存留言
		store.Add(nickname, content)
		// 获取请求 ID（演示 context 使用）
		reqID := r.Context().Value(middleware.RequestIDKey)
		log.Printf("[RequestID: %v] 新留言来自 %s", reqID, nickname)
		// 重定向到首页
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func DeleteMessageHandler(store *store.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "无效留言 ID", http.StatusBadRequest)
			return
		}

		if err := store.DeleteByID(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

package handler

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only 'Get' requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello, welcome to the world of Go Web! Your request path is %s.", r.URL.Path)
}

// HelloWithMiddleware 供中间件示例使用
func HelloWithMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello with Middleware"))
}

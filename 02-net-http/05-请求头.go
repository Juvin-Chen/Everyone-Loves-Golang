/*
lesson 5 请求头

1.请求头的定义与作用。
2.常见请求头字段及其含义。
3.在 Go 中如何读取所有请求头、获取特定头、处理多值头。
4.特殊头 Host 和 Content-Length 的访问方式。
5.实际应用：根据 User-Agent 返回不同内容、获取客户端 IP。
6.安全注意事项：请求头可以被伪造。
*/

/*
1. 什么是请求头？
当客户端（浏览器、curl 等）向服务器发送 HTTP 请求时，除了请求行（方法、路径、版本）和请求体（数据）外，还会发送一组 键值对 来描述请求的附加信息，这就是 请求头（Request Headers）。
你可以把它看作快递包裹上的 寄件人信息、配送要求、物品类型 等标签。服务器根据这些标签来决定如何处理请求。
请求头的格式是：Header-Name: Header-Value
例如：
	User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...
	Accept: text/html,application/xhtml+xml,...
	Content-Type: application/json

2. 常见请求头字段及其作用
字段名			作用										 示例
User-Agent		标识客户端类型（浏览器、爬虫、API 工具等）		Mozilla/5.0 ...
Accept			告诉服务器，客户端能接收的响应内容类型			text/html,application/json
Accept-Language	客户端偏好的语言							  zh-CN,zh;q=0.9
Content-Type	请求体的格式（仅当有请求体时）				   application/x-www-form-urlencoded
Content-Length	请求体的长度（字节数）						  348
Authorization	认证信息（如 Token、Basic Auth）			  Bearer eyJhbGci...
Cookie			携带客户端的 Cookie 信息			         session_id=abc123
Referer			请求来源页面（从哪个页面跳转来的）		       https://www.google.com/
Host			请求的目标主机名和端口（必须存在）		       localhost:8080
Origin			跨域请求时，标识请求来源的域名			       http://example.com
Cache-Control	缓存控制指令							     no-cache
*/

package main

import (
	"fmt"
	"net/http"
	"strings"
)

/*
在 net/http 包中，*http.Request 类型有一个 Header 字段，它是一个 http.Header 类型（本质是 map[string][]string）。
我们可以通过 r.Header 来获取所有请求头，或通过 r.Header.Get(key) 获取某个头的值。
*/

func testRequestHeader() {
	http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
		// 获取所有请求头
		headers := r.Header

		fmt.Fprintf(w, "请求头如下：\n")
		for key, values := range headers {
			// 每个 key 可能对应多个值（例如 Set-Cookie 会有多个，但请求头中不多见）
			fmt.Fprintf(w, "%s: %v\n", key, values)
		}
	})
}

/*
用浏览器访问 http://localhost:8080/headers，你会看到类似：
请求头如下：
User-Agent: [Mozilla/5.0 ...]
Accept: [text/html,application/xhtml+xml,application/xml;q=0.9,...]
Accept-Encoding: [gzip, deflate, br]
Accept-Language: [zh-CN,zh;q=0.9]
...
*/

// 获取单个请求头
func testSingleRequest() {
	http.HandleFunc("/user-agent", func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		fmt.Fprintf(w, "你的 User-Agent 是: %s", ua)
	})
}

/*
测试：
curl -H "User-Agent: MyCustomBot" http://localhost:8080/user-agent
输出：你的 User-Agent 是: MyCustomBot
*/

// 3.获取请求头

// 获取可能的多值请求头
/*
有些头（比如 Accept 或 Cookie）可能包含多个值，但通常用 Get 取第一个就够了。
如果需要所有值，可以用 r.Header.Values(key) 返回 []string。
*/
func testMoreRequest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		values := r.Header.Values("Accept")
		fmt.Fprintf(w, "Accept 的所有值: %v", values)
	})
}

// 特殊头：Host 和 Content-Length
/*
r.Host：返回请求的目标主机名（不包含端口），但 r.Header.Get("Host") 也是可行的。
实际上 r.Host 是单独的字段，因为 Host 头在 HTTP/1.1 中必须存在，Go 单独提供了它。
Content-Length：可以通过 r.ContentLength 直接获取，它是一个 int64，表示请求体的大小。

示例：
fmt.Fprintf(w, "Host: %s\n", r.Host)
fmt.Fprintf(w, "Content-Length: %d\n", r.ContentLength)
*/

/*
4. 请求头的特性与注意事项
4.1 不区分大小写
HTTP 头名称是大小写不敏感的，例如 User-Agent 和 user-agent 是一样的。
r.Header.Get("user-agent") 和 r.Header.Get("User-Agent") 返回相同的结果。

4.2 客户端可以伪造任何头
请求头是客户端发送的，不能完全信任。例如 User-Agent 可以被修改，Referer 可以被伪造。
所以在做安全相关判断时（如身份验证、IP 限制），必须谨慎。

4.3 某些头可能不存在
并非所有请求都包含某些头，例如 Referer 只有在从其他页面跳转时才有；Authorization 只有主动携带时才有。
获取时一定要处理空值情况。

4.4 请求头大小限制
默认情况下，请求头的大小是有限制的（通常 1MB 左右），如果客户端发送了超大的头，服务器可能会拒绝处理。
*/

// 实战示例：根据 User-Agent 返回不同内容
// 假设你想根据客户端类型返回不同的页面：移动设备返回简版，电脑返回完整版。
func testReturnDifferent() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if strings.Contains(ua, "Mobile") {
			fmt.Fprintf(w, "这是移动端页面（简易版）")
		} else {
			fmt.Fprintf(w, "这是桌面端页面（完整版）")
		}
	})
	http.ListenAndServe(":8080", nil)
}

// 获取客户端 IP 地址
/*
虽然 IP 地址不在请求头中，但常常通过请求头来获取真实 IP（尤其是经过代理时）。
常见做法是检查 X-Forwarded-For 或 X-Real-IP 头。
*/
func testGetClientIP(r *http.Request) string {
	// 优先取 X-Forwarded-For
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// 可能有多个 IP，取第一个
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	// 其次取 X-Real-IP
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}
	// 最后取 RemoteAddr（直连 IP）
	return r.RemoteAddr
	// 注意：这些头同样可以被伪造，所以只适用于信任代理的情况。
}

/*
设置请求头（在客户端）
我们主要学习的是 服务端接收请求头，但有时我们也会写客户端（比如调用外部 API），那时需要设置请求头。
虽然不在本课重点，但简单提一下：
client := &http.Client{}
req, _ := http.NewRequest("GET", "https://api.example.com", nil)
req.Header.Set("Authorization", "Bearer token123")
resp, _ := client.Do(req)
*/

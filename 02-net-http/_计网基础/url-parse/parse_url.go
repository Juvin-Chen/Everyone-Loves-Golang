// 解析url
// 在 Go 语言中，我们主要使用 net 和 net/url 标准库
// url 格式：协议://（主机IP）服务器域名:端口号/路径?参数名=参数值。
package main

import (
	"fmt"
	"net/url"
)

// 把一个完整的 URL 字符串（比如百度搜索的 URL）解析成 Go 语言能识别的结构化对象，然后提取出「协议、主机名、路径、查询参数」这些关键信息。
func main() {
	// 这里用的就是教科书版的完整的 url 用于解析，实际上大部分是省略了路径和参数的极简版：https://www.baidu.com
	rawUrl := "https://www.baidu.com/s?wd=go%E8%AF%AD%E8%A8%80&rn=10&tn=baidu"
	// 解析原始 URL 字符串，返回*url.URL对象
	parseUrl, err := url.Parse(rawUrl)
	if err != nil {
		// Go的内置函数，调用后立即停止当前程序的执行，打印错误信息 + 调用栈（方便定位哪里错了）；
		panic(err)
	}

	fmt.Println("协议（Protocol）:", parseUrl.Scheme)
	fmt.Println("主机名（Host）:", parseUrl.Host)
	fmt.Println("路径（Path）:", parseUrl.Path)

	// 获取具体的参数值，把 URL 中 ? 后面的查询参数（wd=go%E8%AF%AD%E8%A8%80&rn=10&tn=baidu）解析成 url.Values 类型（本质是键值对的 map）
	queryParams := parseUrl.Query()
	fmt.Println("wd的值：", queryParams.Get("wd")) // 上面的 url ?后面是 wd
}

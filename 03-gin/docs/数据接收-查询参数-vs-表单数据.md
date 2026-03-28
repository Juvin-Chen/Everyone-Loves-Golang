# Gin 数据接收：查询参数 (Query) vs 表单数据 (Form Data)

## 一、 核心一句话总结
**“长得一样，但灵魂和位置完全不同！”**

无论是拼在 URL 里的 `?name=张三&age=18`，还是塞在 POST 请求体里的 `name=张三&age=18`，它们的**文本排版格式**是一模一样的（学名叫 `x-www-form-urlencoded` 格式）。
正因为它们长得一样，所以在 Gin 框架里定义结构体时，**统一都使用 `` `form:"xxx"` `` 标签**。

## 二、 本质区别

| 对比项 | GET 查询参数 (Query) | POST 表单数据 (Form Body) |
| :--- | :--- | :--- |
| **存放位置** | 暴露在 **URL** 尾部 (`?key=val`) | 隐藏在 HTTP **请求体 (Body)** 内部 |
| **核心作用** | 附加信息：用于**筛选、分页、搜索** | 核心业务：用于**提交真正的实体数据** |
| **Gin 绑定方法** | `c.ShouldBindQuery()` (只去 URL 里找) | `c.ShouldBind()` (自动去 Body 里找) |

## 三、 实战代码对比（以开发博客后端为例）

在写实际的 Gin 接口时，这两种数据有着截然不同的应用场景：

### 场景 1：获取博客文章列表（GET）—— 使用查询参数
前端发来的请求：`GET /api/articles?page=1&limit=10`
数据在 URL 里，属于 Query 参数。

```go
// 1. 定义结构体（标签用 form）
type Pagination struct {
    Page  int `form:"page"`
    Limit int `form:"limit"`
}

// 2. Gin 路由处理
r.GET("/api/articles", func(c *gin.Context) {
    var p Pagination
    // 💡 重点：明确告诉 Gin 去 URL 里面把数据绑过来
    if err := c.ShouldBindQuery(&p); err != nil {
        c.JSON(400, gin.H{"error": "参数错误"})
        return
    }
    c.JSON(200, gin.H{"msg": "获取列表成功", "page": p.Page})
})
```

### 场景 2：发布一篇新博客（POST）—— 使用表单数据
前端发来的请求：`POST /api/articles`
数据在 Body 里：`title=Go语言学习笔记&author=admin`

```go
// 1. 定义结构体（标签同样用 form！）
type Article struct {
    Title  string `form:"title"`
    Author string `form:"author"`
}

// 2. Gin 路由处理
r.POST("/api/articles", func(c *gin.Context) {
    var article Article
    // 💡 重点：ShouldBind 会自动识别 POST 请求，并去 Body 里把表单数据绑过来
    if err := c.ShouldBind(&article); err != nil {
        c.JSON(400, gin.H{"error": "参数错误"})
        return
    }
    c.JSON(200, gin.H{"msg": "发布成功", "title": article.Title})
})
```

## 四、 防坑指南
* **永远不要**在 POST 请求的 Body 里发 `{"title":"Go"}`（这是 JSON！），却指望 Gin 用 `` `form` `` 标签去接收。如果要接 JSON，标签必须换成 `` `json` ``，方法必须换成 `c.ShouldBindJSON()`。
* **只要看到数据长得像 `key=val&key=val`**，结构体闭着眼睛打 `` `form` `` 标签准没错，剩下的就是用 GET 还是 POST 的问题了。

***


# 📖 前后端交互秒懂指南：表单数据 vs JSON

## 一、 核心概念：这俩到底是什么？

无论是「表单数据」还是「JSON 数据」，它们的目的只有一个：**前端给后端传参**。
区别仅仅在于**包装格式**和**使用场景**不同。

| 对比维度 | 表单数据 (Form Data) | JSON 数据 (JSON Body) |
| :--- | :--- | :--- |
| **通俗比喻** | 散装纸条，一张写一个参数 | 打包好的结构化快递盒 |
| **数据长相** | `name=张三&age=18` | `{"name":"张三", "age":18}` |
| **前端场景** | 传统 HTML 页面、带文件的上传 | 现代前后端分离接口 (Vue/React) |
| **能否传文件** |  **能** (必须用表单传文件) |  **不能** (只能传纯文本) |



---

## 二、 终极困惑：`name=张三&age=18` 到底在哪里？

很多后端新手看到 `name=张三&age=18` 这种格式，总搞不清它是在 URL 里，还是在请求体（Body）里。

**结论：它的位置完全取决于前端用的是 `GET` 还是 `POST` 请求！**

### 1. 如果前端用 `GET` 提交（数据在 URL 里）
* **长相**：拼在 URL 问号后面，例如 `http://127.0.0.1:8080/user?name=张三&age=18`
* **特点**：明文可见（不安全），**没有请求体**。
* **Gin 接收**：`c.Query("name")`

### 2. 如果前端用 `POST` 提交（数据在请求体里）
* **长相**：URL 干干净净 `http://127.0.0.1:8080/user`。数据 `name=张三&age=18` 被塞进了请求体 (Body) 中。
* **特点**：URL 看不见参数，这才是标准的“表单提交”。
* **Gin 接收**：`c.PostForm("name")`



> **重要推论：** JSON 数据格式 `{"name":"张三"}` 永远只存在于**请求体**中，绝对不会出现在 URL 里。

---

## 三、 前端怎么发？后端 (Gin) 怎么收？

前后端必须“对暗号”，前端发什么格式，后端就必须用对应的方法解析。

### 场景一：前端发「URL 查询参数」
* **前端代码**：直接请求 `axios.get('/user?name=张三')`
* **Gin 接收**：
    ```go
    // 方式1：单独拿
    name := c.Query("name") 
    
    // 方式2：绑定结构体 (注意标签是 form)
    type User struct {
        Name string `form:"name"` 
    }
    var u User
    c.ShouldBindQuery(&u) 
    ```

### 场景二：前端发「POST 表单数据」
* **前端代码**：
    ```html
    <form action="/user" method="POST">
      <input type="text" name="name" value="张三">
    </form>
    ```
    ```javascript
    // 或者 Vue/React 中使用 FormData
    const fd = new FormData();
    fd.append("name", "张三");
    axios.post("/user", fd);
    ```
* **Gin 接收**：
    ```go
    // 方式1：单独拿
    name := c.PostForm("name") 
    
    // 方式2：绑定结构体 (注意标签依然是 form)
    type User struct {
        Name string `form:"name"` 
    }
    var u User
    c.ShouldBind(&u) // ShouldBind 会自动处理表单
    ```

### 场景三：前端发「JSON 数据」（最常用）
* **前端代码**：
    ```javascript
    // 传一个 JS 对象，axios 会自动转成 JSON 字符串塞进 Body
    axios.post("/user", {
        name: "张三",
        age: 18
    });
    ```
* **Gin 接收**：
    ```go
    // 只能用结构体接收，注意标签是 json！
    type User struct {
        Name string `json:"name"` // 关键：标签变了
        Age  int    `json:"age"`
    }
    var u User
    c.ShouldBindJSON(&u) // 必须用 ShouldBindJSON
    ```

---

## 四、 核心速查表

| 前端发送方式 | 数据长相 / 位置 | Gin 结构体标签 | Gin 核心接收方法 |
| :--- | :--- | :--- | :--- |
| **GET 请求** | `?name=张三` (在 URL 里) | `` `form:"name"` `` | `c.Query()` 或 `c.ShouldBindQuery()` |
| **POST 表单** | `name=张三` (在 Body 里) | `` `form:"name"` `` | `c.PostForm()` 或 `c.ShouldBind()` |
| **POST JSON** | `{"name":"张三"}` (在 Body 里) | `` `json:"name"` `` | **只能用** `c.ShouldBindJSON()` |

***

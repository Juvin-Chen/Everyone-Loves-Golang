// 参数绑定与验证
/*
如何使用 ShouldBind 系列方法绑定不同来源的参数。
结构体标签（form、json、uri）指定字段映射。
验证规则（binding 标签）实现自动校验。
自定义验证器扩展规则。
错误处理与友好提示。
*/

/*
1. 为什么需要参数绑定？
在 net/http 中，我们手动获取参数：
    路径参数：自己解析 r.URL.Path
    查询参数：r.URL.Query().Get("key")
    表单参数：r.ParseForm() + r.FormValue("key")
    JSON：读取 r.Body 并 json.Unmarshal
这样写起来重复且容易出错。Gin 提供了 c.ShouldBind 系列方法，可以将请求中的参数自动绑定到结构体字段，并支持验证规则（如必填、长度范围等）。

2. 绑定不同来源的参数
Gin 支持绑定以下来源：
    路径参数（c.Param）
    查询参数（c.Query）
    表单数据（c.PostForm）
    JSON/XML/YAML 请求体
    头部信息
最常用的统一入口是 c.ShouldBind，它会根据 Content-Type 自动选择绑定方式。
*/

// 3. 基础用法：绑定查询参数
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type LoginRequest_4 struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func loginHandler(c *gin.Context) {
	var req LoginRequest_4
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"username": req.Username})
}

/*
form 标签：指定查询参数或表单参数的字段名。
binding:"required"：验证规则，表示该字段不能为空。
c.ShouldBindQuery 专门绑定查询参数（URL 中的 ?key=value）。
*/

// 4. 绑定表单数据（POST）
type RegisterForm struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

func registerHandler(c *gin.Context) {
	var form RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 处理注册逻辑
}

/*
c.ShouldBind 会根据请求的 Content-Type 自动选择绑定方式，对于 application/x-www-form-urlencoded 会绑定表单数据。
验证规则：email 验证邮箱格式，min=6 验证最小长度。
*/

// 5. 绑定 JSON 请求体
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"gte=0,lte=120"`
	Email string `json:"email" binding:"required,email"`
}

func createUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 处理创建用户逻辑
}

// json 标签：指定 JSON 字段名。
// binding:"gte=0,lte=120"：数值范围验证（大于等于0，小于等于120）。

// 6. 绑定路径参数
type UserParam struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func getUserHandler(c *gin.Context) {
	var param UserParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"user_id": param.ID})
}

/*
路由定义：
r.GET("/user/:id", getUserHandler)
uri 标签：绑定路径中的 :id 参数。

验证规则：min=1 确保 ID 至少为 1。
*/

// 7. 混合绑定
// 有时一个接口可能同时包含多种参数来源，例如路径参数 + 查询参数 + JSON 请求体。Gin 支持分别绑定：
type Path struct {
	UserID int `uri:"user_id" binding:"required"`
}
type Query struct {
	IncludeDetail bool `form:"detail"`
}
type Body struct {
	Name string `json:"name" binding:"required"`
}

func UpdateUserHandler(c *gin.Context) {
	var path Path
	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var query Query
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 使用 path.UserID, query.IncludeDetail, body.Name
}

// 8. 常用验证规则
/*
Gin 默认使用 go-playground/validator/v10 进行验证，支持的常用规则：

    required：字段不能为空（字符串、切片、map、指针等不能为零值）
    min、max：数值最小值/最大值，字符串最小/最大长度，切片最小/最大长度
    len：精确长度
    email：邮箱格式
    url：URL 格式
    oneof：枚举值，例如 oneof=admin user
    gt、gte、lt、lte：大于、大于等于、小于、小于等于
    numeric：数值字符串
    alpha：字母字符
    alphanum：字母数字

更多规则可以查看 validator 文档。
*/

// 9. 自定义验证器
// 如果内置规则不够用，可以注册自定义验证函数。

// 自定义验证：用户名不能包含 "admin"
func usernameNotAdmin(fl validator.FieldLevel) bool {
	return fl.Field().String() != "admin"
}

func testDIY() {
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notadmin", usernameNotAdmin)
	}

	// 给 Username 加上自定义判断
	type User struct {
		Username string `form:"username" binding:"required,notadmin"`
	}

	r.POST("/user", func(ctx *gin.Context) {
		var u User
		if err := ctx.ShouldBind(&u); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"msg": "ok"})
	})
	r.Run()
}

// 10. 错误处理
// ShouldBind 系列方法返回错误，我们可以解析错误信息，返回友好的提示。
// 错误类型是 validator.ValidationErrors，可以遍历并格式化。

// 用map数组是因为支持前端一次性返回多个错误
/*
关键名词解释：
    validator.ValidationErrors
    这是错误切片（数组），里面装着所有校验失败的字段错误，有几个错就存几个
    fieldErr.Field()：校验失败的字段名（比如 Username、Password）
    fieldErr.Tag()：校验失败的规则名（比如 required、notadmin、min）
*/
func handleValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag() // 可以映射为更友好的消息
		}
	}
	return errors
}

// 在处理器中：
func f(c *gin.Context) {
	/*
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"errors": handleValidationError(err)})
			return
		}
	*/
}

/*
11. 对比 net/http
操作	        net/http	                            Gin
获取查询参数	 r.URL.Query().Get("key")	             c.Query("key") 或绑定结构体 // 适合单个数据绑定，也可以c.ShouldBind()
获取表单参数	 r.ParseForm() + r.FormValue("key")	     c.PostForm("key") 或绑定结构体
获取 JSON	    json.NewDecoder(r.Body).Decode(&obj)	 c.ShouldBindJSON(&obj)
验证	        手动 if 判断	                          binding:"required,email" 自动验证
错误处理	    手动返回 400	                          统一返回验证错误
Gin 的参数绑定与验证让代码量减少 70% 以上，且更加清晰。
*/

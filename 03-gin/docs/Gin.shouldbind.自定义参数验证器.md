# Gin 自定义参数验证器

---

## 一、核心作用
Gin 默认的验证标签（`required`/`email`/`len` 等）不够用时，**自定义验证规则**，比如：
- 用户名不能包含 `admin`
- 密码强度校验
- 手机号格式自定义校验

你这段代码的功能：**注册一个自定义验证器，校验用户名不能为 `admin`**

---

## 二、你的代码逐行精讲
```go
// 自定义验证：用户名不能包含 "admin"
func usernameNotAdmin(fl validator.FieldLevel) bool {
	return fl.Field().String() != "admin"
}

func testDIY() {
	r := gin.Default()
	// 获取Gin内置的验证器引擎，并注册自定义规则
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notadmin", usernameNotAdmin)
	}
}
```

### 1. 自定义验证函数（固定写法）
```go
func usernameNotAdmin(fl validator.FieldLevel) bool {
	// 核心逻辑：返回 true=验证通过，false=验证失败
	return fl.Field().String() != "admin"
}
```
- **参数固定**：必须接收 `fl validator.FieldLevel`
- **返回值固定**：`bool` 类型
- `fl.Field().String()`：获取当前校验字段的**字符串值**
- 规则：用户名不等于 `admin` → 通过验证

### 2. 注册自定义验证器（固定语法）
```go
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 第一个参数：自定义标签名（结构体用）
	// 第二个参数：自定义验证函数
	v.RegisterValidation("notadmin", usernameNotAdmin)
}
```
- `binding.Validator`：Gin 内置的验证器
- `(*validator.Validate)`：类型断言，获取底层验证引擎
- `RegisterValidation("标签名", 验证函数)`：**核心注册语句**
  - 标签名：`notadmin`（结构体里写 `binding:"notadmin"`）
  - 验证函数：上面定义的 `usernameNotAdmin`

---

## 三、完整使用步骤
### 步骤1：定义自定义验证函数
### 步骤2：注册验证器（项目启动时注册一次即可）
### 步骤3：结构体绑定自定义标签
### 步骤4：接口中使用验证

---

## 四、完整可运行代码
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// ========= 1. 自定义验证函数 =========
// 用户名不能为 admin
func usernameNotAdmin(fl validator.FieldLevel) bool {
	// 获取字段值，转字符串后判断
	return fl.Field().String() != "admin"
}

// 定义用户结构体，绑定自定义验证标签
type User struct {
	Username string `json:"username" binding:"required,notadmin"` // 自定义标签：notadmin
	Password string `json:"password" binding:"required,min=6"`
}

func testDIY() {
	r := gin.Default()

	// ========= 2. 注册自定义验证器（固定语法） =========
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册：标签名 = notadmin，对应函数 = usernameNotAdmin
		v.RegisterValidation("notadmin", usernameNotAdmin)
	}

	// ========= 3. 使用验证的接口 =========
	r.POST("/user", func(c *gin.Context) {
		var user User
		// 绑定并自动验证
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "验证通过", "data": user})
	})

	r.Run(":8080")
}

func main() {
	testDIY()
}
```

---

## 五、关键固定语法（新手背会即可）
1. **自定义函数格式**
```go
func 自定义函数名(fl validator.FieldLevel) bool {
    // 验证逻辑 return true/false
}
```

2. **注册验证器格式**
```go
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("自定义标签名", 自定义函数名)
}
```

3. **结构体使用**
```go
字段名 类型 `binding:"自定义标签名"`
// 支持多验证：binding:"required,自定义标签"
```

---

## 六、测试效果
| 传入参数                | 结果       | 原因                 |
| ----------------------- | ---------- | -------------------- |
| `{"username":"test"}`   | 验证通过   | 不是 admin           |
| `{"username":"admin"}`  | 验证失败   | 触发自定义 notadmin 规则 |

---

## 七、总结
1. 这段代码是 **Gin 自定义验证器的固定注册语法**
2. 核心：`RegisterValidation("标签名", 验证函数)`
3. 结构体用 `binding:"标签名"` 触发验证
4. 验证函数返回 `true`=通过，`false`=失败
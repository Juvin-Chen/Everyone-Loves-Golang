# Go Web 依赖注入（闭包工厂+结构体接收器）


## 一、核心痛点与解决方案
### 1. 痛点：路由处理函数无法直接传参
Go 标准要求路由处理函数必须是 `func(w http.ResponseWriter, r *http.Request)`（net/http）或 `func(c *gin.Context)`（Gin），但业务中需要用到数据库、缓存等依赖，直接传参就会报错，因此需要用「闭包+工厂模式」解决。

### 2. 核心方案：闭包工厂模式
外层函数（工厂）接收依赖参数，返回符合框架要求的处理函数；内层函数（闭包）“记住”外层依赖，实现参数传递，这就是闭包的核心作用。

## 二、两种核心写法（闭包工厂 vs 结构体接收器）
### （一）闭包工厂模式（Gin 中最常用）
#### 1. 代码示例（Gin 框架适配版）
```go
// 1. 外层工厂函数（接收依赖，如数据库连接、缓存等）
func NewIndexHandler(store *store.MemoryStore) gin.HandlerFunc {
    // 外层函数可接收任意依赖（如数据库、缓存），启动时执行初始化
    return func(c *gin.Context) {
        // 闭包特性：直接使用外层的 store 依赖
        messages := store.GetAll() // 无需额外传参，直接调用
        c.JSON(200, gin.H{"messages": messages})
    }
}

// 2. 路由注册（核心：调用工厂函数，传入依赖）
r.GET("/index", NewIndexHandler(store)) // 调用工厂函数，触发初始化
```

#### 2. 核心逻辑
- 调用 `NewIndexHandler(store)` 时，外层工厂函数先执行（初始化依赖），返回闭包
- 闭包不立即执行，仅作为“处理模板”，等待请求触发
- 本质：用闭包“记住”依赖，实现无直接传参也能使用依赖

### （二）结构体接收器模式（多依赖场景首选）
#### 1. 代码示例（Gin 适配版）
```go
// 1. 定义结构体，存储所有依赖（数据库、缓存等）
type HandlerStruct struct {
    store *store.MemoryStore // 依赖参数
}

// 2. 初始化结构体（注入依赖）
func NewHandlerStruct(store *store.MemoryStore) *HandlerStruct {
    return &HandlerStruct{store: store}
}

// 3. 结构体方法（符合 Gin 处理函数签名）
func (h *HandlerStruct) Index(c *gin.Context) {
    // 直接使用结构体中的 store 依赖
    messages := h.store.GetAll()
    c.JSON(200, gin.H{"messages": messages})
}

// 4. 路由注册
r.GET("/index", handlerStruct.Index) // 直接使用结构体方法
```

## 三、关键区别（一眼分清）
| 方式 | 核心逻辑 | 适用场景 |
|------|----------|----------|
| 闭包工厂 | 外层函数传依赖，内层闭包“记住”依赖 | 依赖少、逻辑简单（如单数据库） |
| 结构体接收器 | 依赖存在结构体中，统一管理 | 多依赖（数据库+缓存+日志）、大型项目 |

## 四、总结
1.  闭包工厂模式：核心是「外层函数传依赖，内层闭包用依赖」，启动时执行外层函数，请求时执行闭包
2.  结构体接收器：核心是「把依赖存到结构体，通过方法调用使用」，更适合多依赖场景
3.  两者本质都是“间接传参”，只是实现方式不同，最终都能让处理函数用上依赖参数
## Go 项目架构与命名规范

### 一、Go 与传统 OOP 的核心思维差异

Go 并非“没有面向对象”，而是通过数据与行为解耦的方式，实现OOP的核心目标（封装、复用、多态）。本质差异在于Go强调简洁、显式和扁平化设计，避免传统OOP的继承复杂性。参考Go谚语：“A little copying is better than a little dependency”（少复制胜过少依赖）。

| 维度       | 传统 OOP（Java/C++）                         | Go                                                           |
| ---------- | -------------------------------------------- | ------------------------------------------------------------ |
| 核心单元   | 类（Class）：封装“数据 + 行为”，程序最小单元 | 包（Package）：以业务领域为核心，数据与行为分离              |
| 数据与行为 | 强绑定：属性和方法必须在同一个类中           | 松耦合：结构体（struct）存纯数据，方法独立定义，通过接收器（receiver）绑定 |
| 复用方式   | 继承（extends）：易导致类层级臃肿            | 组合（struct嵌入）+ 接口（鸭子类型）：扁平、无继承包袱       |
| 架构导向   | 围绕“类层级”设计，易过度抽象                 | 围绕“包”设计，追求极简、扁平、显式                           |

**提示**：Go的多态通过接口实现“鸭子类型”（如果它走路像鸭子，叫声像鸭子，那就是鸭子），无需显式继承。这让代码更灵活，但要求你显式定义依赖。

### 二、Go 项目的主流架构分层

Go 不推崇过度分层（参考Go谚语：“Clear is better than clever”），核心原则是：数据归数据、逻辑归逻辑、接口归接口。按业务领域拆包，而非按角色（Controller/Service/DAO）拆包。依赖方向：api → service → repo/domain → pkg（上层调用下层，避免循环依赖）。

#### 1. 中小项目通用分层（推荐入门）

适合Web服务或CLI工具。每个层职责明确，易测试。

| 分层     | 核心职责               | 包含内容                                      | 对应“数据/行为”思路                |
| -------- | ---------------------- | --------------------------------------------- | ---------------------------------- |
| domain/  | 纯数据定义，无任何逻辑 | 只定义结构体（struct），可选简单getter/setter | 纯数据载体，是行为的“操作对象”     |
| service/ | 核心业务逻辑           | 实现业务规则，调用repo/domain                 | 行为的核心载体（解耦后的业务行为） |
| api/     | 对外交互接口           | HTTP/gRPC/CLI handler，参数校验、响应封装     | 行为的“入口”，不做业务逻辑         |
| repo/    | 数据持久化             | 数据库/缓存操作，封装存储细节                 | 行为的“数据落地”层                 |
| pkg/     | 通用工具               | 日志、配置、加密等，与业务无关                | 通用行为                           |

**目录示例（用户模块）**：


```go
your_project/
├── domain/               # 纯数据结构体
│   └── user.go           # 避免Entity后缀，保持简洁
├── service/              # 业务逻辑
│   └── user_service.go
├── api/                  # HTTP 接口
│   └── http_handler.go   # 明确类型，避免模糊
├── repo/                 # 数据库操作
│   └── user_repo.go
├── pkg/                  # 通用工具
│   ├── logger/
│   └── config/
├── go.mod                # 依赖管理（go mod init）
├── main.go               # 程序入口，初始化依赖（如数据库连接）
└── user_test.go          # 测试文件（可选，放在对应包下）
```

**提示**：使用go mod管理依赖（如go get github.com/gin-gonic/gin）。每个包下添加_test.go文件进行单元测试。

#### 2. 复杂项目：按“业务领域”拆包（DDD思路）

当项目规模大时，按业务领域（如用户、订单、支付）拆包，每个领域内自含分层（类似微服务）。这符合DDD（Domain-Driven Design）的Go适配，避免单一巨包。


```go
your_project/
├── user/                 # 用户领域（独立包，可import）
│   ├── user.go           # 数据结构（domain）
│   ├── service.go        # 业务逻辑
│   ├── http_handler.go   # 接口
│   └── repo.go           # 存储
├── order/                # 订单领域
├── pay/                  # 支付领域
├── pkg/                  # 通用工具（跨领域共享）
├── go.mod
└── main.go
```

**提示**：领域间通信用接口定义契约（如type UserService interface { ... }），便于mock测试。

### 三、结构体与方法分离的命名规范

Go命名原则：简洁、语义化、大小写区分可见性（大写导出，小写私有）。参考Uber Go Style Guide：避免缩写，除非标准（如HTTP）。

#### 1. 包命名

- 规则：全小写、单数、按业务功能命名，避免拼音/无意义缩写。
- ✅ 正确：user、order、repo、service
- ❌ 错误：users、yonghuguanli、svc

#### 2. 结构体命名（domain层）

- 规则：PascalCase（大驼峰）、名词、单数，避免不必要后缀（如Entity，除非冲突）。

- 示例：


  ```go
  // domain/user.go
  package domain
  
  // User 用户核心数据（纯结构体，无逻辑）
  type User struct {
  	ID       int64  // 用户ID
  	Username string // 用户名
  	Password string // 密码哈希（使用bcrypt存储）
  	Age      int    // 年龄
  }
  
  // Admin 组合User的结构体（嵌入实现复用）
  type Admin struct {
  	User            // 匿名嵌入（组合）
  	Level    string // 管理员等级
  	RoleID   int64  // 角色ID
  }
  ```

**提示**：密码绝不存明文，使用golang.org/x/crypto/bcrypt哈希。

#### 3. 文件命名

- 规则：小写、下划线分隔，格式为业务_类型.go，一个文件只做一件事（单一职责）。
- ✅ 正确：user_service.go、http_handler.go、user_repo.go
- ❌ 错误：user.go（塞所有逻辑）、user_business.go（语义模糊）

#### 4. 方法/函数命名

（1）接收器方法（简单操作，绑定到结构体）

- 规则：PascalCase，动词/短语；接收器用1-2字母（如u）；只用于简单操作（如校验、格式化），避免复杂逻辑。

- 示例：


  ```go
  // domain/user.go （简单方法可放domain，避免跨包复杂性）
  package domain
  
  import "golang.org/x/crypto/bcrypt"
  
  // CheckPassword 校验密码（使用哈希比较）
  func (u *User) CheckPassword(inputPwd string) bool {
  	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPwd)) == nil
  }
  
  // FormatUsername 格式化用户名
  func (u *User) FormatUsername() string {
  	return "user_" + u.Username
  }
  ```

（2）业务函数（复杂逻辑，独立函数）

- 规则：PascalCase，按“动作 + 对象”命名（如LoginUser）；接收参数，返回结果 + error；优先用独立函数解耦。

- 示例：


  ```go
  // service/user_service.go
  package service
  
  import (
  	"errors"
  	"your_project/domain"
  	"your_project/repo"
  )
  
  // LoginUser 用户登录（校验+查库+更新）
  func LoginUser(username, password string) (*domain.User, error) {
  	user, err := repo.GetUserByUsername(username)
  	if err != nil {
  		return nil, err
  	}
  	if !user.CheckPassword(password) {
  		return nil, errors.New("密码错误")
  	}
  	if err := repo.UpdateUserLoginTime(user.ID); err != nil {
  		return nil, err
  	}
  	return user, nil
  }
  ```

（3）接口层函数（API/Handler）

- 规则：PascalCase，按“请求方式 + 路径”命名（如PostUserLogin）。

- 示例：


  ```go
  // api/http_handler.go
  package api
  
  import (
  	"github.com/gin-gonic/gin"
  	"your_project/service"
  )
  
  // PostUserLogin 处理 POST /user/login
  func PostUserLogin(c *gin.Context) {
  	var req struct {
  		Username string `json:"username" binding:"required"`
  		Password string `json:"password" binding:"required"`
  	}
  	if err := c.ShouldBindJSON(&req); err != nil {
  		c.JSON(400, gin.H{"error": "参数错误"})
  		return
  	}
  	user, err := service.LoginUser(req.Username, req.Password)
  	if err != nil {
  		c.JSON(400, gin.H{"error": err.Error()})
  		return
  	}
  	c.JSON(200, gin.H{"msg": "登录成功", "user": user})
  }
  ```

**提示**：使用binding:"required"增强校验。错误处理用errors包或fmt.Errorf。

### 四、核心原则总结

1. **架构分层**：优先按“业务领域”拆包，而非“角色”；分层只为隔离，拒绝过度抽象（无用接口/层级）。
2. **数据与行为**：结构体纯数据，行为用接收器方法（简单）或独立函数（复杂），保持解耦。
3. **命名规范**：包小写单数，结构体PascalCase，文件“业务_类型”，方法语义化（动作 + 对象）。
4. **极简主义**：能不用接口就不用（只在需要多态时用）；优先简单函数/结构体。测试覆盖率高，参考官方testing包。
5. **额外建议**：阅读Effective Go（golang.org/doc/effective_go），实践小项目如Todo API。使用VS Code + Go插件调试。
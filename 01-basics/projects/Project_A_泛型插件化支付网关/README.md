# 泛型插件化支付网关项目说明

------

## 项目介绍

本项目是一个基于 **Go 语言** 开发的**模拟企业级泛型插件化支付网关系统**，实现了用户登录、多支付方式选择、余额读取、泛型统一支付处理等完整支付流程。

系统采用**策略模式 + Go 泛型**设计，支持支付宝、微信支付等多种支付方式灵活扩展，支付逻辑与网关逻辑解耦，具备高可维护性与扩展性。

------

## 项目结构

```go
Project_泛型插件化支付网关/
├── configs/                  # 配置文件目录
│   ├── payment_balance.json  # 支付账户余额数据
│   └── users.json            # 用户信息数据
├── internal/                 # 内部业务逻辑包
│   ├── payment/              # 支付相关核心逻辑
│   │   ├── alipay.go         # 支付宝插件
│   │   ├── base_payment.go   # 基础支付结构体与通用方法
│   │   ├── gateway.go        # 支付网关（泛型实现）
│   │   └── wechat.go         # 微信插件
│   └── user/                 # 用户相关逻辑
│       └── user.go           # 用户结构体与操作方法
├── app.go                    # 主程序入口
├── go.mod                    # Go 模块依赖管理
└── README.md                 # 项目说明文档
```
 
------

## 核心功能

该程序模拟了一个小型的支付网关：

1. **用户登录**：根据 `configs/users.json` 中的姓名匹配支付宝 / 微信支付 ID。
2. **多支付方式**：当前支持支付宝、微信；采用插件化策略，新增支付方式只需实现接口即可。
3. **余额持久化**：付款前后读取 / 写入 `configs/payment_balance.json`，模拟账户余额。
4. **泛型网关**：`Gateway[T]` 使用泛型封装支付结果，可适配不同类型的返回数据。
5. **稳健交互**：合理处理输入错误、循环重试、清空缓冲区等。


------


## 核心设计思路（结构体 / 方法 / 接口）

### 1. 接口设计：基于策略模式的「支付插件标准化」

#### 核心接口：`PaymentStrategy`

```go
type PaymentStrategy interface {
    Pay(amount float64) (string, error)
}
```

##### 设计思路：

- **标准化插件能力**：所有支付方式（支付宝 / 微信 / 银联）必须实现 `Pay` 方法，保证网关能「无差别调用」任意支付插件，符合「策略模式」核心思想；
- **最小接口原则**：仅定义核心支付能力（`Pay`），不冗余其他方法，降低插件实现成本；
- **返回值标准化**：统一返回 `(string, error)`（交易号 / 错误信息），让网关能直接处理所有插件的返回结果。

### 2. 结构体设计：「基础层 + 插件层」的分层复用

#### 2.1 基础层：`BasePayment`

```go
type BasePayment struct {
    Paytype   string  // 支付类型（AliPay/WeChat）
    PaymentID string  // 绑定的支付账户ID（如AliPay1001）
    Balance   float64 // 内存缓存的余额（与JSON文件同步）
}
```

##### 设计思路：

- **字段抽象**：提取所有支付方式的通用字段（支付类型、账户 ID、余额），避免支付宝 / 微信插件重复定义；

- 方法挂载（复用核心逻辑）

  ```go
  // 通用余额读取方法
  func (b *BasePayment) GetBalance() (float64, error) {}
  // 通用余额更新方法
  func (b *BasePayment) UpdateBalance(amount float64) bool {}
  // 通用日志方法
  func (b *BasePayment) Log(msg string) {}
  ```

  - 所有通用操作（读余额、更余额、打日志）挂载到 `BasePayment`，插件层通过「嵌入结构体」直接复用，无需重复编码；
  - 符合 Go 「组合优于继承」的设计哲学，通过嵌入实现能力复用，而非传统继承。

#### 2.2 插件层：`AliPay`/`WeChat`（专注专属逻辑）

```go
// 支付宝插件：嵌入BasePayment复用通用能力
type AliPay struct {
    BasePayment
}

// 微信插件：嵌入BasePayment复用通用能力
type WeChat struct {
    BasePayment
}
```

##### 设计思路：

- **轻量插件**：插件结构体仅嵌入 `BasePayment`，不额外定义通用字段 / 方法，专注实现「专属支付逻辑」；

- 方法聚焦插件的 Pay方法只处理「该支付方式的特有逻辑」（如生成专属交易号、调用专属接口），通用逻辑（余额扣减、日志）直接调用 BasePayment 的方法：

  ```go
  func (a *AliPay) Pay(amount float64) (string, error) {
      // 通用逻辑：复用BasePayment的余额扣减
      if !a.UpdateBalance(amount) {
          return "", fmt.Errorf("余额不足")
      }
      // 专属逻辑：生成支付宝交易号
      tradeNo := fmt.Sprintf("ALI_%s_%d", a.PaymentID, amount*100)
      a.Log(fmt.Sprintf("支付宝支付成功：%s", tradeNo)) // 复用BasePayment的日志方法
      return tradeNo, nil
  }
  ```

#### 2.3 泛型层：`Result[T]`/`Gateway[T]`（适配任意结果类型）

```go
// 泛型支付结果：T适配不同支付方式的专属数据
type Result[T any] struct {
    Success bool    // 通用字段：支付是否成功
    PayType string  // 通用字段：支付方式
    Amount  float64 // 通用字段：支付金额
    Data    T       // 泛型字段：专属数据（交易号/结构体）
    Message string  // 通用字段：结果描述
}

// 泛型网关：T适配Result的泛型类型
type Gateway[T any] struct{}
```

##### 设计思路：

- 通用字段 + 泛型字段分离
  - 固定通用字段（Success / PayType / Amount / Message）：所有支付结果的公共属性，统一定义；
  - 泛型字段 `Data`：适配不同支付方式的「专属数据」（如支付宝返回 string 交易号、微信返回结构体订单信息），避免为每种结果定义单独的结构体；
- **网关无感知插件类型**：`Gateway[T]` 的 `ProcessPayment` 方法只依赖 `PaymentStrategy` 接口，不关心具体是支付宝还是微信插件，实现「插件无关性」；
- **类型安全**：通过泛型 `T` 约束返回结果类型，避免类型断言的滥用，同时保留灵活性（支持任意类型的专属数据）。

------

## 支付流程

```
启动系统 → 登录/退出 {判断}
├─ 退出 → 结束
└─ 登录 → 登录成功? {判断}
   ├─ 失败 → 回到登录/退出
   └─ 成功 → 选支付方式（支付宝/微信）→ 输入支付金额 → 读取JSON余额 → 泛型网关处理 → 校验金额→执行支付→更新余额 → 返回支付结果 → 结束
```

------

## 扩展说明

- **新增支付方式**：只需在 `internal/payment/` 新建结构体，实现 `Pay`，嵌入 `BasePayment`，并在 `app.go` 中创建实例加入策略切片。
- **新增用户/余额**：在相应 JSON 文件中追加记录。
- **结果类型**：若希望返回非字符串数据，可将 `Gateway` 的泛型参数改为对应类型

------

> lab4 的 Project 设计
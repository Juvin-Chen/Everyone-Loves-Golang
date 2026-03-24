/*
lesson 7
context 包详解

📦 核心概念：Context 到底是个啥？
你可以把 Context 想象成一个“随身携带的档案袋”。

在 Web 开发里，一个用户的请求进来，可能会经过十几道关卡（比如：先查日志 -> 再查有没有登录 -> 再去查数据库 -> 最后返回数据）。

如果每一个关卡都需要知道这个用户的 RequestID（请求流水号），你要怎么做？
笨办法： 把 RequestID 当作参数，在这十几道关卡的函数里层层传递。这会让代码变得极其臃肿。
Go 的神仙解法： 弄一个“档案袋（Context）”。在第一道关卡把流水号塞进档案袋里，然后只要把这个档案袋往下传就行了。谁需要看流水号，自己打开档案袋看一眼。

context.WithValue 就是用来往档案袋里塞东西的方法。
*/

/*
1. 为什么需要 context？
在 Web 开发中，一个请求可能涉及多个操作：读取数据库、调用外部 API、处理数据等。
如果这个请求超时了，或者客户端主动断开连接，我们希望 中止所有后续操作，避免浪费资源。
同时，我们可能需要在多个处理函数之间传递一些请求范围的数据，比如用户信息、请求 ID 等。

Go 的 context 包正是为了解决这些问题而设计的。它提供了一种机制：
	传递截止时间（deadline）和取消信号（cancel）
	在多个 goroutine 之间传递值（但仅限请求范围内的数据）

2. context 的基本概念
context.Context 是一个接口，定义如下（简化）：
type Context interface {
    Deadline() (deadline time.Time, ok bool)  // 返回截止时间
    Done() <-chan struct{}                    // 返回一个只读 channel，当 context 被取消或超时时关闭
    Err() error                                // 返回取消原因（context.Canceled 或 context.DeadlineExceeded）
    Value(key interface{}) interface{}         // 获取绑定的值
}

（1）根 context：所有 context 的 “祖宗”
context.Background()：最常用的根 context，用于 “确定的场景”（比如 HTTP 服务器、数据库操作）；
ctx := context.Background() 是什么？
// 1. 调用工厂函数 (context.Background)
// 2. 工厂吐出一个具体的、看不见的对象 (实现了 Context 接口的对象)
// 3. 把这个对象赋值给变量 ctx，变量类型被标记为 context.Context

context.TODO()：临时用的根 context（比如还没想好用什么 context，后续要替换）；
特点：空的、不能取消、没有超时、没有值，只能用来派生其他 context。

（2）context 的四大核心能力: 见上所示接口

（3）派生 context：给根 context 加 “能力”
根 context 是空的，我们通过 4 个函数给它加能力，返回 “子 context”（父子关系，父取消子也取消）：
WithCancel：加「手动取消」能力（比如你主动走了，喊 “别做了”）；
WithTimeout/WithDeadline：加「自动超时取消」能力（比如你等了 2 分钟，自动喊 “别做了”）；
WithValue：加「存数据」能力（比如存你的身份证号）。
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// 模块 1：手动取消（WithCancel）—— 解决 “主动停任务” 问题
// 场景：模拟一个耗时任务，手动触发取消，任务立刻停止。

// 模拟耗时任务：每隔1秒打印一次，直到context取消
func longTask(ctx context.Context) {
	i := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("任务被取消，原因：", ctx.Err())
			return // 立刻停止任务
		default:
			fmt.Printf("任务执行中...第%d秒\n", i)
			time.Sleep(1 * time.Second)
			i++
		}
	}
}

func testLongTask() {
	// 1. 创建根context，派生带取消能力的子context
	/*
		它的完整含义是：
			context.Background(): 首先，创建一个最顶层的、永不关闭的“总开关”（根上下文）。
			context.WithCancel(...): 然后，基于这个“总开关”，创建一个新的、可被控制的“任务分支”（派生一个子上下文）。
			ctx, cancel := ...: 这个操作会返回两个东西：
			ctx：就是这个新的“任务分支”本身。你可以把它传递给一个或多个协程（goroutine），让它们开始工作。
			cancel：就是这个“任务分支”的“停止按钮”（取消函数）。当你调用 cancel() 时，所有监听这个 ctx 的协程都会收到一个“停止”信号，然后优雅地退出。
	*/
	ctx, cancel := context.WithCancel(context.Background()) // Withcancel手动取消
	// 重要：最后一定要调用cancel释放资源（defer保证函数退出时执行）
	defer cancel()
	go longTask(ctx)
	// 3. 模拟“用户等了3秒不耐烦，手动取消”
	time.Sleep(3 * time.Second)
	fmt.Println("用户：别做了！")
	cancel() // 触发取消信号

	// 等一下，看任务是否停止
	time.Sleep(1 * time.Second)
	fmt.Println("程序结束")
}

// 模块 2：超时自动取消（WithTimeout）—— 解决 “超时停任务” 问题
// 场景：模拟数据库查询，设置 2 秒超时，超时自动取消。

// 模拟数据库查询（耗时3秒）
func queryDB(ctx context.Context, id int) (string, error) {
	// 模拟耗时操作
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(3 * time.Second):
		// 3秒后查询完成
		return fmt.Sprintf("用户%d的信息：姓名张三", id), nil
	}
}

func testDBquery() {
	// 1. 创建带2秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 必须调用，释放资源

	// 2. 执行查询
	fmt.Println("开始查询数据库...")
	if userInfo, err := queryDB(ctx, 1); err != nil {
		// 超时后，ctx.Done() 自动关闭，任务立刻停止，ctx.Err() 返回 “超时” 原因。
		fmt.Println("查询失败：", err) // 超时会打印 context deadline exceeded
	} else {
		fmt.Println("查询成功：", userInfo)
	}
}

/*
运行结果：
开始查询数据库...
查询失败： context deadline exceeded
*/

// 模块 3：传递请求级数据（WithValue）—— 解决 “层层传参数” 问题
// 场景：模拟中间件给请求加唯一 ID，业务处理器读取这个 ID。

/*
context.WithValue 函数接收三个参数
这三个参数分别是：
1.parent (父上下文)
在你的代码中是： context.Background()
它是什么： 这是新上下文的“父级”或“基础”。所有新的上下文都必须基于一个已有的上下文创建。context.Background() 是一个空的、永不取消的根上下文，通常作为整个调用链的起点。你可以把它想象成一个空的快递箱。
2.key
3.value
*/

// 第一步：定义自定义key（必须用自定义类型，避免和其他包冲突）
type contextKey string

const requestIDKey contextKey = "requestID" // 请求ID的key

// 模拟中间件：生成请求ID，存入context
func requestIDMiddleware(next func(ctx context.Context)) {
	// 生成唯一ID
	requestID := uuid.NewString()
	// 从根context派生，存入请求ID
	ctx := context.WithValue(context.Background(), requestIDKey, requestID)
	// 调用下一个处理函数，传入带ID的context
	next(ctx)
}

// 模拟业务处理器：读取context里的请求ID并打印
func businessHandler(ctx context.Context) {
	// 从context取请求ID（类型断言）
	requestID := ctx.Value(requestIDKey).(string)
	fmt.Println("处理请求，请求ID：", requestID)
}

func testWithValue() {
	// 这里的 businessHandler 被当成一个函数参数传递进去，就相当于 requestIDMiddleware 的 next
	requestIDMiddleware(businessHandler)
}

/*
运行结果（每次 ID 不同）：
处理请求，请求ID： f47ac10b-58cc-4372-a567-0e02b2c3d479

核心点：
	必须用自定义类型当 key（比如 contextKey string）：如果用普通 string，不同包可能用同一个 key，导致数据冲突；
	WithValue 只存 “请求级数据”（比如请求 ID、用户信息），不要存可选参数（比如数据库地址）；
	取值时要做类型断言（.(string)），因为 Value 返回的是interface{}。
*/

// 模块 4：HTTP 场景下的 context
// 场景：HTTP 请求中，监听客户端断开 / 超时，同时传递请求 ID。

// 自定义key
type contextKey_ string

const requestIDKey_ contextKey_ = "requestID"

// 中间件：生成请求ID，存入context
func requeatIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 生成请求ID
		requestID := uuid.NewString()
		ctx := context.WithValue(r.Context(), requestIDKey_, requestID)

		// r.WithContext(ctx) 做了一次“浅拷贝”（复制）。它把原来的 r 复制了一份，但是，把新复制出来的这份的口袋（Context），换成了你刚才做好的 ctx（装着 RequestID 的那个）。
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 处理器：模拟耗时任务，监听请求取消
func longTaskHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(requestIDKey_).(string)
	log.Printf("开始处理任务，请求ID：%s", requestID)
	ctx := r.Context()

	select {
	case <-time.After(10 * time.Second):
		// 任务完成
		fmt.Fprintf(w, "任务完成！请求ID：%s", requestID)
	case <-ctx.Done():
		// 请求被取消（客户端断开/超时）
		log.Printf("请求%s被取消，原因：%v", requestID, ctx.Err())
		w.WriteHeader(http.StatusRequestTimeout) // 408
		fmt.Fprintf(w, "请求超时/被取消！请求ID：%s", requestID)
		return
	}
}

func testContextHttp() {
	http.Handle("/longtask", requeatIDMiddleware(http.HandlerFunc(longTaskHandler)))
	log.Println("服务器启动：http://localhost:8080/longtask")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
核心点：
	r.Context()：每个 HTTP 请求自带的 context，客户端断开 / 服务器超时会自动取消；
	r.WithContext(ctx)：把派生的 context 赋值给请求，保证后续处理器能拿到；
	监听<-ctx.Done()：就能感知客户端是否断开连接。
*/

/*
核心总结：
	1.context 的核心用途：传递「取消信号 / 超时」+ 传递「请求级数据」；
	2.根 context：Background()（正式）/TODO()（临时）；
	3.派生 context：
		WithCancel：手动取消；
		WithTimeout：自动超时取消（必须调用 cancel）；
		WithValue：存请求级数据（自定义 key，只存轻量、不可变数据）；
	4.HTTP 场景：
		r.Context()：请求自带的 context，客户端断开会自动取消；
		r.WithContext(ctx)：更新请求的 context；
必避坑：
	忘记调用cancel：导致资源泄露；
	用WithValue存可选参数：函数签名不清晰；
	用内置类型（如 string）当WithValue的 key：容易冲突。
*/

/*
补充：关于资源泄露
Go context 忘记 cancel () → 资源泄露
ctx, cancel := context.WithTimeout(...)
// 忘记写 cancel()

现象：
内存没有被释放
有一个 goroutine 在后台一直等着 ctx.Done ()
它永远等不到，所以一直活着
占内存、占线程、占资源
→ 垃圾回收收不回去
*/

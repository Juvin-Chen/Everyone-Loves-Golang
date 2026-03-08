package main

import (
	"fmt"
	"sync"
	"time"
)

/*
这一章是 Go 的杀手锏。理解 Goroutine (轻量级线程) 和 Channel (通信管道)。
牢记口号：“不要通过共享内存来通信，要通过通信来共享内存”。
*/

func main() {
	demo1_goroutine()
}

/*------------------------------------------------------------------*/
/*
demo0 前置部分，
对并发的简介：并发是指程序同时执行多个任务的能力。
高并发（High Concurrency）是指：一个系统在同一时间段内，能处理大量的请求 / 任务，并且还能保持稳定、快速的响应。

Goroutines：
1.Go 中的并发执行单位，类似于轻量级的线程。
2.Goroutine 的调度由 Go 运行时管理，用户无需手动分配线程。
3.使用 go 关键字启动 Goroutine。
4.Goroutine 是非阻塞的，可以高效地运行成千上万个 Goroutine。

Channel：
1.Go 中用于在 Goroutine 之间通信的机制。
2.支持同步和数据共享，避免了显式的锁机制。
3.使用 chan 关键字创建，通过 <- 操作符发送和接收数据。

Scheduler（调度器）：
Go 的调度器基于 GMP 模型，调度器会将 Goroutine 分配到系统线程中执行，并通过 M 和 P 的配合高效管理并发。
1.G：Goroutine。
2.M：系统线程（Machine）。
3.P：逻辑处理器（Processor）。

goroutine
语法格式：go 函数名( 参数列表 )
例如：go f(x, y, z)
*/

/*------------------------------------------------------------------*/
/*
demo1 Goroutine (协程)
Go 运行时管理的轻量级线程。
Goroutine（简称 Goroutine / 协程）是 Go 语言特有的轻量级 “线程” —— 由 Go 运行时（runtime）而非操作系统直接管理
可以把它理解成：Java 线程是 “重量级员工”（占资源多、创建慢），Goroutine 是 “轻量级临时工”（占资源极少、创建超快）
*/

func say(s string) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func demo1_goroutine() {
	// 开启一个新的协程执行 say("world")
	go say("world")

	// 当前主协程继续执行 say("hello")
	say("hello")

	// 注意：如果主协程(main)结束了，所有的子协程会被强制终止。
	// 实际开发中会用 WaitGroup 来等待，但阶段一先理解概念。
}

/*
知识补充：
Java 多线程：默认共享内存，需要加sync/Lock锁（容易死锁）；
比如两个线程改同一个变量，必须加synchronized，否则数据乱了；
Go 协程：Go 不推荐 “共享内存 + 加锁”，而是用Channel（管道）通信 —— 一个协程把数据放进管道，另一个从管道取，天然避免死锁（这就是 “不要通过共享内存通信，要通过通信共享内存” 的意思）。
*/

/*------------------------------------------------------------------*/
/*
demo2 Channel (通道)
通道可用于两个 goroutine 之间通过传递一个指定类型的值来同步运行和通讯。
使用 make 函数创建一个 channel，使用 <- 操作符发送和接收数据。如果未指定方向，则为双向通道。
ch <- v    // 把 v 发送到通道 ch
v := <-ch  // 从 ch 接收数据
           // 并把值赋给 v
声明一个通道很简单，我们使用chan关键字即可，通道在使用前必须先创建：
ch := make(chan int)

var ch chan string // ch是nil通道
如果只写上面一行而不写make，那就仅仅是声明，而没写创建，就是nil

基本规则的理解：
Channel —— “管道与接力棒” (最难理解的符号 <-)
这是最容易晕的地方。Go 发明了一个箭头 <-，其实它非常形象，它就是数据的流向。
想象一根水管（Channel）：
1.塞东西进去 (发送)：
c <- sum  // 把 sum 这个球，塞进 c 这根管子里
动作：箭头指向 c，说明数据流进了管子。

2.拿东西出来 (接收)：
x := <-c  // 从 c 这根管子里，把球拿出来，给 x
动作：箭头从 c 出来，说明数据流出了管子。

核心机制：阻塞（Blocking）—— 也就是“死等” 这跟 C++ 的队列不一样。Go 的无缓冲 Channel 是没有仓库的。
1.发送者：如果没人来拿，我就死死地卡在这里，手里拿着数据不敢松手（代码卡在 c <- sum 不动）。
2.接收者：如果管子里没东西，我就死死地卡在这里等（代码卡在 <-c 不动）。

!!! 接收者不一定是主协程，主协程只是最常见的接收者之一—— 通道是协程间的 “通用管道”，任何协程都能当发送方 / 接收方，主协程只是这段示例里刚好扮演了接收者的角色。

图解辅助：
(这张图会展示两个小人，一个把数据放进管子，另一个必须在另一头接住，否则第一个人就不能走。这就是“同步”。)
*/

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // 将 sum 发送到通道 c
}

func demo2_channel() {
	s := []int{7, 2, 8, -9, 4, 0}

	// 创建一个 int 类型的 channel
	c := make(chan int)

	// 分两个协程计算数组前半部分和后半部分
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)

	// 从通道接收数据 (会阻塞直到收到数据)
	x, y := <-c, <-c

	fmt.Println(x, y, x+y)

	// Buffered Channel (缓冲通道)
	// 仅当通道缓冲区满时发送才阻塞，缓冲区空时接收才阻塞
	ch_buf := make(chan int, 2)
	ch_buf <- 1
	ch_buf <- 2
	// ch_buf <- 3 // 如果打开这行会死锁(deadlock)，因为缓冲区只有2，且没有接收者
	fmt.Println(<-ch_buf)
	fmt.Println(<-ch_buf)
}

/*------------------------------------------------------------------*/
/*
demo3 通道缓冲区 & Go 遍历通道与关闭通道
1. 通道缓冲区
通道可以设置缓冲区，通过 make 的第二个参数指定缓冲区大小：
ch := make(chan int, 100)

带缓冲区的通道允许发送端的数据发送和接收端的数据获取处于异步状态，就是说发送端发送的数据可以放在缓冲区里面，可以等待接收端去获取数据，而不是立刻需要接收端去获取数据。
不过由于缓冲区的大小是有限的，所以还是必须有接收端来接收数据的，否则缓冲区一满，数据发送端就无法再发送数据了。

注意：如果通道不带缓冲，发送方会阻塞直到接收方从通道中接收了值。如果通道带缓冲，发送方则会阻塞直到发送的值被拷贝到缓冲区内；如果缓冲区已满，则意味着需要等待直到某个接收方获取到一个值。接收方在有值可以接收之前会一直阻塞。

2. 遍历通道与关闭通道
Go 通过 range 关键字来实现遍历读取到的数据，类似于与数组或切片。格式如下：
v, ok := <-ch
如果通道接收不到数据后 ok 就为 false，这时通道就可以使用 close() 函数来关闭。

!!!
接收者不一定是主协程，主协程只是最常见的接收者之一—— 通道是协程间的 “通用管道”，任何协程都能当发送方 / 接收方，主协程只是这段示例里刚好扮演了接收者的角色。
*/

// 1. 通道缓冲区
func demo3_channel_buffer() {
	// 这里我们定义了一个可以存储整数类型的带缓冲通道
	// 缓冲区大小为2
	ch := make(chan int, 2)

	// 因为 ch 是带缓冲的通道，我们可以同时发送两个数据
	// 而不用立刻需要去同步读取数据
	ch <- 1
	ch <- 2

	// 获取这两个数据
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

// 2. 遍历通道与关闭通道
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}
func demo3_channel_close() {
	// 1. 创建带缓冲的通道c，容量10（最多存10个int）
	c := make(chan int, 10)
	// 2. 开子协程执行fibonacci函数：生成10个斐波那契数（cap(c)=10），发送到c
	go fibonacci(cap(c), c)
	// 3. 主协程遍历通道c，消费里面的数据
	// range会不断从c里取数据，直到通道关闭且数据取完
	for i := range c {
		fmt.Println(i) // 打印从通道里拿到的斐波那契数
	}
	/*
		for i := range c 是 v, ok := <-ch 的简化写法，等价于：
		// 等价的手动遍历写法（帮你理解range的底层）
		for {
			i, ok := <-c // 从通道取数据，ok=true表示有数据，false表示通道关闭且无数据
			if !ok {
				break // 通道关闭，退出循环
			}
			fmt.Println(i)
		}
		range 自动帮你做了 “循环取数据 + 判断通道是否关闭” 的逻辑，更简洁。
	*/
}

/*------------------------------------------------------------------*/
/*
demo4 Select(多路复用)
select 语句使一个 Go 程可以等待多个通信操作，类似于 Switch，但用于 Channel。

大白话翻译： 你现在有两个电话（两个 Channel），你只有一只手。 select 就是让你同时盯着这两个电话：
1.如果电话 A 响了，你就接 A。
2.如果电话 B 响了，你就接 B。
3.如果两个都响了？Go 会随机挑一个接。
这就叫“多路复用”，专门用来处理“不知道哪个管道先来数据”的情况。
*/

func demo4_1_select() {
	// 1. 创建两个无缓冲通道（空快递盒）
	c1 := make(chan string)
	c2 := make(chan string)

	// 2. 启动子协程1：模拟“1秒后给c1放快递（数据）”
	go func() {
		time.Sleep(1 * time.Second) // 等1秒（模拟干活耗时）
		c1 <- "one"                 // 给c1通道发数据“one”（放快递）
	}()

	// 3. 启动子协程2：模拟“2秒后给c2放快递（数据）”
	go func() {
		time.Sleep(2 * time.Second) // 等2秒
		c2 <- "two"                 // 给c2通道发数据“two”
	}()

	// 4. 循环2次：因为要收2个快递（c1和c2各一个）
	for i := 0; i < 2; i++ {
		// 核心：select同时监听c1和c2
		select {
		// 情况1：c1有数据（快递到了），就取出来打印
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		// 情况2：c2有数据（快递到了），就取出来打印
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}

// demo4_2_select()是 “生产 - 消费” 模型：一个协程生产斐波那契数（发数据），另一个协程消费（收数据），用select同时处理 “生产数据” 和 “退出信号”。
func fibonacci_2(c, quit chan int) {
	x, y := 0, 1
	for { // 无限循环：一直生产，直到收到退出信号
		select {
		// case1：如果c通道能发数据（消费者能收），就发x
		case c <- x:
			x, y = y, x+y // 更新斐波那契数
		// case2：如果quit通道有数据（退出信号），就退出循环
		case <-quit:
			fmt.Println("quit")
			return // 退出函数，停止生产
		}
	}
}
func demo4_2_select() {
	// 1. 创建两个无缓冲通道
	c := make(chan int)    // 传斐波那契数的通道
	quit := make(chan int) // 传退出信号的通道

	// 2. 启动子协程（消费者）：收10个数，然后发退出信号
	go func() {
		for i := 0; i < 10; i++ { // 收10个数
			fmt.Println(<-c) // 从c通道取数（消费），打印
		}
		quit <- 0 // 收够10个，给quit通道发0（退出信号）
	}()

	// 3. 主协程执行fibonacci_2（生产者）
	fibonacci_2(c, quit)
	/*以上代码中，fibonacci goroutine 在 channel c 上发送斐波那契数列，当接收到 quit channel 的信号时退出。*/
}

/*
总结：
select 的本质：同时监听多个通道，哪个通道有数据（或能发数据）就执行哪个 case；
1. demo4_1 核心：监听多个 “读通道”，谁先到数据就处理谁；
2. demo4_2 核心：同时监听 “写通道（发数据）” 和 “读通道（退出信号）”，实现协程的可控退出；
新手记住 2 个规则：
select 里的 case 都是针对通道的（要么读、要么写）；
select 会阻塞，直到有一个 case 能执行；如果多个 case 能执行，随机选一个。
*/

/*------------------------------------------------------------------*/
/*
demo5 WaitGroup
sync.WaitGroup 用于等待多个 Goroutine 完成。
同步多个 Goroutine：
*/

// 下例代码用一个班长等待全班同学收拾书包的例子解释
// 其实核心就是：让主协程（班长）等所有子协程（同学）都干完活，自己再走 —— 不会出现 “班长先走了，同学还在教室” 的情况。

// worker函数：代表“同学收拾书包”的动作
// id：同学的编号；wg：班长的点名册（必须传指针，不然每个同学拿的是复印件！）
func worker(id int, wg *sync.WaitGroup) {
	// 关键1：defer wg.Done() → 同学不管有没有收拾完，只要离开教室，就告诉班长“我走了”
	// defer的作用：函数执行完（不管是正常结束还是出错），都会执行这行
	defer wg.Done()
	fmt.Printf("Worker %d started\n", id)  // 同学说：“我开始收拾书包了”
	fmt.Printf("Worker %d finished\n", id) // 同学说：“我收拾完了，要走了”
}

func demo5_waitgroup() {
	var wg sync.WaitGroup // 班长拿出空白点名册（创建WaitGroup）

	// 循环3次：对应等3个同学
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 关键2：点名册上+1 → “要等第i个同学”
		// 启动子协程：让第i个同学去收拾书包（异步执行worker）
		go worker(i, &wg) // &wg：把点名册的“原件”给同学（不是复印件！）
	}

	wg.Wait()                       // 关键3：班长站在门口等 → 直到点名册上的人数减到0
	fmt.Println("All workers done") // 所有同学都走了，班长锁门
}

/*------------------------------------------------------------------*/
/*
demo6 并发编程章节总结 ⭐

一、高级特性
1.Buffered Channel：
创建有缓冲的 Channel。
ch := make(chan int, 2)


2.Context：
用于控制 Goroutine 的生命周期。
context.WithCancel、context.WithTimeout。

（上下文）——“协程遥控器”：给协程配的 “遥控器”，能远程控制协程 “停止干活”；
生活场景：
你启动了 10 个协程爬网页，爬了一半你不想爬了 —— 用 Context 喊一声 “都停”，所有协程就乖乖退出；
没有 Context 的话，你只能等协程自己跑完，或者硬杀（容易出问题）；
核心用法（新手记个大概）：
// 1. 创建带“取消功能”的遥控器
ctx, cancel := context.WithCancel(context.Background())
// 2. 启动协程时把遥控器传进去
go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done(): // 收到“停止”信号
            fmt.Println("协程停止干活")
            return
        default:
            fmt.Println("协程继续干活")
        }
    }
}(ctx)
// 3. 不想干了，按遥控器
cancel()

新手要记：Context 就是 “控制协程生命周期” 的工具，比如设置超时（10 秒没干完就停）、手动取消，不用现在写，知道有这东西就行。


3.Mutex 和 RWMutex：
sync.Mutex 提供互斥锁，用于保护共享资源。
var mu sync.Mutex
mu.Lock()
// critical section
mu.Unlock()

翻译：Mutex=“互斥锁”，就是 “一个位置一个人用”；RWMutex=“读写锁”，更灵活的排队规则；
生活场景：
1. 多个协程抢着改同一个变量（比如共享雨伞的库存），就像多个人抢着用同一个厨房 —— 不加锁的话，有人刚进去一半，另一个人就推门进来，全乱了；
Mutex 就是 “厨房门锁”：一个人进去锁门（mu.Lock()），用完开门（mu.Unlock()），其他人必须等；
核心用法（新手能看懂的例子）：
var mu sync.Mutex // 买一把锁
var count int     // 大家抢着改的变量

func add() {
    mu.Lock()         // 上锁：我先用，别人等
    count++           // 改变量（抢东西）
    mu.Unlock()       // 解锁：用完了，下一个来
}
新手要记：锁是解决 “多个协程改同一个变量不乱” 的工具，和 Channel 二选一（Go 推荐用 Channel，但锁更简单直接）。


二、并发编程小结
Go 语言通过 Goroutine 和 Channel 提供了强大的并发支持，简化了传统线程模型的复杂性。配合调度器和同步工具，可以轻松实现高性能并发程序。
- Goroutines 是轻量级线程，使用 go 关键字启动。
- Channels 用于 goroutines 之间的通信。
- Select 语句 用于等待多个 channel 操作。

常见问题
1.死锁 (Deadlock)：
示例：所有 Goroutine 都在等待，但没有任何数据可用。
解决：避免无限等待、正确关闭通道。
“互相等对方，都卡死”
func main() {
    c := make(chan int)
    c <- 1 // 往无缓冲通道发数据，但没人取 → 主协程卡死（死锁）
}
避坑点：
别让协程 “无限等一个永远不会来的数据”；
通道用完记得关（但不是所有情况都要关，比如用 WaitGroup 的话不用）；

2.数据竞争 (Data Race)：
示例：多个 Goroutine 同时访问同一变量。
解决：使用 Mutex 或 Channel 同步访问。

翻译：多个协程同时改同一个变量，数据被改乱；
例子：两个协程同时给count++，本来该加 2，结果只加了 1；
解决办法：
用 Mutex 锁（排队改）；
用 Channel（一个协程改，其他协程通过通道传指令改）；

解释：
这些 “高级特性” 本质都是为了解决一个问题
并发编程的核心痛点：多个协程同时跑，容易 “乱套”（比如抢东西、卡死、关不掉），这些高级特性就是 Go 给你的 “工具”，专门解决这些乱套问题
*/

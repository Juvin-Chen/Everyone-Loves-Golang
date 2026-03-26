# Go 并发控制学习笔记（锁 vs Channel）

> 目标：讲清楚三个问题。
> 1) 为什么会有并发问题。
> 2) Mutex、RWMutex、Channel 分别解决什么问题。
> 3) 在真实项目里怎么选，怎么避免常见坑。

## 一、先理解“并发问题”到底是什么

并发 bug 的本质：多个 goroutine 在同一时间改同一份数据，导致结果不可预测。

最经典示例：

```go
var count int

func increment() {
    count++
}
```

`count++` 不是一个原子动作，它大致会拆成：

1. 读取 count 当前值。
2. 计算 `count + 1`。
3. 把新值写回去。

如果两个 goroutine 同时做这三步，就可能出现“加了两次但只生效一次”。

## 二、Mutex：最直接的互斥方案

Mutex 的意思是“同一时刻只允许一个 goroutine 进入临界区”。

```go
var (
    count int
    mu    sync.Mutex
)

func safeIncrement() {
    mu.Lock()
    defer mu.Unlock()

    count++
}
```

### 什么时候用 Mutex

- 共享状态很多是写操作。
- 逻辑简单，追求直观和稳定。
- 对读并发性能要求不高。

### 重点规则

1. Lock 和 Unlock 必须成对出现。
2. 推荐 `defer mu.Unlock()`，减少遗漏。
3. 临界区尽量小，别把耗时操作放锁里。

## 三、RWMutex：读多写少时更合适

RWMutex 把“读”和“写”分开处理：

- 读和读可以并发。
- 写和读互斥。
- 写和写互斥。

### 先回答一个高频问题

问题：读锁和写锁是不是同一把锁？

答案：是同一个 `sync.RWMutex` 对象，但内部区分“读计数状态”和“写状态”。

你可以这样理解规则：

1. 已有读锁时，新的读锁通常还能进来。
2. 只要还有任何读锁没释放，写锁就拿不到。
3. 写锁持有期间，新的读锁和写锁都要等待。

结论：你说的“读锁释放后才能获得写锁”是正确的。

```go
type MessageStore struct {
    messages []string
    mu       sync.RWMutex
}

func (s *MessageStore) GetAll() []string {
    s.mu.RLock()
    defer s.mu.RUnlock()

    out := make([]string, len(s.messages))
    copy(out, s.messages)
    return out
}

func (s *MessageStore) Add(msg string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.messages = append(s.messages, msg)
}
```

### 为什么这个写法是对的（时间线）

假设 A 在读，B 在写：

1. A 进入 `GetAll`，执行 `RLock`，开始复制切片。
2. B 进入 `Add`，执行 `Lock`，此时会阻塞等待。
3. A 复制完成并 `RUnlock`。
4. B 这时才能拿到写锁，执行追加。

所以 A 这次看到的是“读取当时的快照”，B 写入的是“后续状态”。
这就是并发下非常常见、也非常合理的结果。

### 什么时候用 RWMutex

- 读操作远多于写操作。
- 比如配置缓存、留言列表、热点读场景。

### 一个容易误解的点

RWMutex 不一定永远比 Mutex 快。

- 读多：RWMutex 常常更快。
- 写多：RWMutex 的管理开销可能不划算。

## 四、Channel：用于“传递数据和协作”

Channel 不是“锁的高级替代品”，它更像是 goroutine 之间的传送带。

```go
ch := make(chan int)

go func() {
    ch <- 42
}()

v := <-ch
fmt.Println(v)
```

### 无缓冲 vs 有缓冲

- 无缓冲：发送和接收要“对上”，常用于同步。
- 有缓冲：可先放入缓冲区，常用于削峰。

```go
ch1 := make(chan int)      // 无缓冲
ch2 := make(chan int, 10)  // 有缓冲
```

### 什么时候优先 Channel

- 任务分发（worker pool）。
- 生产者-消费者。
- 流水线处理。
- 信号通知（done、stop）。

## 五、锁和 Channel 怎么选（实战版）

先问自己：你在做“共享状态保护”还是“任务/数据流转”？

- 共享状态保护：优先锁。
- 任务流转和协作：优先 Channel。

再细分：

- 读多写少：RWMutex。
- 写多或结构简单：Mutex。

一句话记忆：

- 锁保护“同一份数据”。
- Channel 组织“多个 goroutine 的协作流程”。

## 六、结合留言板项目怎么落地

你的 `MemoryStore` 是典型共享状态，适合 RWMutex：

- `GetAll` 用 `RLock`。
- `Add/Delete` 用 `Lock`。
- 返回切片副本，避免外部绕过锁修改内部数据。

示例：

```go
type MemoryStore struct {
    messages []Message
    nextID   int
    mu       sync.RWMutex
}

func (s *MemoryStore) GetAll() []Message {
    s.mu.RLock()
    defer s.mu.RUnlock()

    result := make([]Message, len(s.messages))
    copy(result, s.messages)
    return result
}

func (s *MemoryStore) Add(nickname, content string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // 写入逻辑
}
```

## 七、常见坑（最容易踩）

### 1) 忘记解锁

```go
mu.Lock()
// 中途 return 或 panic
mu.Unlock() // 可能执行不到
```

改法：

```go
mu.Lock()
defer mu.Unlock()
```

### 2) 拷贝了带锁结构体

`sync.Mutex` / `sync.RWMutex` 不应该被复制。

建议：

- 结构体方法用指针接收者。
- 避免把含锁结构体按值传参。

### 3) 锁升级误用导致死锁

先 `RLock` 再直接 `Lock` 会出问题。

- 要先 `RUnlock`，再 `Lock`。

### 4) 关闭 Channel 的责任不清

原则：通常由发送方关闭通道，且只能关闭一次。

### 5) 向已关闭通道发送

会直接 panic，必须避免。

## 八、非常实用的并发模板

### WaitGroup 等待一组任务完成

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // do work
    }()
}

wg.Wait()
```

### select 处理超时

```go
select {
case v := <-ch:
    _ = v
case <-time.After(2 * time.Second):
    fmt.Println("timeout")
}
```

### 使用 race detector 查数据竞争

```bash
go test -race ./...
```

这是排查并发问题最有效的手段之一。

## 九、你可以这样记

1. 共享状态就加锁，任务协作用 Channel。
2. 读多写少用 RWMutex，写多用 Mutex。
3. 锁要成对，推荐 defer 解锁。
4. 带锁结构体不要复制，方法尽量用指针接收者。
5. Channel 一般发送方关闭，关闭只能一次。
6. 并发问题先上 `go test -race`。

## 十、最终总结

并发不是“同时跑很多 goroutine”这么简单，核心是“并发下数据是否一致、流程是否可控”。

- 如果你在维护一份共享状态，先把锁策略设计清楚。
- 如果你在设计多个 goroutine 的协作流程，先把 Channel 的方向、关闭时机和退出机制设计清楚。

把这两件事分开思考，你的并发代码就会稳很多。
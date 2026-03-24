/*
lesson 6 第二课时
演示中间件的链式调用
*/

/*
前置：关于 w r 参数的传递
最外层中间件（假如是中间件A）接受w r
	// 中间件 A
	func middlewareA(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // ① 接收 w/r
			fmt.Println("A:前置代码")
			next.ServeHTTP(w, r) // ② 把 w/r 传给 next（B）
			fmt.Println("A:后置代码")
		})
	}
然后在每个函数的内部 以这个形式进行层层传递 next.ServeHTTP(w, r) // ② 把 w/r 传给 next（B）

1. 问题本质：嵌套调用
假设我们有三个中间件：A、B，和一个最终处理器 Handler。

如果我们这样写：
final := A(B(Handler))
那么 final 实际上是一个新的 Handler，它里面包含了 A、B、Handler 的逻辑。

当请求进来时，执行顺序是：
	1.先执行 A 的“前置代码”
	2.再执行 B 的“前置代码”
	3.再执行 Handler 的业务逻辑
	4.然后执行 B 的“后置代码”
	5.最后执行 A 的“后置代码”
	也就是 先入后出（像洋葱一样，一层包一层）。

2. chain 辅助函数的常见实现
我们常见的 chain 函数可能是这样写的（从后往前循环）：
func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    // 从最后一个中间件开始包装，保证执行顺序与传入顺序一致
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

调用时：
handler := chain(finalHandler, A, B, C)

循环过程：
初始 handler = finalHandler
i = 2（最后一个中间件是 C）：handler = C(handler) → 现在 handler 是 C(finalHandler)
i = 1：handler = B(handler) → 现在 handler 是 B(C(finalHandler))
i = 0：handler = A(handler) → 现在 handler 是 A(B(C(finalHandler)))
最终结果和手动嵌套 A(B(C(finalHandler))) 完全一样。
因此，在这个实现下，chain 的传入顺序从左到右（A, B, C）对应手动嵌套的最外层到最内层，执行顺序也是 A → B → C（前置）。

3. 另一种可能的 chain 实现（从前往后循环）
有些开发者可能写成：
func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for _, mw := range middlewares {
        handler = mw(handler)
    }
    return handler
}
调用 chain(finalHandler, A, B, C)：
第一次：handler = A(finalHandler)
第二次：handler = B(A(finalHandler))
第三次：handler = C(B(A(finalHandler)))
最终结果是 C(B(A(finalHandler)))，执行顺序变为 C → B → A（前置），与传入顺序相反。
这种实现往往不符合直觉，所以大多数库（如 alice、negroni）都会采用从后往前循环，保证传入顺序就是执行顺序。
*/

package main

import (
	"fmt"
	"net/http"
)

// 一个直观的 Demo：用日志打印执行顺序

// 中间件 A
func middlewareA(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("A:前置代码")
		next.ServeHTTP(w, r)
		fmt.Println("A:后置代码")
	})
}

// 中间件 B
func middlewareB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("B: 前置代码")
		next.ServeHTTP(w, r)
		fmt.Println("B: 后置代码")
	})
}

// 最终业务处理器
func finalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("最终业务处理")
	fmt.Fprintf(w, "Hello, World!")
}

func testChain() {
	// 嵌套调用：先 B 包裹 finalHandler，再 A 包裹 (B包裹后的结果)
	handler := middlewareA(middlewareB(http.HandlerFunc(finalHandler)))

	http.Handle("/", handler)
	fmt.Println("服务器启动在 :8080")
	http.ListenAndServe(":8080", nil)
}

/*
运行这个程序，然后在浏览器访问 http://localhost:8080，控制台输出：
	A: 前置代码
	B: 前置代码
	最终业务处理
	B: 后置代码
	A: 后置代码
你看，先进入 A，然后进入 B，然后到业务，然后退回 B，再退回 A。
这就是 洋葱模型，也是中间件链的核心。
*/

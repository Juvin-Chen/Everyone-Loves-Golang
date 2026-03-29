# Go time 包核心：time.Time 与 time.Now () 

## 核心结论

1. `time.Time`：**时间数据类型**（Go 内置结构体类型，等同于 `int`/`string`），专门用于**存储、表示一个具体时间**
2. `time.Now()`：`time` 包**内置函数**，专门用于**获取当前系统时间**
3. `time.Now()` 调用后的**返回值**，就是 `time.Time` 类型的时间值

------

## 一、基础概念拆解

### 1. `time.Time`（类型）

- 归属：Go 标准库 `time` 包

- 本质：**时间类型**，用来定义时间变量、函数参数

- 通俗类比：

  - `int` → 存储整数
  - `string` → 存储文本
  - `time.Time` → 存储时间

  

### 2. `time.Now()`（函数）

- 归属：Go 标准库 `time` 包
- 本质：**工具函数**
- 作用：获取操作系统的当前时间
- 返回值：固定为 `time.Time` 类型

------

## 二、基础代码示例

```
package main

import (
	"fmt"
	"time" // 必须导入time包
)

func main() {
	// 调用time.Now()函数，获取当前时间
	// now 变量的类型 = time.Time
	now := time.Now()

	// 输出时间值
	fmt.Println("当前时间：", now)
	// 输出变量类型（验证：time.Time）
	fmt.Printf("变量类型：%T\n", now)
}
```

------

## 三、实战场景：自定义时间格式化函数

结合 Go Web 模板 / 业务开发的**最常用写法**：

### 1. 定义格式化函数

```
// 参数 t：类型为 time.Time，接收一个时间值
func formatDate(t time.Time) string {
	// 将时间格式化为 年-月-日
	return t.Format("2006-01-02")
}
```

### 2. 函数调用


```
// 传入 time.Now()，返回值刚好匹配 time.Time 类型
formatDate(time.Now())
```

------

## 四、核心区分（表格速记）






| 写法         | 身份     | 核心用途           | 代码示例          |
| :----------- | :------- | :----------------- | :---------------- |
| `time.Time`  | 数据类型 | 定义变量、函数参数 | `var t time.Time` |
| `time.Now()` | 函数     | 获取当前系统时间   | `t := time.Now()` |

------

## 五、新手避坑指南

### ❌ 错误写法



```
// 不存在 time.Time() 这个函数！
t := time.Time()
```

### ✅ 正确写法








```
// 1. 获取当前时间（最常用）
t := time.Now()

// 2. 声明一个时间类型的空变量
var t time.Time
```

------

## 六、极简总结

1. **`time.Time` = 装时间的盒子（类型）**
2. **`time.Now()` = 拿当前时间的工具（函数）**
3. 函数返回值 → 存入类型容器，两者一一对应
4. 时间格式化、传参时，必须用 `time.Time` 接收
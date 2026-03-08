package main

import (
	"fmt"
)

func main() {
	demo1()
}

/*--------------------------------------------------------------------------------------------------*/
/*
demo1 go语言数组
Go 语言提供了数组类型的数据结构。
数组是具有相同唯一类型的一组已编号且长度固定的数据项序列，这种类型可以是任意的原始类型例如整型、字符串或者自定义类型。

1.数组的声明: var arrayName [size]datatype
var balance [10]float32

2.在声明时，数组中的每个元素都会根据其数据类型进行默认初始化，对于整数类型，初始值为 0。
var numbers [5]int

3.还可以使用初始化列表来初始化数组的元素
var numbers_2 = [5]int{1, 2, 3, 4, 5}

4.:=的方式
numbers_3 := [5]int{1, 2, 3, 4, 5}

5.在 Go 语言中，数组的大小是类型的一部分，因此不同大小的数组是不兼容的，也就是说 [5]int 和 [10]int 是不同的类型。

- 以下定义了数组 balance 长度为 5 类型为 float32，并初始化数组的元素：
var balance = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}

- 我们也可以通过字面量在声明数组的同时快速初始化数组：
balance := [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}

- 如果数组长度不确定，可以使用 ... 代替数组的长度，编译器会根据元素个数自行推断数组的长度：
var balance = [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
或
balance := [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}

- 如果设置了数组的长度，我们还可以通过指定下标来初始化元素：
//  将索引为 1 和 3 的元素初始化
balance := [5]float32{1:2.0,3:7.0}
var balance = [5]float63{1:9.08, 3:5.0}

- 初始化数组中 {} 中的元素个数不能大于 [] 中的数字。
如果忽略 [] 中的数字不设置数组大小，Go 语言会根据元素的个数来设置数组的大小

- balance[4] = 50.0
以上实例读取了第五个元素。数组元素可以通过索引（位置）来读取（或者修改），索引从 0 开始，第一个元素索引为 0，第二个索引为 1，以此类推。

6.访问数组元素
数组元素可以通过索引（位置）来读取。格式为数组名后加中括号，中括号中为索引的值。例如：
var salary float32 = balance[9]
*/

func demo1() {
	// Element[%d]=%d 以此格式打印一个数组
	var n [10]int
	var i, j int
	for i = 0; i < 10; i++ {
		n[i] = i + 100
	}
	for j = 0; j < 10; j++ {
		fmt.Printf("Element[%d]=%d\n", j, n[j])
	}

	// 几个打印示例
	var i_, j_, k int
	balance := [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	for i_ = 0; i_ < 5; i_++ {
		fmt.Printf("balance[%d] = %f\n", i_, balance[i_])
	}

	balance2 := [...]float32{100.0, 2.0, 3.4, 7.0}
	for j_ = 0; j_ < 4; j_++ {
		fmt.Printf("balance2[%d] = %f\n", j_, balance2[j_])
	}

	balance3 := [5]float32{1: 2.0, 3: 4.0}
	for k = 0; k < 5; k++ {
		fmt.Printf("balance3[%d] = %f\n", k, balance3[k])
	}

	/* 补充部分：二维数组 & 向函数传递数组 */
	matrix := [2][3]int{
		{1, 2, 3}, // 第一行
		{4, 5, 6}, // 第二行
	}
	fmt.Println("开始遍历二维数组：")
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			fmt.Printf("matrix[%d][%d] = %d\n", i, j, matrix[i][j])
		}
	}

	// 可以让编译器自动推导行数，但不能推导列数
	animals := [...][2]string{
		{"Cat", "Dog"},
	}
	fmt.Println("\n推导后的数组animals:", animals)

	/*
		C++ 的习惯：在 C++ 中，当你把数组传给函数时，实际上传递的是首地址指针。所以在函数里修改数组，外面的数组也会变。
		Go 的特性：在 Go 中，数组是值类型 (Value Type)。
		当你把数组传给函数时，Go 会完整地复制一份数组给函数（Deep Copy）。
		在函数里修改数组，只会改那个副本，原本的数组不会变。
		硬伤：Go 语言中，[5]int 和 [10]int 是完全不同的两种数据类型。如果你的函数定义为 func demo(arr [5]int)，那你绝对不能传一个长度为 10 的数组进去，编译器会报错。
	*/

	// 关于向函数传递数组的测试
	myarr := [3]int{111, 222, 333}
	modifyArr(myarr)
	fmt.Println("Not Pointer:", myarr)
	modifyArrByPointer(&myarr)
	fmt.Println("With Pointer:", myarr)
	/*
		Go 的数组其实挺笨重的：1.传参会发生大量内存拷贝（性能低）。2.长度写死在类型里，函数很难通用（一个函数不能同时处理长度为 5 和 10 的数组）。
		实际代码的使用上我们一般用Slice切片，后面会学，暂时只做一个简单的了解
	*/
}

// 这种方式是值传递，跟C++不一样
func modifyArr(arr [3]int) {
	arr[0] = 999
}

// 如果你想在函数里真的修改外面的数组，必须传指针（和C++类似）
func modifyArrByPointer(arr *[3]int) {
	arr[0] = 888
}

/*--------------------------------------------------------------------------------------------------*/
/*
demo2 go语言指针
我们都知道，变量是一种使用方便的占位符，用于引用计算机内存地址。
Go 语言的取地址符是 &，放到一个变量前使用就会返回相应变量的内存地址。

1.什么是指针
一个指针变量指向了一个值的内存地址。
类似于变量和常量，在使用指针前你需要声明指针。指针声明格式如下：
var var_name *var-type
var ip *int         指向整型
var fp *float32     指向浮点型

2. 如何使用指针
指针使用流程：
-定义指针变量。
-为指针变量赋值。
-访问指针变量中指向地址的值。
--在指针类型前面加上 * 号（前缀）来获取指针所指向的内容。

Go 空指针
当一个指针被定义后没有分配到任何变量时，它的值为 nil。
nil 指针也称为空指针。
nil在概念上和其它语言的null、None、nil、NULL一样，都指代零值或空值。
空指针判断：
if(ptr != nil)     ptr 不是空指针
if(ptr == nil)     ptr 是空指针
*/

func demo2() {
	var a int = 42
	var ip *int
	ip = &a
	fmt.Printf("a 变量的地址是: %x\n", &a)

	/* 指针变量的存储地址 */
	fmt.Printf("ip 变量储存的指针地址: %x\n", ip)

	/* 使用指针访问值 */
	fmt.Printf("*ip 变量的值: %d\n", *ip)

	// 空指针测试
	var ptr_ *int
	fmt.Println("ptr的值为%x", ptr_)

	/* 补充部分：指针数组 & 二级指针** & 向函数传递指针参数（同上） */
	x, y, z := 7, 8, 9
	var ptrArr [3]*int
	ptrArr[0] = &x
	ptrArr[1] = &y
	ptrArr[2] = &z
	for i := 0; i < 3; i++ {
		fmt.Printf("ptrArr[%d] 指向的值是：%d\n", i, *ptrArr[i])
	}
	*ptrArr[0] = 999
	fmt.Printf("修改后变量 x 的值: %d\n", x) // 输出 999

	// 二级指针
	var ptr *int
	var pptr **int
	ptr = &x    // ptr 存 a 的地址
	pptr = &ptr // pptr 存 ptr 的地址
	fmt.Printf("变量 x = %d\n", a)
	fmt.Printf("指针 ptr 指向的值 (*ptr) = %d\n", *ptr)
	fmt.Printf("二级指针 pptr 最终指向的值 (**pptr) = %d\n", **pptr)
	/* C++ 差异提示： 虽然 Go 支持二级指针，但在 Go 的实际业务开发（尤其是 Web 后端）中，极少会像 C++ 那样用到三级、四级指针。通常一级指针处理结构体传递就足够了。 */

	// 向函数传递指针参数
	f, g := 100, 200
	fmt.Println("Before Swap(f, g):(%d, %d)", f, g)
	swap(&f, &g)
	fmt.Println("After Swap(f, g):(%d, %d)", f, g)
	/*
		the same as
		a, b = b, a
	*/
}

func swap(x, y *int) {
	temp := *x
	*x = *y
	*y = temp
}

/*--------------------------------------------------------------------------------------------------*/
/*
demo3 go语言结构体
Go 语言中数组可以存储同一类型的数据，但在结构体中我们可以为不同项定义不同的数据类型。
结构体是由一系列具有相同类型或不同类型的数据构成的数据集合。
结构体表示一项记录，比如保存图书馆的书籍记录，每本书有以下属性：
- Title ：标题
- Author ： 作者
- Subject：学科
- ID：书籍ID

1.定义结构体
结构体定义需要使用 type 和 struct 语句。struct 语句定义一个新的数据类型，结构体中有一个或多个成员。type 语句设定了结构体的名称。结构体的格式如下：
type struct_variable_type struct {
   member definition
   member definition
   ...
   member definition
}

一旦定义了结构体类型，它就能用于变量的声明，语法格式如下：
variable_name := structure_variable_type {value1, value2...valuen}
或
variable_name := structure_variable_type { key1: value1, key2: value2..., keyn: valuen}

2.访问结构体成员
如果要访问结构体成员，需要使用点号 . 操作符，格式为：
结构体.成员名"

3.结构体作为函数参数
你可以像其他数据类型一样将结构体类型作为参数传递给函数。

4.结构体指针
你可以定义指向结构体的指针类似于其他指针变量，格式如下：
var struct_pointer *Books

以上定义的指针变量可以存储结构体变量的地址。查看结构体变量地址，可以将 & 符号放置于结构体变量前：
struct_pointer = &Book1

使用结构体指针访问结构体成员，使用 "." 操作符：
struct_pointer.title
*/
type Book struct {
	title   string
	author  string
	subject string
	book_id int
}

func demo3() {
	// 创建一个新的结构体
	fmt.Println(Book{"Go_Study", "www.book.com", "Go_language", 1001})
	// 也可以使用 key => value 格式
	fmt.Println(Book{title: "Go 语言", author: "www.runoob.com", subject: "Go 语言教程", book_id: 6495407})
	// 忽略的字段为 0 或 空
	fmt.Println(Book{title: "Go 语言", author: "www.runoob.com"})

	var Book1, Book2 Book
	/* book 1 描述 */
	Book1.title = "Go 语言"
	Book1.author = "www.runoob.com"
	Book1.subject = "Go 语言教程"
	Book1.book_id = 6495407

	/* book 2 描述 */
	Book2.title = "Python 教程"
	Book2.author = "www.runoob.com"
	Book2.subject = "Python 语言教程"
	Book2.book_id = 6495700

	/* 打印 Book1 信息 */
	printBook_(&Book1)

	/* 打印 Book2 信息 */
	printBook(Book2)
}

// 结构体作为函数参数,值传递
func printBook(book Book) {
	fmt.Printf("Book title : %s\n", book.title)
	fmt.Printf("Book author : %s\n", book.author)
	fmt.Printf("Book subject : %s\n", book.subject)
	fmt.Printf("Book book_id : %d\n", book.book_id)
}

// 结构体指针作为函数参数
func printBook_(book *Book) {
	fmt.Printf("Book title : %s\n", book.title)
	fmt.Printf("Book author : %s\n", book.author)
	fmt.Printf("Book subject : %s\n", book.subject)
	fmt.Printf("Book book_id : %d\n", book.book_id)
}

/*--------------------------------------------------------------------------------------------------*/
/*
demo4 go语言切片

0.0 和数组的区分！！！
[]有数字 = 数组（定长），[]无数字 = 切片（变长）

Go 语言切片是对数组的抽象。
Go 数组的长度不可改变，在特定场景中这样的集合就不太适用，Go 中提供了一种灵活，功能强悍的内置类型切片("动态数组")，与数组相比切片的长度是不固定的，可以追加元素，在追加时可能使切片的容量增大。

1.定义切片
1.1你可以声明一个未指定大小的数组来定义切片：
var identifier []type	切片不需要说明长度。

1.2或使用 make() 函数来创建切片:
var slice1 []type = make([]type, len)
也可以简写为
slice1 := make([]type, len)
也可以指定容量，其中 capacity 为可选参数。
make([]T, length, capacity)
这里 len 是数组的长度并且也是切片的初始长度。

2.切片初始化
s :=[] int {1,2,3 }
直接初始化切片，[] 表示是切片类型，{1,2,3} 初始化值依次是 1,2,3，其 cap=len=3。
- s := arr[:]
初始化切片 s，是数组 arr 的引用。

- s := arr[startIndex:endIndex]
将 arr 中从下标 startIndex 到 endIndex-1 下的元素创建为一个新的切片。

- s := arr[startIndex:]
默认 endIndex 时将表示一直到arr的最后一个元素。

- s := arr[:endIndex]
默认 startIndex 时将表示从 arr 的第一个元素开始。

- s1 := s[startIndex:endIndex]
通过切片 s 初始化切片 s1。

- s :=make([]int,len,cap)
通过内置函数 make() 初始化切片s，[]int 标识为其元素类型为 int 的切片。

3.空(nil)切片
一个切片在未初始化之前默认为 nil

4.切片截取
可以通过设置下限及上限来设置截取切片 [lower-bound:upper-bound]

5.append() 和 copy() 函数
如果想增加切片的容量，我们必须创建一个新的更大的切片并把原分片的内容都拷贝过来。
下面的代码描述了从拷贝切片的 copy 方法和向切片追加新元素的 append 方法。
*/

func demo4() {
	/*
		切片是可索引的，并且可以由 len() 方法获取长度。
		切片提供了计算容量的方法 cap() 可以测量切片最长可以达到多少。
	*/
	var numbers = make([]int, 3, 5)
	printSlice(numbers)

	// 空(nil)切片
	var numbers_2 []int
	printSlice(numbers_2)
	if numbers == nil {
		fmt.Println("it's a nil Slice")
	}

	// 切片截取
	numbers_3 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	printSlice(numbers_3)
	fmt.Println("numbers_3[1:4] ==", numbers_3[1:4]) // 左闭右开

	// 默认下限为0, 上限为 len(s)
	fmt.Println("numbers_3[:3] ==", numbers_3[:3])
	fmt.Println("numbers_3[4:] ==", numbers_3[4:])

	number1 := numbers_3[:2] // 打印子切片从索引  0(包含) 到索引 2(不包含)
	printSlice(number1)

	number2 := numbers_3[2:5] // 打印子切片从索引 2(包含) 到索引 5(不包含
	printSlice(number2)
	// 上面两个就是创建子切片再打印的方式

	// append() & copy() 函数
	var nums []int
	printSlice(nums)
	nums = append(nums, 0)
	printSlice(nums)

	nums = append(nums, 1) //如果是要添加多个元素就多写几个, 比如append(nums, 1, 2, 3, 4)
	printSlice(nums)

	/* 创建切片 nums1 是之前切片的两倍容量*/
	nums_1 := make([]int, len(nums), (cap(nums) * 2))
	copy(nums, nums_1) // copy
	printSlice(nums_1)
}

// 打印切片函数
func printSlice(x []int) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}

/*--------------------------------------------------------------------------------------------------*/
/*
demo5 go语言范围
Go 语言中 range 关键字用于 for 循环中迭代数组(array)、切片(slice)、通道(channel)或集合(map)的元素。
在数组和切片中它返回元素的索引和索引对应的值，在集合中返回 key-value 对。
for 循环的 range 格式可以对 slice、map、数组、字符串等进行迭代循环。格式如下：
for key, value := range oldMap {
    newMap[key] = value
}
以上代码中的 key 和 value 是可以省略。
如果只想读取 key，格式如下：
- for key := range oldMap
或者这样：
- for key, _ := range oldMap
如果只想读取 value，格式如下：
- for _, value := range oldMap

字符串
range 迭代字符串时，返回每个字符的索引和 Unicode 代码点（rune）。
*/

func demo5() {
	var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	// string range
	for i, c := range "hello" {
		fmt.Printf("index: %d, char: %c\n", i, c)
	}

	map1 := map[int]float32{
		1: 1.0,
		2: 2.0,
		3: 3.0,
		4: 4.0,
	}

	for key, value := range map1 {
		fmt.Printf("key is: %d - value is: %f\n", key, value)
	}
	for key := range map1 {
		fmt.Printf("key is: %d", key)
	}
	for _, value := range map1 {
		fmt.Printf("value is: %f\n", value)
	}
}

func demo5_2() {
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum:", sum)

	for i, num := range nums {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}
	// range 也可以用在map键值对上
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Println("%s -> %s\n", k, v)
	}

	//range也可以用来枚举 Unicode 字符串。第一个参数是字符的索引，第二个是字符（Unicode的值）本身。
	for i, c := range "go" {
		fmt.Println(i, c)
	}
}

/*--------------------------------------------------------------------------------------------------*/
/*
demo6 go语言Map(集合)
-Map 是一种无序的键值对的集合。
--Map 最重要的一点是通过 key 来快速检索数据，key 类似于索引，指向数据的值。
--Map 是一种集合，所以我们可以像迭代数组和切片那样迭代它。不过，Map 是无序的，遍历 Map 时返回的键值对的顺序是不确定的。

在获取 Map 的值时，如果键不存在，返回该类型的零值，例如 int 类型的零值是 0，string 类型的零值是 ""。
Map 是引用类型，如果将一个 Map 传递给一个函数或赋值给另一个变量，它们都指向同一个底层数据结构，因此对 Map 的修改会影响到所有引用它的变量。

1.定义 Map
可以使用内建函数 make 或使用 map 关键字来定义 Map:
// 使用 make 函数
map_variable := make(map[KeyType]ValueType, initialCapacity)
其中 KeyType 是键的类型，ValueType 是值的类型，initialCapacity 是可选的参数，用于指定 Map 的初始容量。Map 的容量是指 Map 中可以保存的键值对的数量，当 Map 中的键值对数量达到容量时，Map 会自动扩容。如果不指定 initialCapacity，Go 语言会根据实际情况选择一个合适的值。

2.创建示例
m := make(map[string]int)  //创建一个空map
m1 := make(map[string]int, 10) // 创建一个初始容量为 10 的 Map

-也可以使用字面量创建 Map：
m := map[string]int{
	"apple": 1,
	"banana": 2,
	"orange": 3,
}

-获取元素
v1 := m["apple"]  // Map 的「键（key）」查找并获取对应「值（value）」
v2, ok := m["pear"]  // 如果键不存在，ok 的值为 false，v2 的值为该类型的零值

-修改元素：
// 修改键值对
m["apple"] = 5

-获取 Map 的长度：
len := len(m)

-遍历 Map：
for k, v := range m {
    fmt.Printf("key=%s, value=%d\n", k, v)
}

-删除元素：
delete(m, "banana")

*/

func demo6() {
	var siteMap map[string]string
	siteMap = make(map[string]string)

	/* map 插入 key - value 对,各个国家对应的首都 */
	siteMap["Google"] = "谷歌"
	siteMap["Baidu"] = "百度"
	siteMap["Wiki"] = "维基百科"

	/*使用键输出地图值 */
	for site := range siteMap {
		fmt.Println(site, "首都是", siteMap[site])
	}

	// 查看元素在集合中是否存在
	name, ok := siteMap["Facebook"]
	if ok {
		fmt.Println("Facebook 的 站点是", name)
	} else {
		fmt.Println("Facebook 站点不存在")
	}

	/* 创建map */
	countryCapitalMap := map[string]string{"France": "Paris", "Italy": "Rome", "Japan": "Tokyo", "India": "New delhi"}
	fmt.Println("原始地图")
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是：", countryCapitalMap[country])
	}

	// 删除元素
	delete(countryCapitalMap, "France")
	fmt.Println("法国条目被删除")
	fmt.Println("删除元素后地图")

	/*打印地图*/
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是", countryCapitalMap[country])
	}
}

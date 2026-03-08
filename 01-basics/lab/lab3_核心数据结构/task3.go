package main

import "fmt"

// 任务3：即时计票系统
func task3() {
	// 包含重复的名字的投票箱
	votes := []string{"Alice", "Bob", "Alice", "Charlie", "Alice", "Bob", "Charlie", "David", "Alice"}
	counter := make(map[string]int)
	for _, v := range votes {
		_, ok := counter[v] //检查key是否存在：value是对应的值，ok是bool
		if !ok {
			counter[v] = 1
		} else {
			counter[v] += 1
		}
	}
	var (
		max_num  int
		max_name string
	)
	for name, count := range counter {
		fmt.Printf("Candidate [%s]: [Count] %d votes\n", name, count)
		if count > max_num {
			max_num = count
			max_name = name
		}
	}
	fmt.Printf("票数最多的候选人为%s,票数为%d\n", max_name, max_num)
}

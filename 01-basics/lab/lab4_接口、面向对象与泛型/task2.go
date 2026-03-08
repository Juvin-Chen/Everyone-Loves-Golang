package main

import (
	"fmt"
)

// 构造一个通用泛型栈，类似于C++的std::stack
type Stack[T any] struct {
	elements []T
}

// Push()入栈，方法接收者也要带 [T]
func (s *Stack[T]) Push(v T) {
	s.elements = append(s.elements, v)
}

// Pop()出栈
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	} else {
		val := s.elements[len(s.elements)-1]
		s.elements = s.elements[:len(s.elements)-1]
		return val, true
	}
}

// 查看栈顶元素但不弹出
func (s *Stack[T]) Peek() T {
	return s.elements[len(s.elements)-1]
}

func task2() {
	fmt.Println("通用泛型栈:")
	stack1 := Stack[int]{}
	for i := 1; i <= 3; i++ {
		stack1.Push(i)
	}
	fmt.Println("开始进行int类型的通用栈的弹出：")
	for {
		if value, ok := stack1.Pop(); ok {
			fmt.Println(value)
		} else {
			break
		}
	}
	stack2 := Stack[string]{}
	stack2.Push("Go")
	stack2.Push("Is")
	stack2.Push("Good")
	fmt.Println("开始进行string类型通用栈的弹出：")
	for {
		if value, ok := stack2.Pop(); ok {
			fmt.Println(value)
		} else {
			break
		}
	}
}

// 定义数据模型和内存存储

package store

import (
	"errors"
	"sync"
	"time"
)

// Message 留言模型
type Message struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MemoryStore 内存存储
type MemoryStore struct {
	// RWMutex：读写锁，比普通锁更高级
	// 读：多人可以同时看（不冲突）
	// 写：只能一个人改（必须独占）
	mu       sync.RWMutex
	messages []Message
	nextID   int // 留言版里的留言 ID
}

// NewMemoryStore 创建存储实例
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		messages: []Message{},
		nextID:   1, // 初始化第一条留言的 id 为1
	}
}

// Add 添加留言
func (s *MemoryStore) Add(nickname, content string) Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	msg := Message{
		ID:        s.nextID,
		Nickname:  nickname,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.messages = append(s.messages, msg)
	s.nextID++
	return msg
}

// GetAll 获取所有留言（按时间倒序）
func (s *MemoryStore) GetAll() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// 复制一份，避免外部修改
	result := make([]Message, len(s.messages))
	copy(result, s.messages)
	// 倒序，分别从首尾出发不断进行交换
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// DeleteByID 删除留言（按留言 ID）
func (s *MemoryStore) DeleteByID(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.messages {
		if s.messages[i].ID == id {
			s.messages = append(s.messages[:i], s.messages[i+1:]...)
			return nil
		}
	}
	return errors.New("留言不存在")
}

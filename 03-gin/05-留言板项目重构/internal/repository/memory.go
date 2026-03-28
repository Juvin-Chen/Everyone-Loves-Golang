package repository

import (
	"message-board-gin/internal/model"
	"sync"
	"time"
)

type MessageRepository struct {
	mu       sync.RWMutex
	messages []model.Message
	nextID   int
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		messages: []model.Message{},
		nextID:   1,
	}
}

func (r *MessageRepository) Add(msg model.Message) model.Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	msg.ID = r.nextID
	msg.CreatedAt = time.Now()
	r.messages = append(r.messages, msg)
	r.nextID++
	return msg
}

func (r *MessageRepository) GetAll() []model.Message {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.Message, len(r.messages))
	copy(result, r.messages)
	// 倒序
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func (r *MessageRepository) DeleteByID(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, msg := range r.messages {
		if msg.ID == id {
			r.messages = append(r.messages[:i], r.messages[i+1:]...)
			return true
		}
	}

	return false
}

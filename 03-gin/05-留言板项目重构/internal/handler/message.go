package handler

import (
	"message-board-gin/internal/model"
	"message-board-gin/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	repo *repository.MessageRepository
}

func NewMessageHandler(repo *repository.MessageRepository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

func (h *MessageHandler) Index(c *gin.Context) {
	messages := h.repo.GetAll()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Messages": messages,
	})
}

func (h *MessageHandler) NewForm(c *gin.Context) {
	c.HTML(http.StatusOK, "new.html", nil)
}

func (h *MessageHandler) Create(c *gin.Context) {
	var msg model.Message
	if err := c.ShouldBind(&msg); err != nil {
		c.HTML(http.StatusBadRequest, "new.html", gin.H{
			"Error": "昵称和内容不能为空",
		})
		return
	}
	h.repo.Add(msg)
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *MessageHandler) Delete(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	h.repo.DeleteByID(id)
	c.Redirect(http.StatusSeeOther, "/")
}

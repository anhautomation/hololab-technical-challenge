package httpadapter

import (
	"hololab-core-system/internal/app/chat"
	"hololab-core-system/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerChat(r *gin.Engine, svc chat.Service) {
	r.GET("/api/bots/:id/messages", func(c *gin.Context) {
		msgs, err := svc.History(c.Request.Context(), c.Param("id"))
		if err == domain.ErrNotFound {
			c.JSON(404, gin.H{"error": "bot not found"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, msgs)
	})

	r.POST("/api/bots/:id/messages", func(c *gin.Context) {
		var req SendMessageRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		reply, err := svc.Send(c.Request.Context(), c.Param("id"), req.Message)
		if err == domain.ErrInvalidInput {
			c.JSON(400, gin.H{"error": "message is required"})
			return
		}
		if err == domain.ErrNotFound {
			c.JSON(404, gin.H{"error": "bot not found"})
			return
		}
		if err == domain.ErrLLMUnavailable {
			c.JSON(502, gin.H{"error": "llm unavailable"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reply": reply})
	})

	r.POST("/api/bots/:id/messages/reset", func(c *gin.Context) {
		err := svc.Reset(c.Request.Context(), c.Param("id"))
		if err == domain.ErrNotFound {
			c.JSON(404, gin.H{"error": "bot not found"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})
}

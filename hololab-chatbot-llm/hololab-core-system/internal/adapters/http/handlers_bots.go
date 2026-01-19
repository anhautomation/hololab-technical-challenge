package httpadapter

import (
	"hololab-core-system/internal/app/bots"
	"hololab-core-system/internal/domain"

	"github.com/gin-gonic/gin"
)

func registerBots(r *gin.Engine, svc bots.Service) {
	r.GET("/api/bots", func(c *gin.Context) {
		out, err := svc.List(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, out)
	})

	r.POST("/api/bots", func(c *gin.Context) {
		var req CreateBotRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		b, err := svc.Create(c.Request.Context(), bots.CreateInput(req))
		if err == domain.ErrInvalidInput {
			c.JSON(400, gin.H{"error": "missing required fields: name, job, bio, style"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, b)
	})

	r.GET("/api/bots/:id", func(c *gin.Context) {
		b, err := svc.Get(c.Request.Context(), c.Param("id"))
		if err == domain.ErrNotFound {
			c.JSON(404, gin.H{"error": "bot not found"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, b)
	})

	r.DELETE("/api/bots/:id", func(c *gin.Context) {
		n, err := svc.Delete(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"ok": true, "deleted": n})
	})
}

package httpadapter

import (
	"hololab-core-system/internal/app/bots"
	"hololab-core-system/internal/app/chat"
	"hololab-core-system/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	Bots bots.Service
	Chat chat.Service
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), CORSMiddleware())

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	registerBots(r, deps.Bots)
	registerChat(r, deps.Chat)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	_ = domain.ErrNotFound
	return r
}

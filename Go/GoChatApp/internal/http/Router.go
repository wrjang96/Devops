package httpx

import (
	"go-chat-mem/internal/auth"
	"go-chat-mem/internal/config"
	"go-chat-mem/internal/rooms"
	"go-chat-mem/internal/users"
	"go-chat-mem/internal/ws"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, uh users.Handler, rh rooms.Handler, hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	// 최소 CORS (Vue dev 서버에서 쿠키 포함 요청 대비)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", cfg.CorsOrigin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// public
	r.POST("/register", uh.Register)
	r.POST("/login", uh.Login)
	r.POST("/refreshToken", uh.RefreshToken)

	// WebSocket
	r.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// Room Stats
	r.GET("/rooms", func(c *gin.Context) {
		c.JSON(200, hub.GetStats())
	})

	// protected
	p := r.Group("/", auth.Middleware(cfg.JWTSecret))
	{
		p.GET("/chatrooms", rh.List)
		p.POST("/chatrooms", rh.Create)
		p.POST("/chatrooms/:id/join", rh.Join)
		p.POST("/chatrooms/:id/leave", rh.Leave)
		p.GET("/chatrooms/:id/messages", rh.GetMessages)
	}
	return r
}

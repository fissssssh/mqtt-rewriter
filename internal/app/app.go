package app

import (
	"mqtt-rewriter/internal/app/config"
	"mqtt-rewriter/internal/rewriter"

	"github.com/gin-gonic/gin"
)

func Run() {
	config.Init()
	client := *rewriter.Instance()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		client.Publish("test", 0, false, "123")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

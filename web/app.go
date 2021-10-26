package web

import (
	"log"
	"mqtt-rewriter/config"
	"mqtt-rewriter/rewriter"

	"github.com/gin-gonic/gin"
)

func Run() {
	err := config.Init()
	if err != nil {
		log.Fatalf("config init failed: %s", err)
	}
	_ = *rewriter.Instance()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

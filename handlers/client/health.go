package client

import "github.com/gin-gonic/gin"

func Health(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}

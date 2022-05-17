package client

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/handlers/client"
	"github.com/gin-gonic/gin"
)

func Routers(g *gin.RouterGroup) {
	r := g.Group("/")
	r.GET("/", client.Health)
	r.PUT("/upload", client.Upload)
	r.POST("/json", client.PostJson)
}

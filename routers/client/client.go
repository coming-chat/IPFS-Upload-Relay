package client

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/handlers/client"
	"github.com/gin-gonic/gin"
)

func Routers(g *gin.RouterGroup) {
	r := g.Group("/")
	r.GET("/", client.Health)
	r.POST("/upload", client.UploadFile)
	r.PUT("/upload", client.UploadFile)
	r.POST("/json", client.UploadJson)
	r.PUT("/json", client.UploadJson)
	r.POST("/video", client.UploadVideo)
}

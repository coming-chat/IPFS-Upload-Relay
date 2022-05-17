package inits

import (
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/routers"
	"github.com/NaturalSelectionLabs/IPFS-Upload-Relay/routers/client"
	"github.com/gin-gonic/gin"
)

func r(e *gin.Engine) {
	gBase := e.Group("/")
	client.Routers(gBase)
}

func Routers() *gin.Engine {
	routers.Include(r)
	return routers.Init()
}

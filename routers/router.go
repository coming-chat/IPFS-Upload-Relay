package routers

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type Option func(*gin.Engine)

var options []Option

// Include : Register routers
func Include(opts ...Option) {
	options = append(options, opts...)
}

func Init(middleware ...gin.HandlerFunc) *gin.Engine {

	if strings.Contains(os.Getenv("MODE"), "prod") {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	for _, mid := range middleware {
		r.Use(mid)
	}

	for _, opt := range options {
		opt(r)
	}

	return r
}

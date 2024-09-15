package mappings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/og-peek/handout/pkg/controllers"
)

var Router *gin.Engine

func CreateUrlMappings() {
	Router = gin.Default()

	Router.Use(controllers.Cors())

	Router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	v1 := Router.Group("/api/v1")
	{
		v1.GET("/get", controllers.GetUrl)
	}
}

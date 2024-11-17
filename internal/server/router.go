package server

import (
	"github.com/gin-gonic/gin"
	"masmaint-cg/internal/controller"
)


func SetRouter(r *gin.Engine) {
	rc := controller.NewRootController()
		
	r.GET("/", rc.IndexPage)
	r.POST("/generate", rc.PostGenerate)
}
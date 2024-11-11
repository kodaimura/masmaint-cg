package controller

import (
	"github.com/gin-gonic/gin"
)


func SetRouter(r *gin.Engine) {

	//render HTML or redirect
	rc := NewRootController()
		
	r.GET("/", rc.indexPage)
	r.POST("/generate", rc.postGenerate)
}
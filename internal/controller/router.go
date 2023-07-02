package controller

import (
	"github.com/gin-gonic/gin"
)


func SetRouter(r *gin.Engine) {

	//render HTML or redirect
	rc := newRootController()
		
	r.GET("/", rc.indexPage)
	r.POST("/csv", rc.postCsv)
}
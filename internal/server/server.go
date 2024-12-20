package server

import (
	"github.com/gin-gonic/gin"

	"masmaint-cg/config"
)

func Run() {
	cf := config.GetConfig()
	r := router()
	r.Run(":" + cf.AppPort)
}

func router() *gin.Engine {
	r := gin.Default()
	
	//TEMPLATE
	r.LoadHTMLGlob("web/template/*.html")

	//STATIC
	r.Static("/css", "web/static/css")
	r.Static("/js", "web/static/js")
	r.Static("/output", "./output")

	SetRouter(r)

	return r
}
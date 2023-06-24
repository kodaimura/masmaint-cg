package server

import (
	"github.com/gin-gonic/gin"

	"masmaint-cg/config"
	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/controller"
)

func Run() {
	cf := config.GetConfig()
	logger.SetAccessLogger()
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

	controller.SetRouter(r)

	return r
}

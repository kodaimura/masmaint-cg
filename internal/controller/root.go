package controller

import (
	"time"
	"github.com/gin-gonic/gin"

	"masmaint/internal/service"
	"masmaint/internal/shared/dto"
	"masmaint/pkg/utils"
)

type DdlParseService interface {
	Parse(path, dbtype string) []dto.Table
}

type rootController struct {
	dpServ *service.DdlParseService
}


func newRootController() *rootController {
	dpServ := service.NewDdlParseService()
	return &rootController{dpServ}
}


//GET /
func (ctr *rootController) indexPage(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

//POST /
func (ctr *rootController) postDdl(c *gin.Context) {
	file, _ := c.FormFile("ddl")
	randStr := utils.GenerateRandomString(10)
	path := "tmp/upload-" + time.Now().Format("2006-01-02-15-04-05") + "-" + randStr + ".sql"
	c.SaveUploadedFile(file, path)
	ctr.dpServ.Parse(path, "")
	c.HTML(200, "index.html", gin.H{})
}
package controller

import (
	"time"
	"github.com/gin-gonic/gin"

	"masmaint/internal/service"
	"masmaint/internal/shared/dto"
	"masmaint/pkg/utils"
)

type CsvParseService interface {
	Parse(path string) ([]dto.Table, []string)
}

type rootController struct {
	cpServ *service.CsvParseService
}


func newRootController() *rootController {
	cpServ := service.NewCsvParseService()
	return &rootController{cpServ}
}


//GET /
func (ctr *rootController) indexPage(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

//POST /
func (ctr *rootController) postDdl(c *gin.Context) {
	file, _ := c.FormFile("ddl")
	randStr := utils.GenerateRandomString(10)
	path := "tmp/upload-" + time.Now().Format("2006-01-02-15-04-05") + "-" + randStr + ".csv"
	c.SaveUploadedFile(file, path)
	ctr.cpServ.Parse(path)
	c.HTML(200, "index.html", gin.H{})
}
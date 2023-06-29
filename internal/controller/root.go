package controller

import (
	"time"
	"github.com/gin-gonic/gin"

	"masmaint-cg/internal/service"
	"masmaint-cg/internal/shared/dto"
	"masmaint-cg/pkg/utils"
)

type CsvParseService interface {
	Parse(path string) ([]dto.Table, []string)
}

type rootController struct {
	cpServ *service.CsvParseService
	genServ *service.GenerateService
}


func newRootController() *rootController {
	cpServ := service.NewCsvParseService()
	genServ := service.NewGenerateService()
	return &rootController{cpServ, genServ}
}


//GET /
func (ctr *rootController) indexPage(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

//POST /
func (ctr *rootController) postDdl(c *gin.Context) {
	file, _ := c.FormFile("ddl")
	randStr := utils.RandomString(10)
	path := "tmp/upload-" + time.Now().Format("2006-01-02-15-04-05") + "-" + randStr + ".csv"
	c.SaveUploadedFile(file, path)
	
	tables, errors := ctr.cpServ.Parse(path)

	if len(errors) != 0 {
		c.HTML(400, "index.html", gin.H{})
		c.Abort()
		return
	}

	zipPath, err := ctr.genServ.Generate(&tables, "golang", "postgresql")

	if err != nil {
		c.HTML(500, "index.html", gin.H{})
		c.Abort()
		return
	}
	 
	c.HTML(200, "index.html", gin.H{path: zipPath})
}
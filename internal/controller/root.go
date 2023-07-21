package controller

import (
	"time"
	"github.com/gin-gonic/gin"

	"masmaint-cg/internal/core/utils"
	"masmaint-cg/internal/service"
	"masmaint-cg/internal/shared/dto"
)

type CsvParseService interface {
	Parse(path string) ([]dto.Table, []string)
}

type GenerateService interface {
	Generate(tables *[]dto.Table, lang, rdbms string) (string, error)
}

type rootController struct {
	cpServ CsvParseService
	genServ GenerateService
}


func NewRootController() *rootController {
	cpServ := service.NewCsvParseService()
	genServ := service.NewGenerateService()
	return &rootController{cpServ, genServ}
}


//GET /
func (ctr *rootController) indexPage(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

//POST /csv
func (ctr *rootController) postCsv(c *gin.Context) {
	file, _ := c.FormFile("file")
	lang := c.PostForm("lang")
	rdbms := c.PostForm("rdbms")

	randStr := utils.RandomString(10)
	path := "tmp/upload-" + time.Now().Format("2006-01-02-15-04-05") + "-" + randStr + ".csv"
	c.SaveUploadedFile(file, path)
	
	tables, errors := ctr.cpServ.Parse(path)

	if len(errors) != 0 {
		c.JSON(400, gin.H{"errors":errors})
		c.Abort()
		return
	}

	zipPath, err := ctr.genServ.Generate(&tables, lang, rdbms)

	if err != nil {
		c.JSON(500, gin.H{"errors":[]string{"生成に失敗しました。"}})
		c.Abort()
		return
	}
	 
	c.JSON(200, gin.H{"path": zipPath})
}
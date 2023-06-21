package controller

import (
	"time"
	"github.com/gin-gonic/gin"

	"masmaint/pkg/utils"
)


type rootController struct {}


func newRootController() *rootController {
	return &rootController{}
}


//GET /
func (ctr *rootController) indexPage(c *gin.Context) {

	c.HTML(200, "index.html", gin.H{})
}

//POST /
func (ctr *rootController) postDdl(c *gin.Context) {
	file, _ := c.FormFile("ddl")
	randStr := utils.GenerateRandomString(10)
	fn := "tmp/upload-" + time.Now().Format("2006-01-02-15-04-05") + "-" + randStr + ".sql"
	c.SaveUploadedFile(file, fn)
	c.HTML(200, "index.html", gin.H{})
}
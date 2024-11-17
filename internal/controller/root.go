package controller

import (
	"io"
	"bytes"
	"github.com/gin-gonic/gin"

	"masmaint-cg/internal/module/generator"
)

type RootController struct {}


func NewRootController() *RootController {
	return &RootController{}
}


//GET /
func (ctr *RootController) IndexPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

//POST /generate
func (ctr *RootController) PostGenerate(c *gin.Context) {
	ddlFile, err := c.FormFile("ddl")
	//lang := c.PostForm("lang")
	rdbms := c.PostForm("rdbms")

	if err != nil {
		c.JSON(400, gin.H{"errors":[]string{"ファイルを取得できませんでした。"}})
		return
	}
	
	file, err := ddlFile.Open()
	if err != nil {
		c.JSON(500, gin.H{"errors":[]string{"ファイルを開けませんでした。"}})
		return
	}
	defer file.Close()
	
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(500, gin.H{"errors":[]string{"ファイル内容を読み込めませんでした。"}})
		return
	}
	
	ddl := buf.String()

	if err != nil {
		c.JSON(400, gin.H{"errors": []string{err.Error()}})
		return
	}
	gen, err := generator.NewGenerator(ddl, rdbms)
	if err != nil {
		c.JSON(400, gin.H{"errors":[]string{err.Error()}})
		return
	}

	zip, err := gen.Generate()
	if err != nil {
		c.JSON(500, gin.H{"errors":[]string{"生成に失敗しました。"}})
		return
	}
	 
	c.JSON(200, gin.H{"zip": zip})
}
package generator

import (
	"os"
	"fmt"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
)


type SourceGeneratorGolang struct {
	tables *[]dto.Table
	rdbms string
	path string
}

func NewSourceGeneratorGolang(tables *[]dto.Table, rdbms, path string) *SourceGeneratorGolang {
	return &SourceGeneratorGolang{
		tables, rdbms, path,
	}
}

func (serv *SourceGeneratorGolang) GenerateSource() error {

	if err := serv.generateSourceCmd(); err != nil {
		return err
	}
	if err := serv.generateSourceConfig(); err != nil {
		return err
	}
	if err := serv.generateSourceCore(); err != nil {
		return err
	}
	if err := serv.generateSourceController(); err != nil {
		return err
	}
	/*
	if err := serv.generateSourceDto(); err != nil {
		return err
	}
	if err := serv.generateSourceModel(); err != nil {
		return err
	}
	if err := serv.generateSourceService(); err != nil {
		return err
	}
	if err := serv.generateSourceWeb(); err != nil {
		return err
	}
	*/
	return nil	
}

func (serv *SourceGeneratorGolang) generateSourceCmd() error {
	source := "_originalcopy_/golang/cmd"
	destination := serv.path + "cmd/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err

}

func (serv *SourceGeneratorGolang) generateSourceConfig() error {
	source := "_originalcopy_/golang/config"
	destination := serv.path + "config/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateSourceCore() error {
	source := "_originalcopy_/golang/core"
	destination := serv.path + "core/"
	
	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateSourceController() error {
	path := serv.path + "controller/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateSourceControllerFiles(path)
}

func (serv *SourceGeneratorGolang) generateSourceControllerFiles(path string) error {
	if err := serv.generateSourceControllerFileRouter(path); err != nil {
		logger.LogError(err.Error())
		return err
	}

	for _, table := range *serv.tables {
		if err := serv.generateSourceControllerFile(&table, path); err != nil {
			logger.LogError(err.Error())
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateSourceControllerFileRouter(path string) error {
	code := "package controller\n\nimport (\n" +
		"\t\"github.com/gin-gonic/gin\"\n\n\t\"masmaint/core/auth\"\n)\n\n\n" +
		"func SetRouter(r *gin.Engine) {\n\n" +
		"\trm := r.Group(\"/mastertables\", auth.NoopAuthMiddleware())\n" +
		"\t{\n\t\trm.GET(\"/\", func(c *gin.Context) {\n" +
		"\t\t\tc.HTML(200, \"index.html\", gin.H{})\n\t\t})\n\n" 

	for _, table := range *serv.tables {
		tn := table.TableName
		tnc := SnakeToCamel(tn)
		tnp := SnakeToPascal(tn)

		code += fmt.Sprintf("\t\t%sController := New%sController()\n", tnc, tnp) + 
			fmt.Sprintf("\t\trm.GET(\"/api/%s\", %sController.Get%sPage)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.GET(\"/api/%s\", %sController.Get%s)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.POST(\"/api/%s\", %sController.Post%s)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.PUT(\"/api/%s\", %sController.Put%s)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.DELETE(\"/api/%s\", %sController.Delete%s)\n", tn, tnc, tnp) + 
			"\n"
	}

	code += "\t}\n}"
	return WriteFile(path + "router.go", code)
}

func (serv *SourceGeneratorGolang) generateSourceControllerFile(table *dto.Table, path string) error {
	code := "package controller\n\nimport (\n" +
		"\t\"github.com/gin-gonic/gin\"\n\n\t\"masmaint/service\"\n\t\"masmaint/dto\"\n)\n\n\n"

	tn := table.TableName
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code += fmt.Sprintf("type %sService interface {\n", tnp) + 
		fmt.Sprintf("\tCreate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%sDto *dto.%sDto) error\n", tni, tnp) +
		fmt.Sprintf("\tGetAll() ([]dto.%sDto, error)\n", tnp) +
		fmt.Sprintf("\tGetOne(%sDto *dto.%sDto) ([]dto.%sDto, error)\n", tni, tnp, tnp) +
		"}\n\n"

	code += fmt.Sprintf("type %sController struct {\n", tnp) +
		fmt.Sprintf("\t%sServ *service.%sService\n", tni, tnp) + 
		"}\n\n\n"

	code += fmt.Sprintf("func New%sController() *%sController {\n", tnp, tnp) +
		fmt.Sprintf("\t%sServ := service.New%sService()\n", tni, tnp) +
		fmt.Sprintf("\treturn &%sController{%sServ}\n", tnp, tni) +
		"}\n\n\n"

	code += fmt.Sprintf("//GET /%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Get%sPage(c *gin.Context) {\n", tnp, tnp) +
		fmt.Sprintf("\tc.HTML(200, \"%s.html\", gin.H{})\n", tn) +
		"}\n\n"

	code += fmt.Sprintf("//GET /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Get%s(c *gin.Context) {\n", tnp, tnp) +
		fmt.Sprintf("\tret, err := ctr.%sServ.GetAll()\n\n", tni) +
		"\tif err != nil {\n\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n\n"

	code += fmt.Sprintf("//POST /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Post%s(c *gin.Context) {\n", tnp, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) + 
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tret, err := ctr.%sServ.Create(&%sDto)\n\n", tni, tni) +
		"\tif err != nil {\n\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n\n"

	code += fmt.Sprintf("//PUT /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Put%s(c *gin.Context) {\n", tnp, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) + 
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tret, err := ctr.%sServ.Update(&%sDto)\n\n", tni, tni) +
		"\tif err != nil {\n\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n\n"

	code += fmt.Sprintf("//DELETE /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Delete%s(c *gin.Context) {\n", tnp, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) +
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tif err := ctr.%sServ.Delete(&%sDto); err != nil {\n", tni, tni) +
		"\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n"

	return WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
}

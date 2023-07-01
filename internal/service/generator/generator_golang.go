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
	if err := serv.generateSourceDto(); err != nil {
		return err
	}
	if err := serv.generateSourceService(); err != nil {
		return err
	}
	/*
	if err := serv.generateSourceModel(); err != nil {
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
		fmt.Sprintf("\tGetAll() ([]dto.%sDto, error)\n", tnp) +
		fmt.Sprintf("\tCreate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%sDto *dto.%sDto) error\n", tni, tnp) +
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

func (serv *SourceGeneratorGolang) generateSourceDto() error {
	path := serv.path + "dto/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateSourceDtoFiles(path)
}

func (serv *SourceGeneratorGolang) generateSourceDtoFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateSourceDtoFile(&table, path); err != nil {
			logger.LogError(err.Error())
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateSourceDtoFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	code := "package dto\n\n\n" + fmt.Sprintf("type %sDto struct {\n", tnp)

	for _, col := range table.Columns {
		cn := col.ColumnName
		cnp := SnakeToPascal(cn)
		code += fmt.Sprintf("\t%s any `json:\"%s\"`\n", cnp, cn)
	}

	code += "}"
	return WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
}

func (serv *SourceGeneratorGolang) generateSourceService() error {
	path := serv.path + "service/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateSourceServiceFiles(path)
}

func (serv *SourceGeneratorGolang) generateSourceServiceFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateSourceServiceFile(&table, path); err != nil {
			logger.LogError(err.Error())
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateSourceServiceFile(table *dto.Table, path string) error {
	code := "package dto\n\nimport (\n" +
		"\t\"errors\"\n\n\t\"masmaint/core/logger\"\n\t\"masmaint/model/entity\"\n" +
		"\t\"masmaint/model/dao\"\n\t\"masmaint/dto\"\n)\n\n\n"

	tn := table.TableName
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code += fmt.Sprintf("type %sDao interface {\n", tnp) + 
		fmt.Sprintf("\tSelectAll() ([]entity.%s, error)\n", tnp) +
		fmt.Sprintf("\tSelect(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tInsert(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%s *entity.%s) error\n", tni, tnp) +
		"}\n\n"

	code += fmt.Sprintf("type %sService struct {\n", tnp) +
		fmt.Sprintf("\t%sDao *dao.%sDao\n", tni, tnp) + 
		"}\n\n\n"

	code += fmt.Sprintf("func New%sService() *%sService {\n", tnp, tnp) +
		fmt.Sprintf("\t%sDao := service.New%sDao()\n", tni, tnp) +
		fmt.Sprintf("\treturn &%sService{%sDao}\n", tnp, tni) +
		"}\n\n\n"

	// *Service.GetAll()
	code += fmt.Sprintf("func (serv *%sService) GetAll() ([]dto.%sDto, error) {\n", tnp, tnp) +
		fmt.Sprintf("\trows, err := serv.%sDao.SelectAll()\n", tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		fmt.Sprintf("\t\treturn []dto.%sDto{}, errors.New(\"取得に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\tvar ret []dto.%sDto\n", tnp) +
		fmt.Sprintf("\tfor _, row := range rows {\n\t\tret = append(ret, row.To%sDto())\n\t}\n\n", tnp) +
		"\treturn ret, nil\n}\n\n\n"

	// *Service.GetOne()
	code += fmt.Sprintf("func (serv *%sService) GetOne(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnp, tni, tnp, tnp) +
		fmt.Sprintf("var %s *entity.%s = entity.New%s()\n", tni, tnp, tnp)
	
	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			cnp := SnakeToPascal(col.ColumnName)
			if isFirst {
				code += fmt.Sprintf("\tif %s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
				isFirst = false
			} else {
				code += "||\n"
				code += fmt.Sprintf("\t%s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
			}
 		}
	}
	code += "{\n"
	code += fmt.Sprintf("\t\treturn dto.%sDto{}, errors.New(\"不正な値があります。\")\n\t}\n\n", tnp)

	code += fmt.Sprintf("\trow, err := serv.%sDao.Select(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		fmt.Sprintf("\t\treturn dto.%sDto{}, errors.New(\"取得に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\treturn row.To%sDto(), nil\n", tnp) +
		"}\n\n\n"

	// *Service.Create()
	code += fmt.Sprintf("func (serv *%sService) Create(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnp, tni, tnp, tnp) +
		fmt.Sprintf("var %s *entity.%s = entity.New%s()\n", tni, tnp, tnp)
	
	isFirst = true
	for _, col := range table.Columns {
		if !col.IsAuto && !col.IsReadOnly {
			cnp := SnakeToPascal(col.ColumnName)
			if isFirst {
				code += fmt.Sprintf("\tif %s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
				isFirst = false
			} else {
				code += "||\n"
				code += fmt.Sprintf("\t%s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
			}
 		}
	}
	code += "{\n"
	code += fmt.Sprintf("\t\treturn dto.%sDto{}, errors.New(\"不正な値があります。\")\n\t}\n\n", tnp)

	code += fmt.Sprintf("\trow, err := serv.%sDao.Insert(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		fmt.Sprintf("\t\treturn dto.%sDto{}, errors.New(\"登録に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\treturn row.To%sDto(), nil\n", tnp) +
		"}\n\n\n"

	// *Service.Update()
	code += fmt.Sprintf("func (serv *%sService) Update(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnp, tni, tnp, tnp) +
		fmt.Sprintf("var %s *entity.%s = entity.New%s()\n", tni, tnp, tnp)
	
	isFirst = true
	for _, col := range table.Columns {
		if !col.IsReadOnly {
			cnp := SnakeToPascal(col.ColumnName)
			if isFirst {
				code += fmt.Sprintf("\tif %s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
				isFirst = false
			} else {
				code += "||\n"
				code += fmt.Sprintf("\t%s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
			}
 		}
	}
	code += "{\n"
	code += fmt.Sprintf("\t\treturn dto.%sDto{}, errors.New(\"不正な値があります。\")\n\t}\n\n", tnp)

	code += fmt.Sprintf("\trow, err := serv.%sDao.Update(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		fmt.Sprintf("\t\treturn dto.%sDto{}, errors.New(\"更新に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\treturn row.To%sDto(), nil\n", tnp) +
		"}\n\n\n"

	// *Service.Delete()
	code += fmt.Sprintf("func (serv *%sService) Delete(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnp, tni, tnp, tnp) +
		fmt.Sprintf("var %s *entity.%s = entity.New%s()\n", tni, tnp, tnp)
	
	isFirst = true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			cnp := SnakeToPascal(col.ColumnName)
			if isFirst {
				code += fmt.Sprintf("\tif %s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
				isFirst = false
			} else {
				code += "||\n"
				code += fmt.Sprintf("\t%s.Set%s(%sDto.%s) != nil ", tni, cnp, tni, cnp)
			}
 		}
	}
	code += "{\n"
	code += fmt.Sprintf("\t\treturn dto.%sDto{}, errors.New(\"不正な値があります。\")\n\t}\n\n", tnp)

	code += fmt.Sprintf("\trow, err := serv.%sDao.Delete(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		"\t\treturn errors.New(\"削除に失敗しました。\")\n\t}\n\n" +
		"\treturn nil\n}\n\n\n"

	return WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
}

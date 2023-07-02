package generator

import (
	"os"
	"fmt"
	"strings"

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

	if err := serv.generateCmd(); err != nil {
		return err
	}
	if err := serv.generateConfig(); err != nil {
		return err
	}
	if err := serv.generateCore(); err != nil {
		return err
	}
	if err := serv.generateLog(); err != nil {
		return err
	}
	if err := serv.generateController(); err != nil {
		return err
	}
	if err := serv.generateDto(); err != nil {
		return err
	}
	if err := serv.generateService(); err != nil {
		return err
	}
	if err := serv.generateModel(); err != nil {
		return err
	}
	if err := serv.generateWeb(); err != nil {
		return err
	}

	return nil	
}

func (serv *SourceGeneratorGolang) generateCmd() error {
	source := "_originalcopy_/golang/cmd"
	destination := serv.path + "cmd/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateConfig() error {
	source := "_originalcopy_/golang/config"
	destination := serv.path + "config/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateCore() error {
	source := "_originalcopy_/golang/core"
	destination := serv.path + "core/"
	
	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateLog() error {
	source := "_originalcopy_/golang/log"
	destination := serv.path + "log/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateController() error {
	path := serv.path + "controller/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateControllerFiles(path)
}

func (serv *SourceGeneratorGolang) generateControllerFiles(path string) error {
	if err := serv.generateControllerFileRouter(path); err != nil {
		return err
	}

	for _, table := range *serv.tables {
		if err := serv.generateControllerFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateControllerFileRouter(path string) error {
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
			fmt.Sprintf("\t\trm.GET(\"/%s\", %sController.Get%sPage)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.GET(\"/api/%s\", %sController.Get%s)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.POST(\"/api/%s\", %sController.Post%s)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.PUT(\"/api/%s\", %sController.Put%s)\n", tn, tnc, tnp) +
			fmt.Sprintf("\t\trm.DELETE(\"/api/%s\", %sController.Delete%s)\n", tn, tnc, tnp) + 
			"\n"
	}

	code += "\t}\n}"

	err := WriteFile(path + "router.go", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateControllerFile(table *dto.Table, path string) error {
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
		"}\n\n"

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
		"\tc.JSON(200, gin.H{})\n}\n"

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateDto() error {
	path := serv.path + "dto/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateDtoFiles(path)
}

func (serv *SourceGeneratorGolang) generateDtoFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateDtoFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateDtoFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	code := "package dto\n\n\n" + fmt.Sprintf("type %sDto struct {\n", tnp)

	for _, col := range table.Columns {
		cn := col.ColumnName
		cnp := SnakeToPascal(cn)
		code += fmt.Sprintf("\t%s any `json:\"%s\"`\n", cnp, cn)
	}

	code += "}"

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateService() error {
	path := serv.path + "service/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateServiceFiles(path)
}

func (serv *SourceGeneratorGolang) generateServiceFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateServiceFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateServiceFile(table *dto.Table, path string) error {
	code := "package service\n\nimport (\n" +
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
		"}\n\n"

	code += fmt.Sprintf("func New%sService() *%sService {\n", tnp, tnp) +
		fmt.Sprintf("\t%sDao := dao.New%sDao()\n", tni, tnp) +
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
		fmt.Sprintf("\tvar %s *entity.%s = entity.New%s()\n\n", tni, tnp, tnp)
	
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
		fmt.Sprintf("\tvar %s *entity.%s = entity.New%s()\n\n", tni, tnp, tnp)
	
	isFirst = true
	for _, col := range table.Columns {
		if col.IsInsAble {
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
		fmt.Sprintf("\tvar %s *entity.%s = entity.New%s()\n\n", tni, tnp, tnp)
	
	isFirst = true
	for _, col := range table.Columns {
		if col.IsPrimaryKey || col.IsUpdAble {
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
	code += fmt.Sprintf("func (serv *%sService) Delete(%sDto *dto.%sDto) error {\n", tnp, tni, tnp) +
		fmt.Sprintf("\tvar %s *entity.%s = entity.New%s()\n\n", tni, tnp, tnp)
	
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
	code += "\t\treturn errors.New(\"不正な値があります。\")\n\t}\n\n"

	code += fmt.Sprintf("\terr := serv.%sDao.Delete(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		"\t\treturn errors.New(\"削除に失敗しました。\")\n\t}\n\n" +
		"\treturn nil\n}\n\n"

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateModel() error {
	if err := serv.generateEntity(); err != nil {
		return err
	}

	if err := serv.generateDao(); err != nil {
		return err
	}

	return nil
}

func (serv *SourceGeneratorGolang) generateEntity() error {
	path := serv.path + "model/entity/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateEntityFiles(path)
}

func (serv *SourceGeneratorGolang) generateEntityFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateEntityFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) getEntityFieldType(col *dto.Column) string {
	isNotNull := col.IsNotNull
	colType := col.ColumnType
	if colType == "s" || colType == "t" {
		if isNotNull {
			return "string"
		}
		return "sql.NullString"
	}
	if colType == "i" {
		if isNotNull {
			return "int64"
		}
		return "sql.NullInt64"
	}
	if colType == "f" {
		if isNotNull {
			return "float64"
		}
		return "sql.NullFloat64"
	}
	return ""
}

func (serv *SourceGeneratorGolang) generateEntityFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	code := "package entity\n\nimport (\n" +
		"\t\"database/sql\"\n\n\t\"masmaint/dto\"\n\t\"masmaint/core/utils\"\n)\n\n\n"

	code += fmt.Sprintf("type %s struct {\n", tnp)
	for _, col := range table.Columns {
		cn := col.ColumnName
		cnp := SnakeToPascal(cn)
		code += fmt.Sprintf("\t%s %s `db:\"%s\"`\n", cnp, serv.getEntityFieldType(&col), cn)
	}
	code += "}\n\n"
	code += serv.generateEntityFileCodeSetters(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateEntityFileCodeSetters(table *dto.Table) string {
	tnp := SnakeToPascal(table.TableName)

	code := fmt.Sprintf("func New%s() *%s {\n\treturn &%s{}\n}\n\n", tnp, tnp, tnp)
	for _, col := range table.Columns {
		code += serv.generateEntityFileCodeSetter(table, &col)
	}

	code += "\n"
	code += serv.generateEntityFileCodeToDto(table)
	
	return code
}

func (serv *SourceGeneratorGolang) generateEntityFileCodeSetter(table *dto.Table, col *dto.Column) string {
	tnp := SnakeToPascal(table.TableName)
	colType := serv.getEntityFieldType(col)
	cnp := SnakeToPascal(col.ColumnName)
	cnc := SnakeToCamel(col.ColumnName)

	code := fmt.Sprintf("func (e *%s) Set%s(%s any) error {\n", tnp, cnp, cnc)

	switch colType {
	case "string":
		code += fmt.Sprintf("\te.%s = utils.ToString(%s)\n\treturn nil\n}\n\n", cnp, cnc)

	case "int64":
		code += fmt.Sprintf("\tx, err := utils.ToInt64(%s)\n\tif err != nil {\n\t\treturn err\n\t}\n", cnc) +
			fmt.Sprintf("\te.%s = x\n\treturn nil\n}\n\n", cnp)

	case "float64":
		code += fmt.Sprintf("\tx, err := utils.ToFloat64(%s)\n\tif err != nil {\n\t\treturn err\n\t}\n", cnc) +
			fmt.Sprintf("\te.%s = x\n\treturn nil\n}\n\n", cnp)
			
	case "sql.NullString":
		if col.ColumnType == "t" {
			code += fmt.Sprintf("\tif %s == nil || %s == \"\" {\n", cnc, cnc)
		} else {
			code += fmt.Sprintf("\tif %s == nil {\n", cnc)
		}
		code += fmt.Sprintf("\t\te.%s.Valid = false\n\t\treturn nil\n\t}\n\n", cnp) +
			fmt.Sprintf("\te.%s.String = utils.ToString(%s)\n", cnp, cnc) +
			fmt.Sprintf("\te.%s.Valid = true\n\treturn nil\n}\n\n", cnp)

	case "sql.NullInt64":
		code += fmt.Sprintf("\tif %s == nil || %s == \"\" {\n", cnc, cnc) +
			fmt.Sprintf("\t\te.%s.Valid = false\n\t\treturn nil\n\t}\n\n", cnp) +
			fmt.Sprintf("\tx, err := utils.ToInt64(%s)\n\tif err != nil {\n\t\treturn err\n\t}\n", cnc) +
			fmt.Sprintf("\te.%s.Int64 = x\n\te.%s.Valid = true\n\treturn nil\n}\n\n", cnp, cnp)

	case "sql.NullFloat64":
		code += fmt.Sprintf("\tif %s == nil || %s == \"\" {\n", cnc, cnc) +
			fmt.Sprintf("\t\te.%s.Valid = false\n\t\treturn nil\n\t}\n\n", cnp) +
			fmt.Sprintf("\tx, err := utils.ToFloat64(%s)\n\tif err != nil {\n\t\treturn err\n\t}\n", cnc) +
			fmt.Sprintf("\te.%s.Float64 = x\n\te.%s.Valid = true\n\treturn nil\n}\n\n", cnp, cnp)
	}

	return code
}

func (serv *SourceGeneratorGolang) generateEntityFileCodeToDto(table *dto.Table) string {
	tnp := SnakeToPascal(table.TableName)
	tni := GetSnakeInitial(table.TableName)

	code := fmt.Sprintf("func (e *%s) To%sDto() dto.%sDto {\n", tnp, tnp, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp)
	for _, col := range table.Columns {
		colType := serv.getEntityFieldType(&col)
		cnp := SnakeToPascal(col.ColumnName)

		switch colType {
		case "string", "int64", "float64":
			code += fmt.Sprintf("\t%sDto.%s = e.%s\n", tni, cnp, cnp)

		case "sql.NullString":
			code += fmt.Sprintf("\tif e.%s.Valid != false {\n", cnp) +
				fmt.Sprintf("\t\t%sDto.%s = e.%s.String\n\t}\n", tni, cnp, cnp)

		case "sql.NullInt64":
			code += fmt.Sprintf("\tif e.%s.Valid != false {\n", cnp) +
				fmt.Sprintf("\t\t%sDto.%s = e.%s.Int64\n\t}\n", tni, cnp, cnp)
				
		case "sql.NullFloat64":
			code += fmt.Sprintf("\tif e.%s.Valid != false {\n", cnp) +
				fmt.Sprintf("\t\t%sDto.%s = e.%s.Float64\n\t}\n", tni, cnp, cnp)	
		}
	}

	code += fmt.Sprintf("\n\treturn %sDto\n}\n", tni)
	return code
}

func (serv *SourceGeneratorGolang) generateDao() error {
	path := serv.path + "model/dao/"

	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	return serv.generateDaoFiles(path)
}

func (serv *SourceGeneratorGolang) generateDaoFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateDaoFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateDaoFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	code := "package dao\n\nimport (\n" +
		"\t\"database/sql\"\n\n\t\"masmaint/core/db\"\n\t\"masmaint/model/entity\"\n)\n\n\n"

	code += fmt.Sprintf("type %sDao struct {\n\tdb *sql.DB\n}\n\n", tnp) +
		fmt.Sprintf("func New%sDao() *%sDao {\n", tnp, tnp) + 
		fmt.Sprintf("\tdb := db.GetDB()\n\treturn &%sDao{db}\n}\n\n\n", tnp)

	code += serv.generateDaoFileCodeSelectAll(table) + "\n\n"
	code += serv.generateDaoFileCodeSelect(table) + "\n\n"
	code += serv.generateDaoFileCodeInsert(table) + "\n\n"
	code += serv.generateDaoFileCodeUpdate(table) + "\n\n"
	code += serv.generateDaoFileCodeDelete(table) + "\n"

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) concatPrimaryKeyWithCommas(cols []dto.Column) string {
	var ls []string
	for _, col := range cols {
		if col.IsPrimaryKey {
			ls = append(ls, col.ColumnName)
		}
	}
	return strings.Join(ls, ", ")
}

func (serv *SourceGeneratorGolang) generateDaoFileCodeSelectAll(table *dto.Table) string {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) SelectAll() ([]entity.%s, error) {\n", tnp, tnp) +
		fmt.Sprintf("\tvar ret []entity.%s\n\n\trows, err := rep.db.Query(\n", tnp)

	code += "\t\t`SELECT\n"
	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\t\t\t%s", col.ColumnName)
		} else {
			code += fmt.Sprintf("\n\t\t\t,%s", col.ColumnName)
		}
	}
	code += fmt.Sprintf("\n\t\t FROM %s\n\t\t ORDER BY %s ASC", tn, serv.concatPrimaryKeyWithCommas(table.Columns))
	code += "`"
	code += ",\n\t)\n\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n"

	code += fmt.Sprintf("\tfor rows.Next() {\n\t\t%s := entity.%s{}\n\t\terr = rows.Scan(\n", tni, tnp)
	for _, col := range table.Columns {
		cnp := SnakeToPascal(col.ColumnName)
		code += fmt.Sprintf("\t\t\t&%s.%s,\n", tni, cnp)
	}
	code += fmt.Sprintf("\t\t)\n\t\tif err != nil {\n\t\t\tbreak\n\t\t}\n\t\tret = append(ret, %s)\n\t}\n\n", tni)

	code += "\treturn ret, err\n}\n"
	return code
}

func (serv *SourceGeneratorGolang) generateDaoFileCodeSelect(table *dto.Table) string {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Select(%s *entity.%s) (entity.%s, error) {\n", tnp, tni, tnp, tnp) +
		fmt.Sprintf("\tvar ret entity.%s\n\n\terr := rep.db.QueryRow(\n", tnp)

	code += "\t\t`SELECT\n"
	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\t\t\t%s\n", col.ColumnName)
		} else {
			code += fmt.Sprintf("\t\t\t,%s\n", col.ColumnName)
		}
	}
	code += fmt.Sprintf("\t\t FROM %s\n\t\t WHERE ", tn)

	bindCount := 0
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if bindCount == 1 {
				code += fmt.Sprintf("%s = $1", col.ColumnName)
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = $%d", col.ColumnName, bindCount)
			}
		}
	}
	code += "`"
	code += ",\n"

	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			code += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}

	code += "\t).Scan(\n"
	for _, col := range table.Columns {
		code += fmt.Sprintf("\t\t&ret.%s,\n", SnakeToPascal(col.ColumnName))
	}
	code += "\t)\n\n\treturn ret, err\n}\n" 

	return code
}

func (serv *SourceGeneratorGolang) concatBindVariableWithCommas(bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, fmt.Sprintf("$%d", i))
	}
	return strings.Join(ls, ",")
}

func (serv *SourceGeneratorGolang) generateDaoFileCodeInsert(table *dto.Table) string {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Insert(%s *entity.%s) (entity.%s, error) {\n", tnp, tni, tnp, tnp) +
		fmt.Sprintf("\tvar ret entity.%s\n\n\terr := rep.db.QueryRow(\n", tnp)

	code += fmt.Sprintf("\t\t`INSERT INTO %s (\n", tn)
	bindCount := 0
	for _, col := range table.Columns {
		if col.IsInsAble {
			bindCount += 1
			if bindCount == 1 {
				code += fmt.Sprintf("\t\t\t%s", col.ColumnName)
			} else {
				code += fmt.Sprintf("\n\t\t\t,%s", col.ColumnName)
			}
		}
	}
	code += fmt.Sprintf("\n\t\t ) VALUES (%s)\n\t\t RETURNING\n", serv.concatBindVariableWithCommas(bindCount))
	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\t\t\t%s", col.ColumnName)
		} else {
			code += fmt.Sprintf("\n\t\t\t,%s", col.ColumnName)
		}
	}
	code += "`"
	code += ",\n"

	for _, col := range table.Columns {
		if col.IsInsAble {
			code += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}

	code += "\t).Scan(\n"
	for _, col := range table.Columns {
		code += fmt.Sprintf("\t\t&ret.%s,\n", SnakeToPascal(col.ColumnName))
	}
	code += "\t)\n\n\treturn ret, err\n}\n" 

	return code
}

func (serv *SourceGeneratorGolang) generateDaoFileCodeUpdate(table *dto.Table) string {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Update(%s *entity.%s) (entity.%s, error) {\n", tnp, tni, tnp, tnp) +
		fmt.Sprintf("\tvar ret entity.%s\n\n\terr := rep.db.QueryRow(\n", tnp)

	code += fmt.Sprintf("\t\t`UPDATE %s\n\t\t SET\n", tn)
	bindCount := 0
	for _, col := range table.Columns {
		if col.IsUpdAble {
			bindCount += 1
			if bindCount == 1 {
				code += fmt.Sprintf("\t\t\t%s = $1", col.ColumnName)
			} else {
				code += fmt.Sprintf("\n\t\t\t,%s = $%d", col.ColumnName, bindCount)
			}
		}
	}
	code += "\n\t\t WHERE "

	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if isFirst {
				code += fmt.Sprintf("%s = $%d", col.ColumnName, bindCount)
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = $%d", col.ColumnName, bindCount)
			}
		}
	}
	code += "\n\t\t RETURNING \n"

	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\t\t\t%s", col.ColumnName)
		} else {
			code += fmt.Sprintf("\n\t\t\t,%s", col.ColumnName)
		}
	}
	code += "`"
	code += ",\n"

	for _, col := range table.Columns {
		if col.IsUpdAble {
			code += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			code += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}

	code += "\t).Scan(\n"
	for _, col := range table.Columns {
		code += fmt.Sprintf("\t\t&ret.%s,\n", SnakeToPascal(col.ColumnName))
	}
	code += "\t)\n\n\treturn ret, err\n}\n" 

	return code
}

func (serv *SourceGeneratorGolang) generateDaoFileCodeDelete(table *dto.Table) string {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Delete(%s *entity.%s) error {\n", tnp, tni, tnp) +
		"\t_, err := rep.db.Exec(\n"

	code += fmt.Sprintf("\t\t`DELETE FROM %s\n\t\t WHERE ", tn)

	bindCount := 0
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if bindCount == 1 {
				code += fmt.Sprintf("%s = $1", col.ColumnName)
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = $%d", col.ColumnName, bindCount)
			}
		}
	}
	code += "`"
	code += ",\n"

	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			code += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}
	code += "\t)\n\n\treturn err\n}\n" 

	return code
}

func (serv *SourceGeneratorGolang) generateWeb() error {
	source := "_originalcopy_/golang/web"
	destination := serv.path + "web/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}

	if err := serv.generateStatic(); err != nil {
		return err
	}
	if err := serv.generateTemplate(); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return nil
}

func (serv *SourceGeneratorGolang) generateStatic() error {
	if err := serv.generateCss(); err != nil {
		return err
	}

	if err := serv.generateJs(); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return nil
}

func (serv *SourceGeneratorGolang) generateCss() error {
	return nil
}

func (serv *SourceGeneratorGolang) generateJs() error {
	path := serv.path + "web/static/js/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateJsFiles(path)
}

func (serv *SourceGeneratorGolang) generateJsFiles(path string) error {
	for _, table := range *serv.tables {
		code := GenerateJsCode(&table)
		if err := WriteFile(fmt.Sprintf("%s%s.js", path, table.TableName), code); err != nil {
			logger.LogError(err.Error())
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateTemplate() error {
	path := serv.path + "web/template/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateTemplateFiles(path)
}

func (serv *SourceGeneratorGolang) generateTemplateFiles(path string) error {
	if err := serv.generateTemplateFileHeader(path); err != nil {
		return err
	}
	if err := serv.generateTemplateFileFooter(path); err != nil {
		return err
	}
	if err := serv.generateTemplateFileIndex(path); err != nil {
		return err
	}
	for _, table := range *serv.tables {
		if err := serv.generateTemplateFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

func (serv *SourceGeneratorGolang) generateTemplateFileHeader(path string) error {
	content := GenerateHtmlCodeHeader(serv.tables)
	code := `{{define "header"}}` + content + `{{end}}`

	err := WriteFile(path + "_header.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateTemplateFileFooter(path string) error {
	content := GenerateHtmlCodeFooter()
	code := `{{define "footer"}}` + content + `{{end}}`

	err := WriteFile(path + "_footer.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateTemplateFileIndex(path string) error {
	content := "\n"
	code := `{{template "header" .}}` + content + `{{template "footer" .}}`

	err := WriteFile(path + "index.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateTemplateFile(table *dto.Table, path string) error {
	content := GenerateHtmlCodeMain(table)
	code := `{{template "header" .}}` + content + `{{template "footer" .}}`

	err := WriteFile(fmt.Sprintf("%s%s.html", path, table.TableName), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

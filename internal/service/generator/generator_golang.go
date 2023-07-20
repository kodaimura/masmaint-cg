package generator

import (
	"os"
	"fmt"
	"strings"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
	"masmaint-cg/internal/shared/constant"
)


type sourceGeneratorGolang struct {
	tables *[]dto.Table
	rdbms string
	path string
}

func NewSourceGeneratorGolang(tables *[]dto.Table, rdbms, path string) *sourceGeneratorGolang {
	return &sourceGeneratorGolang{
		tables, rdbms, path,
	}
}

// Goソース生成
func (serv *sourceGeneratorGolang) GenerateSource() error {

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

// cmd生成
func (serv *sourceGeneratorGolang) generateCmd() error {
	source := "_originalcopy_/golang/cmd"
	destination := serv.path + "cmd/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// config生成
func (serv *sourceGeneratorGolang) generateConfig() error {
	source := "_originalcopy_/golang/config"
	destination := serv.path + "config/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// core生成
func (serv *sourceGeneratorGolang) generateCore() error {
	source := "_originalcopy_/golang/core"
	destination := serv.path + "core/"
	
	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}
	return serv.generateDb()
}

// db生成
func (serv *sourceGeneratorGolang) generateDb() error {
	rdbmsCls := "postgresql"
	if serv.rdbms == constant.MYSQL || serv.rdbms == constant.MYSQL_8021 {
		rdbmsCls = "mysql"
	} else if serv.rdbms == constant.SQLITE_3350 {
		rdbmsCls = "sqlite3"
	}

	source := fmt.Sprintf("_originalcopy_/golang/core-sub/db/%s.go", rdbmsCls)
	destination := serv.path + fmt.Sprintf("core/db/db.go")

	err := CopyFile(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// log生成
func (serv *sourceGeneratorGolang) generateLog() error {
	source := "_originalcopy_/golang/log"
	destination := serv.path + "log/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// controller生成
func (serv *sourceGeneratorGolang) generateController() error {
	path := serv.path + "controller/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateControllerFiles(path)
}


func (serv *sourceGeneratorGolang) generateControllerFiles(path string) error {
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

const GO_CONTROLLER_ROUTER_FORMAT =
`
package controller

import (
	"github.com/gin-gonic/gin"

	"masmaint/core/auth"
)


func SetRouter(r *gin.Engine) {

	rm := r.Group("/mastertables", auth.NoopAuthMiddleware())
	{
		rm.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{})
		})
		%s
	}
}
`
func (serv *sourceGeneratorGolang) generateControllerFileRouter(path string) error {
	code := ""
	for _, table := range *serv.tables {
		tn := table.TableName
		tnc := SnakeToCamel(tn)
		tnp := SnakeToPascal(tn)

		code += fmt.Sprintf("\n\t\t%sController := New%sController()", tnc, tnp) +
			fmt.Sprintf("\n\t\trm.GET(\"/%s\", %sController.Get%sPage)", tn, tnc, tnp) +
			fmt.Sprintf("\n\t\trm.GET(\"/api/%s\", %sController.Get%s)", tn, tnc, tnp) +
			fmt.Sprintf("\n\t\trm.POST(\"/api/%s\", %sController.Post%s)", tn, tnc, tnp) +
			fmt.Sprintf("\n\t\trm.PUT(\"/api/%s\", %sController.Put%s)", tn, tnc, tnp) +
			fmt.Sprintf("\n\t\trm.DELETE(\"/api/%s\", %sController.Delete%s)\n", tn, tnc, tnp)
	}

	code = fmt.Sprintf(GO_CONTROLLER_ROUTER_FORMAT, code)

	err := WriteFile(path + "router.go", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) generateControllerFile(table *dto.Table, path string) error {
	code := "package controller\n\nimport (\n" +
		"\t\"github.com/gin-gonic/gin\"\n\n\t\"masmaint/service\"\n\t\"masmaint/dto\"\n)\n\n\n"

	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code += fmt.Sprintf("type %sService interface {\n", tnp) +
		fmt.Sprintf("\tGetAll() ([]dto.%sDto, error)\n", tnp) +
		fmt.Sprintf("\tCreate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%sDto *dto.%sDto) error\n", tni, tnp) +
		"}\n\n"

	code += fmt.Sprintf("type %sController struct {\n", tnc) +
		fmt.Sprintf("\t%sServ %sService\n", tni, tnp) +
		"}\n\n"

	code += fmt.Sprintf("func New%sController() *%sController {\n", tnp, tnc) +
		fmt.Sprintf("\t%sServ := service.New%sService()\n", tni, tnp) +
		fmt.Sprintf("\treturn &%sController{%sServ}\n", tnc, tni) +
		"}\n\n\n"

	code += serv.generateControllerFileCodeGetPage(table) + "\n\n"
	code += serv.generateControllerFileCodeGet(table) + "\n\n"
	code += serv.generateControllerFileCodePost(table) + "\n\n"
	code += serv.generateControllerFileCodePut(table) + "\n\n"
	code += serv.generateControllerFileCodeDelete(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) generateControllerFileCodeGetPage(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	return fmt.Sprintf("//GET /%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Get%sPage(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tc.HTML(200, \"%s.html\", gin.H{})\n", tn) +
		"}\n"
}


func (serv *sourceGeneratorGolang) generateControllerFileCodeGet(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//GET /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Get%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tret, err := ctr.%sServ.GetAll()\n\n", tni) +
		"\tif err != nil {\n\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n"
}


func (serv *sourceGeneratorGolang) generateControllerFileCodePost(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//POST /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Post%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) +
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tret, err := ctr.%sServ.Create(&%sDto)\n\n", tni, tni) +
		"\tif err != nil {\n\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n"
}


func (serv *sourceGeneratorGolang) generateControllerFileCodePut(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//PUT /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Put%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) +
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tret, err := ctr.%sServ.Update(&%sDto)\n\n", tni, tni) +
		"\tif err != nil {\n\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n"
}


func (serv *sourceGeneratorGolang) generateControllerFileCodeDelete(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//DELETE /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Delete%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) +
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tif err := ctr.%sServ.Delete(&%sDto); err != nil {\n", tni, tni) +
		"\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, gin.H{})\n}\n"
}

// dto生成
func (serv *sourceGeneratorGolang) generateDto() error {
	path := serv.path + "dto/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateDtoFiles(path)
}


func (serv *sourceGeneratorGolang) generateDtoFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateDtoFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}


func (serv *sourceGeneratorGolang) generateDtoFile(table *dto.Table, path string) error {
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

// service生成
func (serv *sourceGeneratorGolang) generateService() error {
	path := serv.path + "service/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateServiceFiles(path)
}


func (serv *sourceGeneratorGolang) generateServiceFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateServiceFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}


func (serv *sourceGeneratorGolang) generateServiceFile(table *dto.Table, path string) error {
	code := "package service\n\nimport (\n" +
		"\t\"errors\"\n\n\t\"masmaint/core/logger\"\n\t\"masmaint/model/entity\"\n" +
		"\t\"masmaint/model/dao\"\n\t\"masmaint/dto\"\n)\n\n\n"

	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code += fmt.Sprintf("type %sDao interface {\n", tnp) +
		fmt.Sprintf("\tSelectAll() ([]entity.%s, error)\n", tnp) +
		fmt.Sprintf("\tSelect(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tInsert(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%s *entity.%s) error\n", tni, tnp) +
		"}\n\n"

	code += fmt.Sprintf("type %sService struct {\n", tnc) +
		fmt.Sprintf("\t%sDao %sDao\n", tni, tnp) +
		"}\n\n"

	code += fmt.Sprintf("func New%sService() *%sService {\n", tnp, tnc) +
		fmt.Sprintf("\t%sDao := dao.New%sDao()\n", tni, tnp) +
		fmt.Sprintf("\treturn &%sService{%sDao}\n", tnc, tni) +
		"}\n\n\n"

	code += serv.generateServiceFileCodeGetAll(table) + "\n\n"
	code += serv.generateServiceFileCodeGetOne(table) + "\n\n"
	code += serv.generateServiceFileCodeCreate(table) + "\n\n"
	code += serv.generateServiceFileCodeUpdate(table) + "\n\n"
	code += serv.generateServiceFileCodeDelete(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) generateServiceFileCodeGetAll(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	return fmt.Sprintf("func (serv *%sService) GetAll() ([]dto.%sDto, error) {\n", tnc, tnp) +
		fmt.Sprintf("\trows, err := serv.%sDao.SelectAll()\n", tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		fmt.Sprintf("\t\treturn []dto.%sDto{}, errors.New(\"取得に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\tvar ret []dto.%sDto\n", tnp) +
		fmt.Sprintf("\tfor _, row := range rows {\n\t\tret = append(ret, row.To%sDto())\n\t}\n\n", tnp) +
		"\treturn ret, nil\n}\n"
}


func (serv *sourceGeneratorGolang) generateServiceFileCodeGetOne(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (serv *%sService) GetOne(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnc, tni, tnp, tnp) +
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
		"}\n"

	return code
}


func (serv *sourceGeneratorGolang) generateServiceFileCodeCreate(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (serv *%sService) Create(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnc, tni, tnp, tnp) +
		fmt.Sprintf("\tvar %s *entity.%s = entity.New%s()\n\n", tni, tnp, tnp)
	
	isFirst := true
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
		"}\n"

	return code
}


func (serv *sourceGeneratorGolang) generateServiceFileCodeUpdate(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (serv *%sService) Update(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnc, tni, tnp, tnp) +
		fmt.Sprintf("\tvar %s *entity.%s = entity.New%s()\n\n", tni, tnp, tnp)
	
	isFirst := true
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
		"}\n"

	return code
}


func (serv *sourceGeneratorGolang) generateServiceFileCodeDelete(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (serv *%sService) Delete(%sDto *dto.%sDto) error {\n", tnc, tni, tnp) +
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
	code += "\t\treturn errors.New(\"不正な値があります。\")\n\t}\n\n"
	code += fmt.Sprintf("\terr := serv.%sDao.Delete(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.LogError(err.Error())\n" +
		"\t\treturn errors.New(\"削除に失敗しました。\")\n\t}\n\n" +
		"\treturn nil\n}\n"

	return code
}

// model生成
func (serv *sourceGeneratorGolang) generateModel() error {
	if err := serv.generateEntity(); err != nil {
		return err
	}

	if err := serv.generateDao(); err != nil {
		return err
	}

	return nil
}

// entity生成
func (serv *sourceGeneratorGolang) generateEntity() error {
	path := serv.path + "model/entity/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateEntityFiles(path)
}


func (serv *sourceGeneratorGolang) generateEntityFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateEntityFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}


func (serv *sourceGeneratorGolang) getEntityFieldType(col *dto.Column) string {
	isNotNull := col.IsNotNull
	isPrimaryKey := col.IsPrimaryKey
	colType := col.ColumnType

	if colType == "s" || colType == "t" {
		if isNotNull || isPrimaryKey {
			return "string"
		}
		return "sql.NullString"
	}
	if colType == "i" {
		if isNotNull || isPrimaryKey {
			return "int64"
		}
		return "sql.NullInt64"
	}
	if colType == "f" {
		if isNotNull || isPrimaryKey {
			return "float64"
		}
		return "sql.NullFloat64"
	}
	return ""
}


func (serv *sourceGeneratorGolang) generateEntityFile(table *dto.Table, path string) error {
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
	code += serv.generateEntityFileCodeSetters(table) + "\n"
	code += serv.generateEntityFileCodeToDto(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) generateEntityFileCodeSetters(table *dto.Table) string {
	tnp := SnakeToPascal(table.TableName)

	code := fmt.Sprintf("func New%s() *%s {\n\treturn &%s{}\n}\n\n", tnp, tnp, tnp)
	for _, col := range table.Columns {
		code += serv.generateEntityFileCodeSetter(table, &col)
	}
	
	return code
}


func (serv *sourceGeneratorGolang) generateEntityFileCodeSetter(table *dto.Table, col *dto.Column) string {
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


func (serv *sourceGeneratorGolang) generateEntityFileCodeToDto(table *dto.Table) string {
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

// dao生成
func (serv *sourceGeneratorGolang) generateDao() error {
	path := serv.path + "model/dao/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateDaoFiles(path)
}


func (serv *sourceGeneratorGolang) generateDaoFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateDaoFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}


func (serv *sourceGeneratorGolang) generateDaoFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	code := "package dao\n\nimport (\n" +
		"\t\"database/sql\"\n\n\t\"masmaint/core/db\"\n\t\"masmaint/model/entity\"\n)\n\n\n"

	code += fmt.Sprintf("type %sDao struct {\n\tdb *sql.DB\n}\n\n", tnc) +
		fmt.Sprintf("func New%sDao() *%sDao {\n", tnp, tnc) +
		fmt.Sprintf("\tdb := db.GetDB()\n\treturn &%sDao{db}\n}\n\n\n", tnc)

	code += serv.generateDaoFileCodeSelectAll(table) + "\n\n"
	code += serv.generateDaoFileCodeSelect(table) + "\n\n"
	if serv.rdbms == constant.MYSQL {
		// RETURNING が使えない場合
		code += serv.generateDaoFileCodeInsert_MySQL(table) + "\n\n"
		code += serv.generateDaoFileCodeUpdate_MySQL(table) + "\n\n"
	} else {
		code += serv.generateDaoFileCodeInsert(table) + "\n\n"
		code += serv.generateDaoFileCodeUpdate(table) + "\n\n"
	}
	code += serv.generateDaoFileCodeDelete(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) concatPrimaryKeyWithCommas(cols []dto.Column) string {
	var ls []string
	for _, col := range cols {
		if col.IsPrimaryKey {
			ls = append(ls, col.ColumnName)
		}
	}
	return strings.Join(ls, ", ")
}


func (serv *sourceGeneratorGolang) getBindVariable(n int) string {
	if serv.rdbms == constant.POSTGRESQL {
		return fmt.Sprintf("$%d", n)
	} else {
		return "?"
	}
}


func (serv *sourceGeneratorGolang) generateDaoFileCodeSelectAll(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) SelectAll() ([]entity.%s, error) {\n", tnc, tnp) +
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
	code += "`,\n\t)\n\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n"

	code += fmt.Sprintf("\tfor rows.Next() {\n\t\t%s := entity.%s{}\n\t\terr = rows.Scan(\n", tni, tnp)
	for _, col := range table.Columns {
		cnp := SnakeToPascal(col.ColumnName)
		code += fmt.Sprintf("\t\t\t&%s.%s,\n", tni, cnp)
	}
	code += fmt.Sprintf("\t\t)\n\t\tif err != nil {\n\t\t\tbreak\n\t\t}\n\t\tret = append(ret, %s)\n\t}\n\n", tni)

	code += "\treturn ret, err\n}\n"
	return code
}


func (serv *sourceGeneratorGolang) generateDaoFileCodeSelect(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Select(%s *entity.%s) (entity.%s, error) {\n", tnc, tni, tnp, tnp) +
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
				code += fmt.Sprintf("%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			}
		}
	}
	code += "`,\n"

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


func (serv *sourceGeneratorGolang) concatBindVariableWithCommas(bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, serv.getBindVariable(i))
	}
	return strings.Join(ls, ",")
}


func (serv *sourceGeneratorGolang) generateDaoFileCodeInsert(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Insert(%s *entity.%s) (entity.%s, error) {\n", tnc, tni, tnp, tnp) +
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
	code += "`,\n"

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


func (serv *sourceGeneratorGolang) generateDaoFileCodeUpdate(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Update(%s *entity.%s) (entity.%s, error) {\n", tnc, tni, tnp, tnp) +
		fmt.Sprintf("\tvar ret entity.%s\n\n\terr := rep.db.QueryRow(\n", tnp)

	code += fmt.Sprintf("\t\t`UPDATE %s\n\t\t SET\n", tn)
	bindCount := 0
	for _, col := range table.Columns {
		if col.IsUpdAble {
			bindCount += 1
			if bindCount == 1 {
				code += fmt.Sprintf("\t\t\t%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t\t,%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			}
		}
	}
	code += "\n\t\t WHERE "

	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if isFirst {
				code += fmt.Sprintf("%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, serv.getBindVariable(bindCount))
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
	code += "`,\n"

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


func (serv *sourceGeneratorGolang) getAutoIncrementColumn(table *dto.Table) (dto.Column, bool) {
	for _, col := range table.Columns {
		if col.IsPrimaryKey && !col.IsInsAble && (col.ColumnType == "i") {
			return col, true
		}
	}
	return dto.Column{}, false
}


func (serv *sourceGeneratorGolang) generateDaoFileCodeInsert_MySQL(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	aicol, hasAICol := serv.getAutoIncrementColumn(table)
	code := fmt.Sprintf("func (rep *%sDao) Insert(%s *entity.%s) (entity.%s, error) {\n", tnc, tni, tnp, tnp)

	if hasAICol {
		code += "\tresult, err := rep.db.Exec(\n"
	} else {
		code += "\t_, err := rep.db.Exec(\n"
	}
	
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
	code += fmt.Sprintf("\n\t\t ) VALUES (%s)`,\n", serv.concatBindVariableWithCommas(bindCount))

	for _, col := range table.Columns {
		if col.IsInsAble {
			code += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}
	code += fmt.Sprintf("\t)\n\n\tif err != nil {\n\t\treturn entity.%s{}, err\n\t}\n\n", tnp)

	if hasAICol {
		code += "\tlastInsertId, err := result.LastInsertId()\n" + 
			fmt.Sprintf("\tif err != nil {\n\t\treturn entity.%s{}, err\n\t}\n\n", tnp) +
			fmt.Sprintf("\t%s.Set%s(lastInsertId)\n\n", tni, SnakeToPascal(aicol.ColumnName))
	}

	code += fmt.Sprintf("\treturn rep.Select(%s)\n}\n", tni)
	return code
}


func (serv *sourceGeneratorGolang) generateDaoFileCodeUpdate_MySQL(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Update(%s *entity.%s) (entity.%s, error) {\n", tnc, tni, tnp, tnp)
	code += "\t_, err := rep.db.Exec(\n"

	code += fmt.Sprintf("\t\t`UPDATE %s\n\t\t SET\n", tn)
	bindCount := 0
	for _, col := range table.Columns {
		if col.IsUpdAble {
			bindCount += 1
			if bindCount == 1 {
				code += fmt.Sprintf("\t\t\t%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t\t,%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			}
		}
	}
	code += "\n\t\t WHERE "

	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if isFirst {
				code += fmt.Sprintf("%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			}
		}
	}
	code += "`,\n"

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

	code += fmt.Sprintf("\t)\n\n\tif err != nil {\n\t\treturn entity.%s{}, err\n\t}\n\n", tnp)
	code += fmt.Sprintf("\treturn rep.Select(%s)\n}\n", tni)

	return code
}


func (serv *sourceGeneratorGolang) generateDaoFileCodeDelete(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	code := fmt.Sprintf("func (rep *%sDao) Delete(%s *entity.%s) error {\n", tnc, tni, tnp) +
		"\t_, err := rep.db.Exec(\n"

	code += fmt.Sprintf("\t\t`DELETE FROM %s\n\t\t WHERE ", tn)

	bindCount := 0
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if bindCount == 1 {
				code += fmt.Sprintf("%s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, serv.getBindVariable(bindCount))
			}
		}
	}
	code += "`,\n"

	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			code += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}
	code += "\t)\n\n\treturn err\n}\n" 

	return code
}

// web生成
func (serv *sourceGeneratorGolang) generateWeb() error {
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
		return err
	}

	return nil
}

// static生成
func (serv *sourceGeneratorGolang) generateStatic() error {
	if err := serv.generateCss(); err != nil {
		return err
	}

	if err := serv.generateJs(); err != nil {
		return err
	}

	return nil
}

// css生成
func (serv *sourceGeneratorGolang) generateCss() error {
	return nil
}

// js生成
func (serv *sourceGeneratorGolang) generateJs() error {
	path := serv.path + "web/static/js/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateJsFiles(path)
}


func (serv *sourceGeneratorGolang) generateJsFiles(path string) error {
	for _, table := range *serv.tables {
		code := GenerateJsCode(&table)
		if err := WriteFile(fmt.Sprintf("%s%s.js", path, table.TableName), code); err != nil {
			logger.LogError(err.Error())
			return err
		}
	}
	return nil
}

// template生成
func (serv *sourceGeneratorGolang) generateTemplate() error {
	path := serv.path + "web/template/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateTemplateFiles(path)
}


func (serv *sourceGeneratorGolang) generateTemplateFiles(path string) error {
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


func (serv *sourceGeneratorGolang) generateTemplateFileHeader(path string) error {
	content := GenerateHtmlCodeHeader(serv.tables)
	code := `{{define "header"}}` + content + `{{end}}`

	err := WriteFile(path + "_header.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) generateTemplateFileFooter(path string) error {
	content := GenerateHtmlCodeFooter()
	code := `{{define "footer"}}` + content + `{{end}}`

	err := WriteFile(path + "_footer.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) generateTemplateFileIndex(path string) error {
	content := "\n"
	code := `{{template "header" .}}` + content + `{{template "footer" .}}`

	err := WriteFile(path + "index.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}


func (serv *sourceGeneratorGolang) generateTemplateFile(table *dto.Table, path string) error {
	content := GenerateHtmlCodeMain(table)
	code := `{{template "header" .}}` + content + `{{template "footer" .}}`

	err := WriteFile(fmt.Sprintf("%s%s.html", path, table.TableName), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

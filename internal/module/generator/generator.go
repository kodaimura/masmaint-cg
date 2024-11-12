package generator

import (
	"os"
	"fmt"
	"strings"
	"github.com/kodaimura/ddlparse"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
	"masmaint-cg/internal/shared/constant"
)


type generator struct {
	tables []ddlparse.Table
	rdbms string
	path string
}

type Generator interface {
	Generate() (string, error)
}

func NewGenerator([]ddlparse.Table, rdbms, path string) Generator {
	return &generator{
		tables, rdbms, path,
	}
}

// Goソース生成
func (gen *generator) Generate() error {
	if err := os.MkdirAll(gen.path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	if err := gen.generateGitignore(); err != nil {
		return err
	}
	if err := gen.generateCmd(); err != nil {
		return err
	}
	if err := gen.generateConfig(); err != nil {
		return err
	}
	if err := gen.generateCore(); err != nil {
		return err
	}
	if err := gen.generateLog(); err != nil {
		return err
	}
	if err := gen.generateController(); err != nil {
		return err
	}
	if err := gen.generateDto(); err != nil {
		return err
	}
	if err := gen.generategenice(); err != nil {
		return err
	}
	if err := gen.generateModel(); err != nil {
		return err
	}
	if err := gen.generateWeb(); err != nil {
		return err
	}

	return nil	
}

// .gitignore生成
func (gen *generator) generateGitignore() error {
	code := "*.log\n*.db\n*.sqlite3\n.DS_Store\nmain\n.env\nlocal.env"
	err := WriteFile(fmt.Sprintf("%s.gitignore", gen.path), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// cmd生成
func (gen *generator) generateCmd() error {
	source := "_originalcopy_/golang/cmd"
	destination := gen.path + "cmd/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// config生成
func (gen *generator) generateConfig() error {
	path := gen.path + "config/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}
	source := "_originalcopy_/golang/config"
	destination := gen.path + "config/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.Error(err.Error())
	}
	return gen.generateEnv()
}

// env生成
func (gen *generator) generateEnv() error {
	path := gen.path + "config/env/"

	rdbmsCls := "postgresql"
	if gen.rdbms == constant.MYSQL {
		rdbmsCls = "mysql"
	} else if gen.rdbms == constant.SQLITE_3350 {
		rdbmsCls = "sqlite3"
	}

	source := fmt.Sprintf("_originalcopy_/golang/config-sub/env/local.%s.env", rdbmsCls)
	destination := fmt.Sprintf("%slocal.env", path)

	err := CopyFile(source, destination)
	if err != nil {
		logger.Error(err.Error())
	}

	return err
}

// core生成
func (gen *generator) generateCore() error {
	source := "_originalcopy_/golang/core"
	destination := gen.path + "core/"
	
	err := CopyDir(source, destination)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return gen.generateDb()
}

// db生成
func (gen *generator) generateDb() error {
	path := gen.path + "core/db/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	rdbmsCls := "postgresql"
	if gen.rdbms == constant.MYSQL {
		rdbmsCls = "mysql"
	} else if gen.rdbms == constant.SQLITE_3350 {
		rdbmsCls = "sqlite3"
	}

	source := fmt.Sprintf("_originalcopy_/golang/core-sub/db/%s.go", rdbmsCls)
	destination := fmt.Sprintf("%sdb.go", path)

	err := CopyFile(source, destination)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// log生成
func (gen *generator) generateLog() error {
	source := "_originalcopy_/golang/log"
	destination := gen.path + "log/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// controller生成
func (gen *generator) generateController() error {
	path := gen.path + "controller/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	return gen.generateControllerFiles(path)
}

// controller内のファイル生成
func (gen *generator) generateControllerFiles(path string) error {
	if err := gen.generateControllerFileRouter(path); err != nil {
		return err
	}

	for _, table := range *gen.tables {
		if err := gen.generateControllerFile(&table, path); err != nil {
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

// controller/router.go生成
func (gen *generator) generateControllerFileRouter(path string) error {
	code := ""
	for _, table := range *gen.tables {
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
		logger.Error(err.Error())
	}
	return err
}

// controller/*.go生成
func (gen *generator) generateControllerFile(table *dto.Table, path string) error {
	code := "package controller\n\nimport (\n\t\"github.com/gin-gonic/gin\"\n\n" +
		"\tcerror \"masmaint/core/error\"\n\t\"masmaint/genice\"\n\t\"masmaint/dto\"\n)\n\n\n"

	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code += gen.generateCodegeniceInterface(table) + "\n"

	code += fmt.Sprintf("type %sController struct {\n", tnc) +
		fmt.Sprintf("\t%sgen %sgenice\n", tni, tnp) +
		"}\n\n"

	code += fmt.Sprintf("func New%sController() *%sController {\n", tnp, tnc) +
		fmt.Sprintf("\t%sgen := genice.New%sgenice()\n", tni, tnp) +
		fmt.Sprintf("\treturn &%sController{%sgen}\n", tnc, tni) +
		"}\n\n\n"

	code += gen.generateControllerFileCodeGetPage(table) + "\n\n"
	code += gen.generateControllerFileCodeGet(table) + "\n\n"
	code += gen.generateControllerFileCodePost(table) + "\n\n"
	code += gen.generateControllerFileCodePut(table) + "\n\n"
	code += gen.generateControllerFileCodeDelete(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// genice interfaceプログラム生成
func (gen *generator) generateCodegeniceInterface(table *dto.Table) string {
	tnp := SnakeToPascal(table.TableName)
	tni := GetSnakeInitial(table.TableName)
	return fmt.Sprintf("type %sgenice interface {\n", tnp) +
		fmt.Sprintf("\tGetAll() ([]dto.%sDto, error)\n", tnp) +
		fmt.Sprintf("\tCreate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%sDto *dto.%sDto) (dto.%sDto, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%sDto *dto.%sDto) error\n", tni, tnp) +
		"}\n"
}

// controllerのGetPageメソッドプログラム生成
func (gen *generator) generateControllerFileCodeGetPage(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	return fmt.Sprintf("//GET /%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Get%sPage(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tc.HTML(200, \"%s.html\", gin.H{})\n", tn) +
		"}\n"
}

// controllerのGetメソッドプログラム生成
func (gen *generator) generateControllerFileCodeGet(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//GET /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Get%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tret, err := ctr.%sgen.GetAll()\n\n", tni) +
		"\tif err != nil {\n\t\tc.JSON(500, gin.H{})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n"
}

// controllerのPostメソッドプログラム生成
func (gen *generator) generateControllerFileCodePost(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//POST /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Post%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) +
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tret, err := ctr.%sgen.Create(&%sDto)\n\n", tni, tni) +
		"\tif err != nil {\n\t\tif _, ok := err.(*cerror.InvalidArgumentError); ok {\n" +
		"\t\t\tc.JSON(400, gin.H{})\n\t\t} else {\n\t\t\tc.JSON(500, gin.H{})\n\t\t}" +
		"\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n"
}

// controllerのPutメソッドプログラム生成
func (gen *generator) generateControllerFileCodePut(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//PUT /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Put%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) +
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tret, err := ctr.%sgen.Update(&%sDto)\n\n", tni, tni) +
		"\tif err != nil {\n\t\tif _, ok := err.(*cerror.InvalidArgumentError); ok {\n" +
		"\t\t\tc.JSON(400, gin.H{})\n\t\t} else {\n\t\t\tc.JSON(500, gin.H{})\n\t\t}" +
		"\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, ret)\n}\n"
}

// controllerのDeleteメソッドプログラム生成
func (gen *generator) generateControllerFileCodeDelete(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	return fmt.Sprintf("//DELETE /api/%s\n", tn) +
		fmt.Sprintf("func (ctr *%sController) Delete%s(c *gin.Context) {\n", tnc, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp) +
		fmt.Sprintf("\tif err := c.ShouldBindJSON(&%sDto); err != nil {\n", tni) +
		"\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		fmt.Sprintf("\tif err := ctr.%sgen.Delete(&%sDto); err != nil {\n", tni, tni) +
		"\t\tif _, ok := err.(*cerror.InvalidArgumentError); ok {\n" +
		"\t\t\tc.JSON(400, gin.H{})\n\t\t} else {\n\t\t\tc.JSON(500, gin.H{})\n\t\t}" +
		"\n\t\tc.Abort()\n\t\treturn\n\t}\n\n" +
		"\tc.JSON(200, gin.H{})\n}\n"
}

// dto生成
func (gen *generator) generateDto() error {
	path := gen.path + "dto/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	return gen.generateDtoFiles(path)
}

// dto内のファイル生成
func (gen *generator) generateDtoFiles(path string) error {
	for _, table := range *gen.tables {
		if err := gen.generateDtoFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// dto/*.go生成
func (gen *generator) generateDtoFile(table *dto.Table, path string) error {
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
		logger.Error(err.Error())
	}
	return err
}

// genice生成
func (gen *generator) generategenice() error {
	path := gen.path + "genice/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	return gen.generategeniceFiles(path)
}

// genice内のファイル生成
func (gen *generator) generategeniceFiles(path string) error {
	for _, table := range *gen.tables {
		if err := gen.generategeniceFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// genice/*.go生成
func (gen *generator) generategeniceFile(table *dto.Table, path string) error {
	code := "package genice\n\nimport (\n" +
		"\tcerror \"masmaint/core/error\"\n\n\t\"masmaint/core/logger\"\n\t\"masmaint/model/entity\"\n" +
		"\t\"masmaint/model/dao\"\n\t\"masmaint/dto\"\n)\n\n\n"

	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code += gen.generateCodeDaoInterface(table) + "\n"

	code += fmt.Sprintf("type %sgenice struct {\n", tnc) +
		fmt.Sprintf("\t%sDao %sDao\n", tni, tnp) +
		"}\n\n"

	code += fmt.Sprintf("func New%sgenice() *%sgenice {\n", tnp, tnc) +
		fmt.Sprintf("\t%sDao := dao.New%sDao()\n", tni, tnp) +
		fmt.Sprintf("\treturn &%sgenice{%sDao}\n", tnc, tni) +
		"}\n\n\n"

	code += gen.generategeniceFileCodeGetAll(table) + "\n\n"
	code += gen.generategeniceFileCodeGetOne(table) + "\n\n"
	code += gen.generategeniceFileCodeCreate(table) + "\n\n"
	code += gen.generategeniceFileCodeUpdate(table) + "\n\n"
	code += gen.generategeniceFileCodeDelete(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// dao interfaceプログラム生成
func (gen *generator) generateCodeDaoInterface(table *dto.Table) string {
	tnp := SnakeToPascal(table.TableName)
	tni := GetSnakeInitial(table.TableName)

	return fmt.Sprintf("type %sDao interface {\n", tnp) +
		fmt.Sprintf("\tSelectAll() ([]entity.%s, error)\n", tnp) +
		fmt.Sprintf("\tSelect(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tInsert(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tUpdate(%s *entity.%s) (entity.%s, error)\n", tni, tnp, tnp) +
		fmt.Sprintf("\tDelete(%s *entity.%s) error\n", tni, tnp) + "}\n"
}

// geniceのGetAllメソッドプログラム生成
func (gen *generator) generategeniceFileCodeGetAll(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	return fmt.Sprintf("func (gen *%sgenice) GetAll() ([]dto.%sDto, error) {\n", tnc, tnp) +
		fmt.Sprintf("\trows, err := gen.%sDao.SelectAll()\n", tni) +
		"\tif err != nil {\n\t\tlogger.Error(err.Error())\n" +
		fmt.Sprintf("\t\treturn []dto.%sDto{}, cerror.NewDaoError(\"取得に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\tvar ret []dto.%sDto\n", tnp) +
		fmt.Sprintf("\tfor _, row := range rows {\n\t\tret = append(ret, row.To%sDto())\n\t}\n\n", tnp) +
		"\treturn ret, nil\n}\n"
}

// geniceのGetOneメソッドプログラム生成
func (gen *generator) generategeniceFileCodeGetOne(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (gen *%sgenice) GetOne(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnc, tni, tnp, tnp) +
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
	code += fmt.Sprintf("\t\treturn dto.%sDto{}, cerror.NewInvalidArgumentError(\"不正な値があります。\")\n\t}\n\n", tnp)
	code += fmt.Sprintf("\trow, err := gen.%sDao.Select(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.Error(err.Error())\n" +
		fmt.Sprintf("\t\treturn dto.%sDto{}, cerror.NewDaoError(\"取得に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\treturn row.To%sDto(), nil\n", tnp) +
		"}\n"

	return code
}

// geniceのCreateメソッドプログラム生成
func (gen *generator) generategeniceFileCodeCreate(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (gen *%sgenice) Create(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnc, tni, tnp, tnp) +
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
	code += fmt.Sprintf("\t\treturn dto.%sDto{}, cerror.NewInvalidArgumentError(\"不正な値があります。\")\n\t}\n\n", tnp)

	code += fmt.Sprintf("\trow, err := gen.%sDao.Insert(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.Error(err.Error())\n" +
		fmt.Sprintf("\t\treturn dto.%sDto{}, cerror.NewDaoError(\"登録に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\treturn row.To%sDto(), nil\n", tnp) +
		"}\n"

	return code
}

// geniceのUpdateメソッドプログラム生成
func (gen *generator) generategeniceFileCodeUpdate(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (gen *%sgenice) Update(%sDto *dto.%sDto) (dto.%sDto, error) {\n", tnc, tni, tnp, tnp) +
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
	code += fmt.Sprintf("\t\treturn dto.%sDto{}, cerror.NewInvalidArgumentError(\"不正な値があります。\")\n\t}\n\n", tnp)
	code += fmt.Sprintf("\trow, err := gen.%sDao.Update(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.Error(err.Error())\n" +
		fmt.Sprintf("\t\treturn dto.%sDto{}, cerror.NewDaoError(\"更新に失敗しました。\")\n\t}\n\n", tnp) +
		fmt.Sprintf("\treturn row.To%sDto(), nil\n", tnp) +
		"}\n"

	return code
}

// geniceのDeleteメソッドプログラム生成
func (gen *generator) generategeniceFileCodeDelete(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	code := fmt.Sprintf("func (gen *%sgenice) Delete(%sDto *dto.%sDto) error {\n", tnc, tni, tnp) +
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
	code += "\t\treturn cerror.NewInvalidArgumentError(\"不正な値があります。\")\n\t}\n\n"
	code += fmt.Sprintf("\terr := gen.%sDao.Delete(%s)\n", tni, tni) +
		"\tif err != nil {\n\t\tlogger.Error(err.Error())\n" +
		"\t\treturn cerror.NewDaoError(\"削除に失敗しました。\")\n\t}\n\n" +
		"\treturn nil\n}\n"

	return code
}

// model生成
func (gen *generator) generateModel() error {
	if err := gen.generateEntity(); err != nil {
		return err
	}
	if err := gen.generateDao(); err != nil {
		return err
	}

	return nil
}

// entity生成
func (gen *generator) generateEntity() error {
	path := gen.path + "model/entity/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	return gen.generateEntityFiles(path)
}

// entity内のファイル生成
func (gen *generator) generateEntityFiles(path string) error {
	for _, table := range *gen.tables {
		if err := gen.generateEntityFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// CSVフォーマットのカラム型からentityフィールド用の型取得
func (gen *generator) getEntityFieldType(col *dto.Column) string {
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

// entity/*.go生成
func (gen *generator) generateEntityFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnp := SnakeToPascal(tn)
	code := "package entity\n\nimport (\n" +
		"\t\"database/sql\"\n\n\t\"masmaint/dto\"\n\t\"masmaint/core/utils\"\n)\n\n\n"

	code += fmt.Sprintf("type %s struct {\n", tnp)
	for _, col := range table.Columns {
		cn := col.ColumnName
		cnp := SnakeToPascal(cn)
		code += fmt.Sprintf("\t%s %s `db:\"%s\"`\n", cnp, gen.getEntityFieldType(&col), cn)
	}
	code += "}\n\n"
	code += gen.generateEntityFileCodeSetters(table) + "\n"
	code += gen.generateEntityFileCodeToDto(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// entityのセッタープログラム生成
func (gen *generator) generateEntityFileCodeSetters(table *dto.Table) string {
	tnp := SnakeToPascal(table.TableName)

	code := fmt.Sprintf("func New%s() *%s {\n\treturn &%s{}\n}\n\n", tnp, tnp, tnp)
	for _, col := range table.Columns {
		code += gen.generateEntityFileCodeSetter(table, &col)
	}
	
	return code
}

// entityのセッタープログラム生成
func (gen *generator) generateEntityFileCodeSetter(table *dto.Table, col *dto.Column) string {
	tnp := SnakeToPascal(table.TableName)
	colType := gen.getEntityFieldType(col)
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

// entityのToDtoメソッドプログラム生成
func (gen *generator) generateEntityFileCodeToDto(table *dto.Table) string {
	tnp := SnakeToPascal(table.TableName)
	tni := GetSnakeInitial(table.TableName)

	code := fmt.Sprintf("func (e *%s) To%sDto() dto.%sDto {\n", tnp, tnp, tnp) +
		fmt.Sprintf("\tvar %sDto dto.%sDto\n\n", tni, tnp)
	for _, col := range table.Columns {
		colType := gen.getEntityFieldType(&col)
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
func (gen *generator) generateDao() error {
	path := gen.path + "model/dao/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	return gen.generateDaoFiles(path)
}

// dao内のファイル生成
func (gen *generator) generateDaoFiles(path string) error {
	for _, table := range *gen.tables {
		if err := gen.generateDaoFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// dao/*.go生成
func (gen *generator) generateDaoFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	code := "package dao\n\nimport (\n" +
		"\t\"database/sql\"\n\n\t\"masmaint/core/db\"\n\t\"masmaint/model/entity\"\n)\n\n\n"

	code += fmt.Sprintf("type %sDao struct {\n\tdb *sql.DB\n}\n\n", tnc) +
		fmt.Sprintf("func New%sDao() *%sDao {\n", tnp, tnc) +
		fmt.Sprintf("\tdb := db.GetDB()\n\treturn &%sDao{db}\n}\n\n\n", tnc)

	code += gen.generateDaoFileCodeSelectAll(table) + "\n\n"
	code += gen.generateDaoFileCodeSelect(table) + "\n\n"
	if gen.rdbms == constant.MYSQL {
		// RETURNING が使えない場合
		code += gen.generateDaoFileCodeInsert_MySQL(table) + "\n\n"
		code += gen.generateDaoFileCodeUpdate_MySQL(table) + "\n\n"
	} else {
		code += gen.generateDaoFileCodeInsert(table) + "\n\n"
		code += gen.generateDaoFileCodeUpdate(table) + "\n\n"
	}
	code += gen.generateDaoFileCodeDelete(table)

	err := WriteFile(fmt.Sprintf("%s%s.go", path, tn), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// SELECT ORDER BYで指定するpkのカンマ区切り文字列生成
func (gen *generator) concatPrimaryKeyWithCommas(cols []dto.Column) string {
	var ls []string
	for _, col := range cols {
		if col.IsPrimaryKey {
			ls = append(ls, col.ColumnName)
		}
	}
	return strings.Join(ls, ", ")
}

// バインド変数文字列生成
func (gen *generator) getBindVariable(n int) string {
	if gen.rdbms == constant.POSTGRESQL {
		return fmt.Sprintf("$%d", n)
	} else {
		return "?"
	}
}

// daoのSelectAllメソッド生成
func (gen *generator) generateDaoFileCodeSelectAll(table *dto.Table) string {
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
	code += fmt.Sprintf("\n\t\t FROM %s\n\t\t ORDER BY %s ASC", tn, gen.concatPrimaryKeyWithCommas(table.Columns))
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

// daoのSelectメソッド生成
func (gen *generator) generateDaoFileCodeSelect(table *dto.Table) string {
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
				code += fmt.Sprintf("%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, gen.getBindVariable(bindCount))
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

// INSERT VALUES (...)用のバインド変数のカンマ区切り文字列生成
func (gen *generator) concatBindVariableWithCommas(bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, gen.getBindVariable(i))
	}
	return strings.Join(ls, ",")
}

// daoのInsertメソッド生成
func (gen *generator) generateDaoFileCodeInsert(table *dto.Table) string {
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
	code += fmt.Sprintf("\n\t\t ) VALUES (%s)\n\t\t RETURNING\n", gen.concatBindVariableWithCommas(bindCount))
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

// daoのUpdateメソッド生成
func (gen *generator) generateDaoFileCodeUpdate(table *dto.Table) string {
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
				code += fmt.Sprintf("\t\t\t%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t\t,%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
			}
		}
	}
	code += "\n\t\t WHERE "

	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if isFirst {
				code += fmt.Sprintf("%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, gen.getBindVariable(bindCount))
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

// daoのDeleteメソッド生成
func (gen *generator) generateDaoFileCodeDelete(table *dto.Table) string {
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
				code += fmt.Sprintf("%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, gen.getBindVariable(bindCount))
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

// AUTO_INCREMENTのカラム取得
// このシステムではPK・入力不可・整数型のカラムはAUTO_INCREMENTのカラムと判定する
func (gen *generator) getAutoIncrementColumn(table *dto.Table) (dto.Column, bool) {
	for _, col := range table.Columns {
		if col.IsPrimaryKey && !col.IsInsAble && (col.ColumnType == "i") {
			return col, true
		}
	}
	return dto.Column{}, false
}

// daoのInsertメソッド生成
func (gen *generator) generateDaoFileCodeInsert_MySQL(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	aicol, hasAICol := gen.getAutoIncrementColumn(table)
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
	code += fmt.Sprintf("\n\t\t ) VALUES (%s)`,\n", gen.concatBindVariableWithCommas(bindCount))

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

// daoのUpdateメソッド生成
func (gen *generator) generateDaoFileCodeUpdate_MySQL(table *dto.Table) string {
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
				code += fmt.Sprintf("\t\t\t%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
			} else {
				code += fmt.Sprintf("\n\t\t\t,%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
			}
		}
	}
	code += "\n\t\t WHERE "

	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			bindCount += 1
			if isFirst {
				code += fmt.Sprintf("%s = %s", col.ColumnName, gen.getBindVariable(bindCount))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t    AND %s = %s", col.ColumnName, gen.getBindVariable(bindCount))
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

// web生成
func (gen *generator) generateWeb() error {
	source := "_originalcopy_/golang/web"
	destination := gen.path + "web/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if err := gen.generateStatic(); err != nil {
		return err
	}
	if err := gen.generateTemplate(); err != nil {
		return err
	}

	return nil
}

// static生成
func (gen *generator) generateStatic() error {
	if err := gen.generateCss(); err != nil {
		return err
	}
	if err := gen.generateJs(); err != nil {
		return err
	}

	return nil
}

// css生成
func (gen *generator) generateCss() error {
	//path := gen.path + "public/static/css/"
	// _originalcopy_からコピー
	return nil
}

// js生成
func (gen *generator) generateJs() error {
	path := gen.path + "web/static/js/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	return gen.generateJsFiles(path)
}

// js/*.js生成
func (gen *generator) generateJsFiles(path string) error {
	for _, table := range *gen.tables {
		code := GenerateJsCode(&table)
		if err := WriteFile(fmt.Sprintf("%s%s.js", path, table.TableName), code); err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	return nil
}

// template生成
func (gen *generator) generateTemplate() error {
	path := gen.path + "web/template/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	return gen.generateTemplateFiles(path)
}

// template内のファイル生成
func (gen *generator) generateTemplateFiles(path string) error {
	if err := gen.generateTemplateFileHeader(path); err != nil {
		return err
	}
	if err := gen.generateTemplateFileFooter(path); err != nil {
		return err
	}
	if err := gen.generateTemplateFileIndex(path); err != nil {
		return err
	}

	for _, table := range *gen.tables {
		if err := gen.generateTemplateFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// template/_header.html生成
func (gen *generator) generateTemplateFileHeader(path string) error {
	content := GenerateHtmlCodeHeader(gen.tables)
	code := `{{define "header"}}` + content + `{{end}}`

	err := WriteFile(path + "_header.html", code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// template/_footer.html生成
func (gen *generator) generateTemplateFileFooter(path string) error {
	content := GenerateHtmlCodeFooter()
	code := `{{define "footer"}}` + content + `{{end}}`

	err := WriteFile(path + "_footer.html", code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// template/index.html生成
func (gen *generator) generateTemplateFileIndex(path string) error {
	content := "\n"
	code := `{{template "header" .}}` + content + `{{template "footer" .}}`

	err := WriteFile(path + "index.html", code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// template/*.html生成
func (gen *generator) generateTemplateFile(table *dto.Table, path string) error {
	content := GenerateHtmlCodeMain(table)
	code := `{{template "header" .}}` + content + `{{template "footer" .}}`

	err := WriteFile(fmt.Sprintf("%s%s.html", path, table.TableName), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// Jsコード生成
func GenerateJsCode(table *dto.Table) string {
	tn := table.TableName
	code := JS_COMMON_CODE + "\n"
	code += fmt.Sprintf(
		JS_FORMAT, 
		generateJsCode_createTrNew(table),
		generateJsCode_createTr(table),
		tn,
		generateJsCode_setUp(table),
		generateJsCode_doPutAll(table),
		generateJsCode_doPutAll_for(table),
		generateJsCode_doPutAll_if(table),
		generateJsCode_doPutAll_requestBody(table),
		tn,
		generateJsCode_doPutAll_then(table),
		generateJsCode_doPost(table),
		generateJsCode_doPost_if(table),
		generateJsCode_doPost_requestBody(table),
		tn,
		tn,
	)

	return code
}

// HTMLコード生成
func GenerateHtmlCodeMain(table *dto.Table) string {
	return fmt.Sprintf(
		HTML_FORMAT,
		generateHtmlCode_h2(table),
		"%", //要改善
		generateHtmlCode_tr(table),
		table.TableName,
	)
}


// HTMLコード生成（共通フッタ）
func GenerateHtmlCodeFooter() string {
	return HTML_FOOTER_CODE
}


// HTMLコード生成（共通ヘッダ）
func GenerateHtmlCodeHeader(tables *[]dto.Table) string {
	return fmt.Sprintf(
		HTML_HEADER_FORMAT,
		generateHtmlCodeHeader_ul(tables),
	)
}

// -------------------------------------------------------------------------------------------------//

func generateJsCode_createTrNew(table *dto.Table) string {
	code := "`<tr id='new'><td></td>`"
	for _, col := range table.Columns {
		if col.IsInsAble {
			code += fmt.Sprintf("\n\t\t+ `<td><input type='text' id='%s_new'></td>`", col.ColumnName)
		} else {
			code += "\n\t\t+ `<td><input type='text' disabled></td>`"
		}
	}
	return code + ";"
}

func generateJsCode_createTr(table *dto.Table) string {
	code := "`<tr><td><input class='form-check-input' type='checkbox' name='del' value='${JSON.stringify(elem)}'></td>`"
	for _, col := range table.Columns {
		cn := col.ColumnName
		if col.IsUpdAble {
			code += fmt.Sprintf("\n\t\t+ `<td><input type='text' name='%s' value='${nullToEmpty(elem.%s)}'><input type='hidden' name='%s_bk' value='${nullToEmpty(elem.%s)}'></td>`", cn, cn, cn, cn)
		} else {
			code += fmt.Sprintf("\n\t\t+ `<td><input type='text' name='%s' value='${nullToEmpty(elem.%s)}' disabled></td>`", cn, cn)
		}
	}
	return code + ";"
}

func generateJsCode_setUp(table *dto.Table) string {
	code := ""
	for _, col := range table.Columns {
		if col.IsUpdAble {
			code += fmt.Sprintf("\n\t\taddChangedAction('%s');", col.ColumnName)
		}
	}
	return code
}

func generateJsCode_doPutAll(table *dto.Table) string {
	code := ""
	for _, col := range table.Columns {
		cn := col.ColumnName
		code += fmt.Sprintf("\n\tlet %s = document.getElementsByName('%s');", cn, cn)
		if col.IsUpdAble {
			code += fmt.Sprintf("\n\tlet %s_bk = document.getElementsByName('%s_bk');", cn, cn)
		}
	}
	return code
}

func generateJsCode_doPutAll_for(table *dto.Table) string {
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			return col.ColumnName
		}
	}
	return table.Columns[0].ColumnName
}

func generateJsCode_doPutAll_if(table *dto.Table) string {
	code := ""
	isFirst := true
	for _, col := range table.Columns {
		if col.IsUpdAble {
			cn := col.ColumnName
			if isFirst {
				code += fmt.Sprintf("(%s[i].value !== %s_bk[i].value)", cn, cn)
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t\t|| (%s[i].value !== %s_bk[i].value)", cn, cn)
			}
		}
	}
	return code
}

func generateJsCode_doPutAll_requestBody(table *dto.Table) string {
	code := ""
	for _, col := range table.Columns {
		cn := col.ColumnName
		code += fmt.Sprintf("\n\t\t\t\t%s: %s[i].value,", cn, cn)
	}
	return code
}

func generateJsCode_doPutAll_then(table *dto.Table) string {
	code := ""
	for _, col := range table.Columns {
		cn := col.ColumnName
		code += fmt.Sprintf("\n\t\t\t\t%s[i].value = data.%s;", cn, cn)
		if col.IsUpdAble {
			code += fmt.Sprintf("\n\t\t\t\t%s_bk[i].value = data.%s;", cn, cn)
		}
	}
	code += "\n"
	for _, col := range table.Columns {
		if col.IsUpdAble {
			code += fmt.Sprintf("\n\t\t\t\t%s[i].classList.remove('changed');", col.ColumnName)
		}
	}
	return code
}

func generateJsCode_doPost(table *dto.Table) string {
	code := ""
	for _, col := range table.Columns {
		if col.IsInsAble {
			cn := col.ColumnName
			code += fmt.Sprintf("\n\tlet %s = document.getElementById('%s_new').value;", cn, cn)
		}
	}
	return code
}

func generateJsCode_doPost_if(table *dto.Table) string {
	code := ""
	isFirst := true
	for _, col := range table.Columns {
		if col.IsInsAble {
			if isFirst {
				code += fmt.Sprintf("(%s !== '')", col.ColumnName)
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t|| (%s !== '')", col.ColumnName)
			}
		}
	}
	return code
}

func generateJsCode_doPost_requestBody(table *dto.Table) string {
	code := ""
	for _, col := range table.Columns {
		if col.IsInsAble {
			cn := col.ColumnName
			code += fmt.Sprintf("\n\t\t\t%s: %s,", cn, cn)
		}
	}
	return code
}

var JS_COMMON_CODE = ""
const JS_FORMAT =
`
/* <tr></tr>を作成 （tbody末尾の新規登録用レコード）*/
const createTrNew = (elem) => {
	return %s
} 

/* <tr></tr>を作成 */
const createTr = (elem) => {
	return %s
} 


/* セットアップ */
const setUp = () => {
	fetch('api/%s')
	.then(response => response.json())
	.then(data  => renderTbody(data))
	.then(() => {%s
	});
}


/* 一括更新 */
const doPutAll = async () => {
	let successCount = 0;
	let errorCount = 0;
	%s

	for (let i = 0; i < %s.length; i++) {
		if (%s) {

			let requestBody = {%s
			}

			await fetch('api/%s', {
				method: 'PUT',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify(requestBody)
			})
			.then(response => {
				if (!response.ok){
					throw new Error(response.statusText);
				}
  				return response.json();
  			})
			.then(data => {%s

				successCount += 1;
			}).catch(error => {
				errorCount += 1;				
			})
		}
	}

	renderMessage('更新', successCount, true);
	renderMessage('更新', errorCount, false);
} 


/* 新規登録 */
const doPost = () => {%s

	if (%s)
	{
		let requestBody = {%s
		}

		fetch('api/%s', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(requestBody)
		})
		.then(response => {
			if (!response.ok){
				throw new Error(response.statusText);
			}
  			return response.json();
  		})
		.then(data => {
			document.getElementById('new').remove();

			let tmpElem = document.createElement('tbody');
			tmpElem.innerHTML = createTr(data);
			tmpElem.firstChild.addEventListener('change', changeAction);
			document.getElementById('records').appendChild(tmpElem.firstChild);

			tmpElem = document.createElement('tbody');
			tmpElem.innerHTML = createTrNew();
			document.getElementById('records').appendChild(tmpElem.firstChild);

			renderMessage('登録', 1, true);
		}).catch(error => {
			renderMessage('登録', 1, false);
		})
	}
}


/* 一括削除 */
const doDeleteAll = async () => {
	let ls = getDeleteTarget();
	let successCount = 0;
	let errorCount = 0;

	for (let x of ls) {
		await fetch('api/%s', {
			method: 'DELETE',
			headers: {'Content-Type': 'application/json'},
			body: x
		})
		.then(response => {
			if (!response.ok){
				throw new Error(response.statusText);
			}
			successCount += 1;
  		}).catch(error => {
			errorCount += 1;
		});
	}

	setUp();

	renderMessage('削除', successCount, true);
	renderMessage('削除', errorCount, false);
}
`

func generateHtmlCode_h2(table *dto.Table) string {
	if table.TableNameJp != "" {
		return fmt.Sprintf("%s（%s）", table.TableName, table.TableNameJp)
	} else {
		return table.TableName
	}
}

func generateHtmlCode_tr(table *dto.Table) string {
	code := ""
	for _, col := range table.Columns {
		if (col.IsNotNull || col.IsPrimaryKey) && (col.IsUpdAble || col.IsInsAble) {
			code += fmt.Sprintf("\n\t\t\t\t<th>%s<spnn style='color:red;'>*</spnn></th>", col.ColumnName)
		} else {
			code += fmt.Sprintf("\n\t\t\t\t<th>%s</th>", col.ColumnName)
		}
	}
	return code
}

const HTML_FORMAT =
`
<h2 class="ps-2 my-2">%s</h2>
<hr class="mt-0 mb-2">
<div class="w-100 vh-100 px-3">
	<div id=message></div>
	<button type="button" class="btn btn-danger" data-bs-toggle="modal" data-bs-target="#ModalDeleteAll">削除</button>
	<button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#ModalSaveAll">保存</button>
	<button type="button" class="btn btn-secondary" id="reload">リロード</button>
	<div class="table-responsive mt-2" style="height:70%s">
	<table class="table table-hover table-bordered table-sm">
		<thead>
			<tr class="fixed-table-header bg-light">
				<th>削除</th>%s
			</tr>
		</thead>
		<tbody id="records">
		</tbody>
	</table>
	</div>
</div>
<script src="/static/js/%s.js"></script>
`

func generateHtmlCodeHeader_ul(tables *[]dto.Table) string {
	code := ""
	for _, table := range *tables {
		tn := table.TableName
		code += fmt.Sprintf("\n\t\t\t<li class='nav-item'><a href='/mastertables/%s' class='nav-link text-white'>%s</a></li>", tn, tn)
	}
	return code
}

const HTML_HEADER_FORMAT = 
`<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name=”description“ content=““ />
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<link rel="stylesheet" href="/static/css/style.css">
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">
	<title>マスタメンテナンス</title>
</head>
<body>

<!-- 削除確認モーダル -->
<div class="modal" tabindex="-1" id="ModalDeleteAll">
<div class="modal-dialog">
	<div class="modal-content">
		<div class="modal-header">
			<h4 class="modal-title">削除</h4>
			<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
		</div>
		<div class="modal-body">
			<p>この操作は元には戻せません。よろしいですか？</p>
		</div>
		<div class="modal-footer">
			<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">キャンセル</button>
			<button type="button" class="btn btn-danger" data-bs-dismiss="modal" id="ModalDeleteAllOk">削除</button>
		</div>
	</div>
</div>
</div>

<!-- 保存確認モーダル -->
<div class="modal" tabindex="-1" id="ModalSaveAll">
<div class="modal-dialog">
	<div class="modal-content">
		<div class="modal-header">
			<h4 class="modal-title">保存</h4>
			<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
		</div>
		<div class="modal-body">
			<p>この操作は元には戻せません。よろしいですか？</p>
		</div>
		<div class="modal-footer">
			<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">キャンセル</button>
			<button type="button" class="btn btn-primary" data-bs-dismiss="modal" id="ModalSaveAllOk">保存</button>
		</div>
	</div>
</div>
</div>

<main>
	<!-- サイドバー -->
	<div class="d-flex flex-column flex-shrink-0 p-3 text-white bg-secondary" style="width: 280px;">
		<span class="fs-4">テーブル一覧</span>
		<hr class="mt-0 mb-2">
		<ul class="nav nav-pills flex-column mb-auto">%s
		</ul>
	</div>
	<!-- メインコンテンツ -->
	<div class="w-100 vh-100">
` 

const HTML_FOOTER_CODE = 
`
	</div>
	</div>
</main>
<footer>
	Copyright &copy; kodaimurakami. 2023. 
</footer>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/js/bootstrap.bundle.min.js" integrity="sha384-w76AqPfDkMBDXo30jS1Sgez6pr3x5MlQ1ZAGC+nuZB+EYdgRZgiwxhTBTkF7CXvN" crossorigin="anonymous"></script>
</body>
</html>
`

package generator

import (
	"os"
	"fmt"
	"strings"
	"time"
	"os/exec"
	"github.com/kodaimura/ddlparse"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/core/utils"
)


type generator struct {
	ddl string
	tables []ddlparse.Table
	rdbms string
	output string
}

type Generator interface {
	Generate() (string, error)
}

func NewGenerator(ddl string, rdbms string) (Generator, error) {
	var tables []ddlparse.Table
	var err error
	if (rdbms == "postgresql") {
		tables, err = ddlparse.ParsePostgreSQL(ddl)
	} else if (rdbms == "mysql") {
		tables, err = ddlparse.ParseMySQL(ddl)
	} else if (rdbms == "sqlite3") {
		tables, err = ddlparse.ParseSQLite(ddl)
	} else {
		tables, err = ddlparse.ParseSQLite(ddl)
	}
	if err != nil {
		return &generator{}, err
	}

	return &generator{
		ddl: ddl,
		tables: tables,
		rdbms: rdbms,
		output: "./output",
	}, nil
}

func (gen *generator) Generate() (string, error) {
	dir, path, err := gen.createWorkDir()
	if err != nil {
		return "", err
	}
	if err := gen.generateSource(path); err != nil {
		return "", err
	}
	filename, err := gen.zipWorkDir(dir)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (gen *generator) createWorkDir() (string, string, error) {
	dir := fmt.Sprintf(
		"%s-%s", 
		time.Now().Format("2006-01-02-15-04-05"), 
		utils.RandomString(10),
	)
	path := fmt.Sprintf("%s/%s", gen.output, dir)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return "", "", err
	}
	return dir, path, nil
}

func (gen *generator) zipWorkDir(dir string) (string, error) {
	current, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer os.Chdir(current)

	if err := os.Chdir(gen.output); err != nil {
		logger.Error(err.Error())
		return "", err
	}

	zip := fmt.Sprintf("%s.zip", dir)
	if err := exec.Command("zip", "-rm", zip, dir).Run(); err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return zip, nil
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////  生成処理  ///////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func (gen *generator) generateSource(path string) error {
	path = fmt.Sprintf("%s/masmaint", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.copyTemplate(path); err != nil {
		return err
	}
	if err := gen.copySomeFiles(path); err != nil {
		return err
	}
	if err := gen.generateInternal(path); err != nil {
		return err
	}
	if err := gen.generateWeb(path); err != nil {
		return err
	}
	if err := gen.generateScripts(path); err != nil {
		return err
	}
	return nil	
}

func (gen *generator) copyTemplate(path string) error {
	origin := "_template/masmaint"

	if err := CopyDir(origin, path); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (gen *generator) copySomeFiles(path string) error {
	origin := fmt.Sprintf("_template/%s", gen.rdbms)

	if err := CopyDir(origin, path); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////  internal  //////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// internal 生成
func (gen *generator) generateInternal(path string) error {
	path = fmt.Sprintf("%s/internal", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateModule(path); err != nil {
		return err
	}
	if err := gen.generateServer(path); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  internal/module  //////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// module 生成
func (gen *generator) generateModule(path string) error {
	path = fmt.Sprintf("%s/module", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateTableModules(path); err != nil {
		return err
	}
	return nil
}

func (gen *generator) generateTableModules(path string) error {
	for _, table := range gen.tables {
		if err := gen.generateTableModule(path, table); err != nil {
			return err
		}
	}
	return nil
}

// module/table_name 生成
func (gen *generator) generateTableModule(path string, table ddlparse.Table) error {
	tn := strings.ToLower(table.Name)
	path = fmt.Sprintf("%s/%s", path, tn)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateTableModuleFiles(path, table); err != nil {
		return err
	}
	return nil
}

func (gen *generator) generateTableModuleFiles(path string, table ddlparse.Table) error {
	if err := gen.generateControllerGoFile(path, table); err != nil {
		return err
	}
	if err := gen.generateModelGoFile(path, table); err != nil {
		return err
	}
	if err := gen.generateRequestGoFile(path, table); err != nil {
		return err
	}
	if err := gen.generateServiceGoFile(path, table); err != nil {
		return err
	}
	if err := gen.generateRepositoryGoFile(path, table); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  controller.go  ////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// controller.go 生成
func (gen *generator) generateControllerGoFile(path string, table ddlparse.Table) error {
	path = fmt.Sprintf("%s/controller.go", path)
	code := gen.codeControllerGo(table)
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// controller.go コード生成
func (gen *generator) codeControllerGo(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	return fmt.Sprintf(
		FORMAT_CONTROLLER, 
		tn, tn, tn, tn, tn, tn, tn,
	)
}

/////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////  model.go  //////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// model.go 生成
func (gen *generator) generateModelGoFile(path string, table ddlparse.Table) error {
	path = fmt.Sprintf("%s/model.go", path)
	code := gen.codeModelGo(table)
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// model.go コード生成
func (gen *generator)codeModelGo(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	
	fields := ""
	for _, c := range table.Columns {
		cn := strings.ToLower(c.Name)
		fields += "\t" + gen.getFieldName(cn ,tn) + " ";
		if gen.isNullColumn(c, table.Constraints) {
			fields += "*"
		}
		fields += gen.dataTypeToGoType(c.DataType.Name) + " "
		fields += "`db:\"" + cn + "\" json:\"" + cn + "\"`\n"
	}
	fields = strings.TrimSuffix(fields, "\n")
	return fmt.Sprintf(
		FORMAT_MODEL, 
		tn, tnp, fields,
	)
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////  request.go  /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// request.go 生成
func (gen *generator) generateRequestGoFile(path string, table ddlparse.Table) error {
	path = fmt.Sprintf("%s/request.go", path)
	code := gen.codeRequestGo(table)
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// request.go コード生成
func (gen *generator)codeRequestGo(table ddlparse.Table) string {
	return fmt.Sprintf(
		FORMAT_REQUEST, 
		strings.ToLower(table.Name), 
		gen.codeRequestPostBodyFields(table), 
		gen.codeRequestPutBodyFields(table),
		gen.codeRequestDeleteBodyFields(table),
	)
}

func (gen *generator)codeRequestPostBodyFields(table ddlparse.Table) string {
    tn := strings.ToLower(table.Name)
    
    code := ""
    for _, c := range gen.getInsertColumns(table) {
        cn := strings.ToLower(c.Name)
        code += fmt.Sprintf("\t%s ", gen.getFieldName(cn ,tn))
        if gen.isNullColumn(c, table.Constraints) {
            code += "*"
        }
        code += fmt.Sprintf("%s `json:\"%s\"", gen.dataTypeToGoType(c.DataType.Name), cn)
        if !gen.isNullColumn(c, table.Constraints) {
            code += " binding:\"required\""
        }
        code += "`\n"
    }
    return strings.TrimSuffix(code, "\n")
}

func (gen *generator)codeRequestPutBodyFields(table ddlparse.Table) string {
    tn := strings.ToLower(table.Name)
    
    code := ""
    for _, c := range table.Columns {
        if strings.Contains(c.Name, "_at") || strings.Contains(c.Name, "_AT") {
            continue
        }
        cn := strings.ToLower(c.Name)
        code += fmt.Sprintf("\t%s ", gen.getFieldName(cn ,tn))
        if gen.isNullColumn(c, table.Constraints) {
            code += "*"
        }
        code += fmt.Sprintf("%s `json:\"%s\"", gen.dataTypeToGoType(c.DataType.Name), cn)
        if !gen.isNullColumn(c, table.Constraints) {
            code += " binding:\"required\""
        }
        code += "`\n"
    }
    return strings.TrimSuffix(code, "\n")
}

func (gen *generator)codeRequestDeleteBodyFields(table ddlparse.Table) string {
    tn := strings.ToLower(table.Name)
    code := ""
    for _, c := range gen.getPrimaryKeyColumns(table) {
        cn := strings.ToLower(c.Name)
        code += fmt.Sprintf("\t%s ", gen.getFieldName(cn ,tn))
        if gen.isNullColumn(c, table.Constraints) {
            code += "*"
        }
        code += fmt.Sprintf("%s `json:\"%s\"", gen.dataTypeToGoType(c.DataType.Name), cn)
        if !gen.isNullColumn(c, table.Constraints) {
            code += " binding:\"required\""
        }
        code += "`\n"
    }
    return strings.TrimSuffix(code, "\n")
}

/////////////////////////////////////////////////////////////////////////////////
////////////////////////////////  repository.go  ///////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// repository.go 生成
func (gen *generator)generateRepositoryGoFile(path string, table ddlparse.Table) error {
	path = fmt.Sprintf("%s/repository.go", path)
	code := gen.codeRepositoryGo(table)
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// repository.go コード生成
func (gen *generator)codeRepositoryGo(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)

	return fmt.Sprintf(
		FORMAT_REPOSITORY,
		tn,
		gen.codeRepositoryInterface(table),
		gen.codeRepositoryGet(table),
		gen.codeRepositoryGetOne(table),
		gen.codeRepositoryInsert(table),
		gen.codeRepositoryUpdate(table),
		gen.codeRepositoryDelete(table),
	)
}

func (gen *generator)codeRepositoryInterface(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	aicol, found := gen.getAutoIncrementColumn(table)

	retType := "error"
	if found {
		retType = fmt.Sprintf("(%s, error)", gen.dataTypeToGoType(aicol.DataType.Name))
	}

	return fmt.Sprintf(
		FORMAT_REPOSITORY_INTERFACE,
		tni, tnp, tnp, 
		tni, tnp, tnp, 
		tni, tnp, retType,
		tni, tnp, 
		tni, tnp,
	)
}

func (gen *generator)codeRepositoryGet(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	query := "\n\t`SELECT"
	for i, c := range table.Columns {
		if i == 0 {
			query += fmt.Sprintf("\n\t\t%s", c.Name)
		} else {
			query += fmt.Sprintf("\n\t\t,%s", c.Name)
		}
	}
	query += fmt.Sprintf("\n\t FROM %s `", tn)

	scan := "\n"
	for _, c := range table.Columns {
		scan += fmt.Sprintf("\t\t\t&%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
	}
	scan += "\t\t"

	return fmt.Sprintf(
		FORMAT_REPOSITORY_GET,
		tni, tnp, tnp, tni, 
		query,
		tnp, tnp, tni, tnp,
		scan,
		tnp, tni,
	) 
}

func (gen *generator)codeRepositoryGetOne(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	query := "\n\t`SELECT"
	for i, c := range table.Columns {
		if i == 0 {
			query += fmt.Sprintf("\n\t\t%s", c.Name)
		} else {
			query += fmt.Sprintf("\n\t\t,%s", c.Name)
		}
	}
	query += fmt.Sprintf("\n\t FROM %s `", tn)

	scan := "\n"
	for _, c := range table.Columns {
		scan += fmt.Sprintf("\t\t&ret.%s,\n", gen.getFieldName(c.Name ,tn))
	}
	scan += "\t"

	return fmt.Sprintf(
		FORMAT_REPOSITORY_GETONE,
		tni, tnp, tnp, tnp, tni, 
		query, scan,
	) 
}

func (gen *generator)getBindVar(n int) string {
	if gen.rdbms == "postgresql" {
		return fmt.Sprintf("$%d", n)
	} else {
		return "?"
	}
}

func (gen *generator)concatBindVariableWithCommas(bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, gen.getBindVar(i))
	}
	return strings.Join(ls, ",")
}

func (gen *generator)codeRepositoryInsert(table ddlparse.Table) string {
	_, found := gen.getAutoIncrementColumn(table)
	if found {
		if gen.rdbms == "mysql" {
			return gen.codeRepositoryInsertAIMySQL(table)
		} else {
			return gen.codeRepositoryInsertAI(table)
		}	
	}
	return gen.codeRepositoryInsertNomal(table)
}

func (gen *generator)codeRepositoryInsertNomal(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	inscols := gen.getInsertColumns(table)

	query := fmt.Sprintf("\n\t`INSERT INTO %s (", tn)
	bindCount := 0
	for i, c := range inscols {
		bindCount += 1
		if i == 0 {
			query += fmt.Sprintf("\n\t\t%s", c.Name)
		} else {
			query += fmt.Sprintf("\n\t\t,%s", c.Name)
		}	
	}
	query += fmt.Sprintf("\n\t ) VALUES(%s)`\n", gen.concatBindVariableWithCommas(bindCount))

	binds := "\n"
	for _, c := range inscols {
		binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
	}
	binds += "\t"

	return fmt.Sprintf(
		FORMAT_REPOSITORY_INSERT,
		tni, tnp,
		query, binds,
	) 
}

func (gen *generator)codeRepositoryInsertAI(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	inscols := gen.getInsertColumns(table)
	aicol, _ := gen.getAutoIncrementColumn(table)
	aicn := strings.ToLower(aicol.Name)
	aicnc := SnakeToCamel(aicn)

	query := fmt.Sprintf("\n\t`INSERT INTO %s (", tn)
	bindCount := 0
	for i, c := range inscols {
		bindCount += 1
		if i == 0 {
			query += fmt.Sprintf("\n\t\t%s", c.Name)
		} else {
			query += fmt.Sprintf("\n\t\t,%s", c.Name)
		}
	}
	query += fmt.Sprintf("\n\t ) VALUES(%s)", gen.concatBindVariableWithCommas(bindCount))
	query += fmt.Sprintf("\n\t RETURNING %s`\n", aicn)

	binds := "\n"
	for _, c := range inscols {
		binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
	}
	binds += "\t"

	return fmt.Sprintf(
		FORMAT_REPOSITORY_INSERT_AI,
		tni, tnp,
		query, binds,
		aicnc, aicnc, aicnc, aicnc,
	) 
}

func (gen *generator)codeRepositoryInsertAIMySQL(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	inscols := gen.getInsertColumns(table)
	aicol, _ := gen.getAutoIncrementColumn(table)
	aicn := strings.ToLower(aicol.Name)
	aicnc := SnakeToCamel(aicn)

	query := fmt.Sprintf("\n\t`INSERT INTO %s (", tn)
	bindCount := 0
	for i, c := range inscols {
		bindCount += 1
		if i == 0 {
			query += fmt.Sprintf("\n\t\t%s", c.Name)
		} else {
			query += fmt.Sprintf("\n\t\t,%s", c.Name)
		}
		
	}
	query += fmt.Sprintf("\n\t ) VALUES(%s)`\n", gen.concatBindVariableWithCommas(bindCount))

	binds := "\n"
	for _, c := range inscols {
		binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
	}
	binds += "\t"

	return fmt.Sprintf(
		FORMAT_REPOSITORY_INSERT_AI_MYSQL,
		tni, tnp,
		query, binds,
		aicnc, aicnc, aicnc, aicnc,
	) 
}

func (gen *generator)codeRepositoryUpdate(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	updcols := gen.getUpdateColumns(table)
	pkcols := gen.getPrimaryKeyColumns(table)

	bindCount := 0
	query := fmt.Sprintf("\n\t`UPDATE %s\n\t SET ", tn)
	for i, c := range updcols {
		bindCount += 1
		if i > 0 {
			query += "\t\t,"
		}
		query += fmt.Sprintf("%s = %s\n", c.Name, gen.getBindVar(bindCount))
	}
	query += "\t WHERE "
	for i, c := range pkcols {
		bindCount += 1
		if i > 0 {
			query += "\n\t   AND "
		}
		query += fmt.Sprintf("%s = %s", c.Name, gen.getBindVar(bindCount))
	}
	query += "`"

	binds := "\n"
	for _, c := range updcols {
		binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
	}
	for _, c := range pkcols {
		binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
	}
	binds += "\t"

	return fmt.Sprintf(
		FORMAT_REPOSITORY_UPDATE,
		tni, tnp,
		query, binds,
	) 
}

func (gen *generator)codeRepositoryDelete(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	return fmt.Sprintf(
		FORMAT_REPOSITORY_DELETE, 
		tni, tnp, tni, tn,
	) 
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////  service.go  /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// service.go 生成
func (gen *generator) generateServiceGoFile(path string, table ddlparse.Table) error {
	path = fmt.Sprintf("%s/service.go", path)
	code := gen.codeServiceGo(table)
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// service.go コード生成
func (gen *generator) codeServiceGo(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	return fmt.Sprintf(
		FORMAT_SERVICE, 
		tn, tnp, tnp, tnp,
		gen.codeServiceGet(table),
		gen.codeServiceCreate(table),
		gen.codeServiceUpdate(table),
		gen.codeServiceDelete(table),
	)
}

func (gen *generator)codeServiceGet(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)

	return fmt.Sprintf(
		FORMAT_SERVICE_GET,
		tnp, tnp, tnp,
	) 
}

func (gen *generator)codeServiceCreate(table ddlparse.Table) string {
	_, found := gen.getAutoIncrementColumn(table)
	if found {
		return gen.codeServiceCreateAI(table)
	}
	return gen.codeServiceCreateNomal(table)
}

func (gen *generator)codeServiceCreateNomal(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)

	s1 := ""
	for i, c := range gen.getPrimaryKeyColumns(table) {
		if i > 0 {
			s1 += ", "
		} 
		fn := gen.getFieldName(c.Name ,tn)
		s1 += fmt.Sprintf("%s: input.%s", fn, fn)
	}

	return fmt.Sprintf(
		FORMAT_SERVICE_CREATE,
		tnp, tnp, tnp, tnp, tnp, s1, tnp,
	) 
}

func (gen *generator)codeServiceCreateAI(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	aicol, _ := gen.getAutoIncrementColumn(table)
	aicn := strings.ToLower(aicol.Name)
	aicnc := SnakeToCamel(aicn)
	fn := gen.getFieldName(aicn ,tn)

	return fmt.Sprintf(
		FORMAT_SERVICE_CREATE_AI,
		tnp, tnp, aicnc, tnp, tnp, tnp, fmt.Sprintf("%s: %s", fn, aicnc), tnp,
	) 
}

func (gen *generator)codeServiceUpdate(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)

	s1 := ""
	for i, c := range gen.getPrimaryKeyColumns(table) {
		if i > 0 {
			s1 += ", "
		} 
		fn := gen.getFieldName(c.Name ,tn)
		s1 += fmt.Sprintf("%s: input.%s", fn, fn)
	}

	return fmt.Sprintf(
		FORMAT_SERVICE_UPDATE,
		tnp, tnp, tnp, tnp, tnp, s1, tnp,
	) 
}

func (gen *generator)codeServiceDelete(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)

	return fmt.Sprintf(FORMAT_SERVICE_DELETE, tnp) 
}

/////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  internal/server  //////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// server 生成
func (gen *generator) generateServer(path string) error {
	path = fmt.Sprintf("%s/server", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateRouterGoFile(path); err != nil {
		return err
	}
	return nil
}

// router.go 生成
func (gen *generator) generateRouterGoFile(path string) error {
	path = fmt.Sprintf("%s/router.go", path)
	code := gen.codeRouterGo()
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// router.go コード生成
func (gen *generator) codeRouterGo() string {
	s1 := ""
	for _, table := range gen.tables {
		s1 += fmt.Sprintf("\t\"masmaint/internal/module/%s\"\n", strings.ToLower(table.Name))
	}
	s1 = strings.TrimSuffix(s1, "\n")
	return fmt.Sprintf(
		FORMAT_ROUTER, 
		s1,
		gen.CodeRouterSetWebRouter(),
		gen.CodeRouterSetApiRouter(),
	)
}

func (gen *generator) CodeRouterSetWebRouter() string {
	s1 := ""
	for _, table := range gen.tables {
		tn := strings.ToLower(table.Name)
		tnc := SnakeToCamel(tn)
		s1 += fmt.Sprintf("\t%sController := %s.NewController()\n", tnc, tn)
	}
	s1 = strings.TrimSuffix(s1, "\n")
	s2 := ""
	for _, table := range gen.tables {
		tn := strings.ToLower(table.Name)
		tnc := SnakeToCamel(tn)
		s2 += fmt.Sprintf("\t\tauth.GET(\"/%s\", %sController.GetPage)\n", tn, tnc)
	}
	s2 = strings.TrimSuffix(s2, "\n")
	return fmt.Sprintf(
		FORMAT_ROUTER_SETWEB, 
		s1, s2,
	)
}

func (gen *generator) CodeRouterSetApiRouter() string {
	s1 := ""
	for _, table := range gen.tables {
		tn := strings.ToLower(table.Name)
		tnc := SnakeToCamel(tn)
		s1 += fmt.Sprintf("\t%sController := %s.NewController()\n", tnc, tn)
	}
	s1 = strings.TrimSuffix(s1, "\n")
	s2 := ""
	for _, table := range gen.tables {
		tn := strings.ToLower(table.Name)
		tnc := SnakeToCamel(tn)
		s2 += fmt.Sprintf("\t\tauth.GET(\"/%s\", %sController.Get)\n", tn, tnc)
		s2 += fmt.Sprintf("\t\tauth.POST(\"/%s\", %sController.Post)\n", tn, tnc)
		s2 += fmt.Sprintf("\t\tauth.PUT(\"/%s\", %sController.Put)\n", tn, tnc)
		s2 += fmt.Sprintf("\t\tauth.DELETE(\"/%s\", %sController.Delete)\n\n", tn, tnc)
	}
	s2 = strings.TrimSuffix(s2, "\n\n")
	return fmt.Sprintf(
		FORMAT_ROUTER_SETAPI, 
		s1, s2,
	)
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////  web  ////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// web 生成
func (gen *generator) generateWeb(path string) error {
	path = fmt.Sprintf("%s/web", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateStatic(path); err != nil {
		return err
	}
	if err := gen.generateTemplate(path); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////  web/static  /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// static 生成
func (gen *generator) generateStatic(path string) error {
	path = fmt.Sprintf("%s/static", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateJs(path); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
////////////////////////////////  web/static/js  ///////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// js 生成
func (gen *generator) generateJs(path string) error {
	path = fmt.Sprintf("%s/js", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateTableJsFiles(path); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
////////////////////////////////  table_name.js  ///////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func (gen *generator) generateTableJsFiles(path string) error {
	for _, table := range gen.tables {
		if err := gen.generateTableJsFile(path, table); err != nil {
			return err
		}
	}
	return nil
}

// table_name.js 生成
func (gen *generator) generateTableJsFile(path string, table ddlparse.Table) error {
	tn := strings.ToLower(table.Name)
	path = fmt.Sprintf("%s/%s.js", path, tn)
	code := gen.codeTableJs(table)
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// table_name.js コード生成
func (gen *generator) codeTableJs(table ddlparse.Table) string {
	return fmt.Sprintf(
		FORMAT_JS, 
		gen.codeJsCreateTrNew(table),
		gen.codeJsCreateTr(table),
		gen.codeJsGetRows(table),
		gen.codeJsPutRows(table),
		gen.codeJsPostRow(table),
		gen.codeJsDeleteRows(table),
	)
}

func (gen *generator) codeJsCreateTrNew(table ddlparse.Table) string {
	s1 := "\n\t\t<td></td>"
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			s1 += fmt.Sprintf("\n\t\t<td><input type='text' id='%s_new'></td>", strings.ToLower(c.Name))
		} else {
			s1 += "\n\t\t<td><input type='text' disabled></td>"
		}
	}
	return fmt.Sprintf(FORMAT_JS_CREATETRNEW, s1)
}

func (gen *generator) codeJsCreateTr(table ddlparse.Table) string {
	s1 := "\n\t\t<td><input class='form-check-input' type='checkbox' name='del' value='${JSON.stringify(elem)}'></td>"
	for _, c := range table.Columns {
		cn := strings.ToLower(c.Name)
		if gen.isUpdateColumn(c) {
			s1 += fmt.Sprintf(
				"\n\t\t<td><input type='text' name='%s' value='${nullToEmpty(elem.%s)}'><input type='hidden' name='%s_bk' value='${nullToEmpty(elem.%s)}'></td>",
				cn, cn, cn, cn,
			)
		} else {
			s1 += fmt.Sprintf(
				"\n\t\t<td><input type='text' name='%s' value='${nullToEmpty(elem.%s)}' disabled></td>", 
				cn, cn,
			)
		}
	}
	return fmt.Sprintf(FORMAT_JS_CREATETR, s1)
}

func (gen *generator) codeJsGetRows(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	s1 := ""
	for _, c := range gen.getUpdateColumns(table) {
		cn := strings.ToLower(c.Name)
		s1 += fmt.Sprintf("\taddChangeEvent('%s');\n", cn)
	}
	s1 = strings.TrimSuffix(s1, "\n")
	return fmt.Sprintf(FORMAT_JS_GETROWS, tn, s1)
}

func (gen *generator) codeJsPutRows(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	updcols := gen.getUpdateColumns(table)
	s1 := ""
	for _, c := range table.Columns {
		cn := strings.ToLower(c.Name)
		s1 += fmt.Sprintf("\tconst %s = document.getElementsByName('%s');\n", cn, cn)
	}
	s1 = strings.TrimSuffix(s1, "\n")
	s2 := ""
	for _, c := range updcols {
		cn := strings.ToLower(c.Name)
		s2 += fmt.Sprintf("\tconst %s_bk = document.getElementsByName('%s_bk');\n", cn, cn)
	}
	s2 = strings.TrimSuffix(s2, "\n")
	s3 := strings.ToLower(table.Columns[0].Name)
	s4 := ""
	for _, c := range updcols {
		cn := strings.ToLower(c.Name)
		s4 += fmt.Sprintf("\t\t\t'%s': %s[i],\n", cn, cn)
	}
	s4 = strings.TrimSuffix(s4, "\n")
	s5 := ""
	for _, c := range updcols {
		cn := strings.ToLower(c.Name)
		s5 += fmt.Sprintf("\t\t\t'%s': %s_bk[i],\n", cn, cn)
	}
	s5 = strings.TrimSuffix(s5, "\n")
	s6 := ""
	for _, c := range table.Columns {
		cn := strings.ToLower(c.Name)
		ctype := gen.dataTypeToGoType(c.DataType.Name)
		if ctype == "int" {
			s6 += fmt.Sprintf("\t\t\t\t%s: parseIntOrReturnOriginal(%s[i].value),\n", cn, cn)
		} else if ctype == "float64" {
			s6 += fmt.Sprintf("\t\t\t\t%s: parseFloatOrReturnOriginal(%s[i].value),\n", cn, cn)
		} else if gen.isNullColumn(c, table.Constraints) {
			s6 += fmt.Sprintf("\t\t\t\t%s: emptyToNull(%s[i].value),\n", cn, cn)
		} else {
			s6 += fmt.Sprintf("\t\t\t\t%s: %s[i].value,\n", cn, cn)
		}
	}
	s6 = strings.TrimSuffix(s6, "\n")
	s7 := ""
	for _, c := range table.Columns {
		cn := strings.ToLower(c.Name)
		s7 += fmt.Sprintf("\t\t\t\t%s[i].value = data.%s;\n", cn, cn)
	}
	for _, c := range updcols {
		cn := strings.ToLower(c.Name)
		s7 += fmt.Sprintf("\t\t\t\t%s_bk[i].value = data.%s;\n", cn, cn)
	}
	s7 = strings.TrimSuffix(s7, "\n")

	return fmt.Sprintf(
		FORMAT_JS_PUTROWS, 
		s1, s2, s3, s4, s5, s6, tn, s7, tn,
	)
}

func (gen *generator) codeJsPostRow(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	inscols := gen.getInsertColumns(table)
	s1 := ""
	for _, c := range inscols {
		cn := strings.ToLower(c.Name)
		s1 += fmt.Sprintf("\t\t'%s': document.getElementById('%s_new'),\n", cn, cn)
	}
	s1 = strings.TrimSuffix(s1, "\n")
	s2 := ""
	for _, c := range inscols{
		cn := strings.ToLower(c.Name)
		ctype := gen.dataTypeToGoType(c.DataType.Name)
		if ctype == "int" {
			s2 += fmt.Sprintf("\t\t\t%s: parseIntOrReturnOriginal(rowMap.%s.value),\n", cn, cn)
		} else if ctype == "float64" {
			s2 += fmt.Sprintf("\t\t\t%s: parseFloatOrReturnOriginal(rowMap.%s.value),\n", cn, cn)
		} else if gen.isNullColumn(c, table.Constraints) {
			s2 += fmt.Sprintf("\t\t\t%s: emptyToNull(rowMap.%s.value),\n", cn, cn)
		} else {
			s2 += fmt.Sprintf("\t\t\t%s: rowMap.%s.value,\n", cn, cn)
		}
	}
	s2 = strings.TrimSuffix(s2, "\n")

	return fmt.Sprintf(
		FORMAT_JS_POSTROW, 
		s1, s2, tn, tn,
	)
}

func (gen *generator) codeJsDeleteRows(table ddlparse.Table) string {
	return fmt.Sprintf(
		FORMAT_JS_DELETEROWS,
		strings.ToLower(table.Name),
	)
}

/////////////////////////////////////////////////////////////////////////////////
////////////////////////////////  web/template  ////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// template 生成
func (gen *generator) generateTemplate(path string) error {
	path = fmt.Sprintf("%s/template", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateMenuHtmlFile(path); err != nil {
		return err
	}
	if err := gen.generateTableHtmlFiles(path); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////  _menu.html  /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// _menu.html 生成
func (gen *generator) generateMenuHtmlFile(path string) error {
	path = fmt.Sprintf("%s/_menu.html", path)
	code := gen.codeMenuHtml()
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// _menu.html コード生成
func (gen *generator) codeMenuHtml() string {
	s1 := ""
	for _, table := range gen.tables {
		tn := strings.ToLower(table.Name)
		s1 += fmt.Sprintf("\t\t<li class='nav-item'><a href='/%s' class='nav-link py-1'>%s</a></li>\n", tn, tn)
	}
	s1 = strings.TrimSuffix(s1, "\n")
	return fmt.Sprintf(FORMAT_TEMPLATE_MENU, s1)
}

/////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  table_name.html  //////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func (gen *generator) generateTableHtmlFiles(path string) error {
	for _, table := range gen.tables {
		if err := gen.generateTableHtmlFile(path, table); err != nil {
			return err
		}
	}
	return nil
}

// table_name.html 生成
func (gen *generator) generateTableHtmlFile(path string, table ddlparse.Table) error {
	tn := strings.ToLower(table.Name)
	path = fmt.Sprintf("%s/%s.html", path, tn)
	code := gen.codeTableHtml(table)
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

// table_name.html コード生成
func (gen *generator) codeTableHtml(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	s1 := ""
	for _, c := range table.Columns {
		cn := strings.ToLower(c.Name)
		if gen.isNullColumn(c, table.Constraints) || !gen.isInsertColumn(c) {
			s1 += fmt.Sprintf("\t\t\t\t\t\t\t\t<th>%s</th>\n", cn)
		} else {
			s1 += fmt.Sprintf("\t\t\t\t\t\t\t\t<th>%s<spnn class=\"text-danger\">*</spnn></th>\n", cn)
		}
	}
	s1 = strings.TrimSuffix(s1, "\n")
	return fmt.Sprintf(
		FORMAT_TEMPLATE, 
		tn, s1, tn,
	)
}

/////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////  scripts  //////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// scripts 生成
func (gen *generator) generateScripts(path string) error {
	path = fmt.Sprintf("%s/scripts", path)
	if err := MakeDirAll(path); err != nil {
		logger.Error(err.Error())
		return err
	}
	if err := gen.generateCreateTableSqlFile(path); err != nil {
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////  create-table.sql  ///////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// create-table.sql 生成
func (gen *generator) generateCreateTableSqlFile(path string) error {
	path = fmt.Sprintf("%s/create-table.sql", path)
	code := gen.ddl
	if err := WriteFile(path, code); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
///////////////////////////////  コード生成用共通  ///////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// カラム名 -> Goフィールド名
func (gen *generator) getFieldName(columnName, tableName string) string {
	cn := strings.ToLower(columnName)
	tn := strings.ToLower(tableName)
	pf := tn + "_"
	if (strings.HasPrefix(cn, pf)) {
		cn = cn[len(pf):]
	}
	return SnakeToPascal(cn)
}

// データ型 -> Goデータ型
func (gen *generator) dataTypeToGoType(dataType string) string {
	dataType = strings.ToUpper(dataType)

	if (strings.Contains(dataType, "INT") || strings.Contains(dataType, "BIT") || strings.Contains(dataType, "SERIAL")) {
		return "int"
	} else if strings.Contains(dataType, "NUMERIC") ||
		strings.Contains(dataType, "DECIMAL") ||
		strings.Contains(dataType, "FLOAT") ||
		strings.Contains(dataType, "REAL") ||
		strings.Contains(dataType, "DOUBLE") {
		return "float64"
	} else {
		return "string"
	}
}

// Null許容のカラムか判定
func (gen *generator) isNullColumn(column ddlparse.Column, constraints ddlparse.TableConstraint) bool {
	if (column.Constraint.IsNotNull) {
		return false
	}
	if (column.Constraint.IsPrimaryKey) {
		return false
	}
	if (column.Constraint.IsAutoincrement) {
		return false
	}

	for _, pk := range constraints.PrimaryKey {
		for _, name := range pk.ColumnNames {
			if (column.Name == name) {
				return false
			}
		}
	}
	return true
}

// INSERTで指定するカラムか判定
func (gen *generator)isInsertColumn(c ddlparse.Column) bool {
	if c.Constraint.IsAutoincrement {
		return false
	}
	if strings.Contains(strings.ToUpper(c.DataType.Name), "SERIAL") {
		return false
	}
	if strings.Contains(c.Name, "_at") || strings.Contains(c.Name, "_AT") {
		return false
	}
	return true
}

// UPDATEで指定するカラムか判定
func (gen *generator)isUpdateColumn(c ddlparse.Column) bool {
	if c.Constraint.IsAutoincrement {
		return false
	}
	if strings.Contains(strings.ToUpper(c.DataType.Name), "SERIAL") {
		return false
	}
	if c.Constraint.IsPrimaryKey {
		return false
	}
	if strings.Contains(c.Name, "_at") || strings.Contains(c.Name, "_AT") {
		return false
	}
	return true
}

// INSERTするカラムのリストを取得
func (gen *generator)getInsertColumns(table ddlparse.Table) []ddlparse.Column {
	ret := []ddlparse.Column{}
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			ret = append(ret, c)
		}	
	}
	return ret
}

// UPDATEするカラムのリストを取得
func (gen *generator)getUpdateColumns(table ddlparse.Table) []ddlparse.Column {
	ret := []ddlparse.Column{}
	for _, c := range table.Columns {
		if gen.isUpdateColumn(c) {
			ret = append(ret, c)
		}	
	}
	return ret
}

// 主キーカラムのリストを取得
func (gen *generator)getPrimaryKeyColumns(table ddlparse.Table) []ddlparse.Column {
	pkcols := []string{}
	for _, pk := range table.Constraints.PrimaryKey {
		for _, name := range pk.ColumnNames {
			pkcols = append(pkcols, name)
		}
	}

	names := []string{}
	ret := []ddlparse.Column{}
	for _, c := range table.Columns {
		if c.Constraint.IsPrimaryKey || Contains(pkcols, c.Name) || strings.Contains(strings.ToUpper(c.DataType.Name), "SERIAL"){
			if !Contains(names, c.Name) {
				names = append(names, c.Name)
				ret = append(ret, c)
			}
		}
	}
	return ret
}

// AUTO_INCREMENTのカラムを取得（1つ以下である前提）
func (gen *generator)getAutoIncrementColumn(table ddlparse.Table) (ddlparse.Column, bool) {
	for _, c := range table.Columns {
		if c.Constraint.IsAutoincrement {
			return c, true
		}
		if strings.Contains(strings.ToUpper(c.DataType.Name), "SERIAL") {
			return c, true
		}
	}
	return ddlparse.Column{}, false
}
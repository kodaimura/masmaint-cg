package generator

import (
	"os"
	"fmt"
	"strings"
	"github.com/kodaimura/ddlparse"

	"masmaint-cg/internal/core/logger"
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

	if err := gen.copyTemplate(); err != nil {
		return err
	}
	if err := gen.generateModule(); err != nil {
		return err
	}
	if err := gen.generateWeb(); err != nil {
		return err
	}

	return nil	
}

// templateコピー
func (gen *generator) copyTemplate() error {
	source := "_template/masmaint"
	destination := gen.path

	err := CopyDir(source, destination)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// module生成
func (gen *generator) generateModules() error {
	path := gen.path + "_template/masmaint/internal/module/"
	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}

	for table := range gen.tables {
		if err := gen.generateModule(path, table); err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	return nil
}

// module/:table_name生成
func (gen *generator) generateModule(path string, table ddlparse.Table) error {
	pkg := strings.ToLower(table.Name)
	path = path + pkg
	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return err
	}
	return gen.generateModuleFiles(path)
}

// module/:table_name内のファイル生成
func (gen *generator) generateModuleFiles(path string, table ddlparse.Table) error {
	if err := gen.generateControllerFile(path, table); err != nil {
		return err
	}
	if err := gen.generateModelFile(path, table); err != nil {
		return err
	}
	if err := gen.generateRequestFile(path, table); err != nil {
		return err
	}
	if err := gen.generateRepositoryFile(path, table); err != nil {
		return err
	}
	if err := gen.generateServiceFile(path, table); err != nil {
		return err
	}
	return nil
}

// controller.go生成
func (gen *generator) generateControllerFile(table ddlparse.Table, path string) error {
	code := codeController(table)
	err := WriteFile(fmt.Sprintf("%scontroller.go", path), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// controller.go コード生成
func (gen *generator) codeController(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	return fmt.Sprintf(
		CONTROLLER_FORMAT, 
		tn, tn, tn, tn, tn, tn, tn,
	)
}

// model.go生成
func (gen *generator) generateModelFile(table ddlparse.Table, path string) error {
	code := codeModel(table)
	err := WriteFile(fmt.Sprintf("%smodel.go", path), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

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


func (gen *generator) getFieldName(columnName, tableName string) string {
	cn := strings.ToLower(columnName)
	tn := strings.ToLower(tableName)
	pf := tn + "_"
	if (strings.HasPrefix(cn, pf)) {
		cn = cn[len(pf):]
	}
	return SnakeToPascal(cn)
}

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

// model.go コード
func (gen *generator)codeModel(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	
	fields := ""
	for _, column := range table.Columns {
		cn := strings.ToLower(column.Name)
		fields += "\t" + gen.getFieldName(cn ,tn) + " ";
		if gen.isNullColumn(column, table.Constraints) {
			fields += "*"
		}
		fields += gen.dataTypeToGoType(column.DataType.Name) + " "
		fields += "`db:\"" + cn + "\" json:\"" + cn + "\"`\n"
	}
	fields = strings.TrimSuffix(fields, "\n")
	return fmt.Sprintf(
		MODEL_FORMAT, 
		tn, tnp, fields,
	)
}

// request.go生成
func (gen *generator) generateRequestFile(table ddlparse.Table, path string) error {
	code := codeRequest(table)
	err := WriteFile(fmt.Sprintf("%srequest.go", path), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// request.go コード
func (gen *generator)codeRequest(table ddlparse.Table) string {
	return fmt.Sprintf(
		REQUEST_FORMAT, 
		strings.ToLower(table.Name), 
		gen.codeRequestPostBodyFields(table), 
		gen.codeRequestPutBodyFields(table),
		gen.codeRequestDeleteBodyFields(table),
	)
	return code
}

func (gen *generator)codeRequestPostBodyFields(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	
	code := ""
	for _, column := range table.Columns {
		if !gen.isInsertColumn(column) {
			continue
		}
		cn := strings.ToLower(column.Name)
		code += "\t" + gen.getFieldName(cn ,tn) + " ";
		if gen.isNullColumn(column, table.Constraints) {
			code += "*"
		}
		code += gen.dataTypeToGoType(column.DataType.Name) + " "
		code += "`json:\"" + cn + "\""
		if gen.isNullColumn(column, table.Constraints) {
			code += " binding:\"required\"`\n"
		}
	}
	return strings.TrimSuffix(code, "\n")
}

// request.go コード
func (gen *generator)codeRequestPutBodyFields(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	
	code := ""
	for _, column := range table.Columns {
		if strings.Contains(c.Name, "_at") || strings.Contains(c.Name, "_AT") {
			return false
		}
		cn := strings.ToLower(column.Name)
		code += "\t" + gen.getFieldName(cn ,tn) + " ";
		if gen.isNullColumn(column, table.Constraints) {
			code += "*"
		}
		code += gen.dataTypeToGoType(column.DataType.Name) + " "
		code += "`json:\"" + cn + "\""
		if gen.isNullColumn(column, table.Constraints) {
			code += " binding:\"required\"`\n"
		}
	}
	return strings.TrimSuffix(code, "\n")
}

func (gen *generator)getPKColumns(table ddlparse.Table) []ddlparse.Column {
	pkcols := []string{}
	for _, pk := range table.Constraints.PrimaryKey {
		for _, name := range pk.ColumnNames {
			pkcols = append(pkcols, name)
		}
	}

	names := []string{}
	ret := []ddlparse.Column{}
	for _, c := range table.Columns {
		if c.Constraint.IsPrimaryKey || contains(pkcols, c.Name) || strings.Contains(strings.ToUpper(c.DataType.Name), "SERIAL"){
			if !contains(names, c.Name) {
				names = append(names, c.Name)
				ret = append(ret, c)
			}
		}
	}
	return ret
}

func (gen *generator)codeRequestDeleteBodyFields(table ddlparse.Table) string {
	code := ""
	for _, column := range gen.getPKColumns(table) {
		cn := strings.ToLower(column.Name)
		code += "\t" + gen.getFieldName(cn ,tn) + " ";
		if gen.isNullColumn(column, table.Constraints) {
			code += "*"
		}
		code += gen.dataTypeToGoType(column.DataType.Name) + " "
		code += "`json:\"" + cn + "\""
		if gen.isNullColumn(column, table.Constraints) {
			code += " binding:\"required\"`\n"
		}
	}
	return strings.TrimSuffix(code, "\n")
}

func (gen *generator)generateRepositoryFile(table ddlparse.Table, path string) error {
	code := codeRepository(table)
	err := WriteFile(fmt.Sprintf("%srepository.go", path), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

func (gen *generator)codeRepository(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	return fmt.Sprintf(
		REQPOSITORY_FORMAT,
		tn, tnp, 
		tni, tnp, tnp, 
		tni, tnp, tnp, 
		tni, tnp,  gen.codeRepositoryInsertReturnType(table),
		tni, tnp, 
		tni, tnp, 
		tnc, tnp, tnp, tnc,
		gen.codeRepositoryGet(table),
		gen.codeRepositoryGetOne(table),
		gen.codeRepositoryInsert(table),
		gen.codeRepositoryUpdate(table),
		gen.codeRepositoryDelete(table),
	)
}

func (gen *generator)codeRepositoryInsertReturnType(table ddlparse.Table) string {
	aiColumn, found := gen.getAutoIncrementColumn(table)
	if found {
		return fmt.Sprintf("(%s, error)", gen.dataTypeToGoType(aiColumn.DataType.Name))
	}
	return "error"
}

func (gen *generator)codeRepositoryGet(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	query := "\n\t`SELECT\n"
	for i, c := range table.Columns {
		if i == 0 {
			query += fmt.Sprintf("\t\t%s", c.Name)
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
		REQPOSITORY_FORMAT_GET,
		tnc, tni, tnp, tnp, tni, 
		query,
		tnp, tnp, tni, tnp,
		scan,
		tnp, tni,
	) 
}

func (gen *generator)codeRepositoryGetOne(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	query := "\n\t`SELECT\n"
	for i, c := range table.Columns {
		if i == 0 {
			query += fmt.Sprintf("\t\t%s", c.Name)
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
		REQPOSITORY_FORMAT_GETONE,
		tnc, tni, tnp, tnp, tnp, tni, 
		query,
		scan,
	) 
}


func (gen *generator)concatBindVariableWithCommas(bindCount int) string {
	bindVar := "?"
	if gen.rdbms == "postgres" {
		bindVar := fmt.Sprintf("$%d", n)
	}
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, bindVar)
	}
	return strings.Join(ls, ",")
}


func (gen *generator)gen.isInsertColumn(c ddlparse.Column) bool {
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


func (gen *generator)codeRepositoryInsert(table ddlparse.Table) string {
	_, found := gen.getAutoIncrementColumn(table)
	if found {
		if cf.DBDriver == "mysql" {
			return gen.codeRepositoryInsertAIMySQL(table)
		} else {
			return gen.codeRepositoryInsertNomal(table)
		}	
	}
	return gen.codeRepositoryInsertNomal(table)
}


func (gen *generator)codeRepositoryInsertNomal(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	query := fmt.Sprintf("\n\t`INSERT INTO %s (\n", tn)
	bindCount := 0
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			bindCount += 1
			if bindCount == 1 {
				query += fmt.Sprintf("\t\t%s", c.Name)
			} else {
				query += fmt.Sprintf("\n\t\t,%s", c.Name)
			}
		}	
	}
	query += fmt.Sprintf("\n\t ) VALUES(%s)`\n", gen.concatBindVariableWithCommas(bindCount))

	binds := "\n"
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
		}
	}
	binds += "\t"

	return fmt.Sprintf(
		REQPOSITORY_FORMAT_INSERT,
		tnc, tni, tnp,
		query,
		binds,
	) 
}


func (gen *generator)codeRepositoryInsertAI(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	aiColumn, _ := gen.getAutoIncrementColumn(table)
	aicn := strings.ToLower(aiColumn.Name)
	aicnc := SnakeToCamel(aicn)

	query := fmt.Sprintf("\n\t`INSERT INTO %s (\n", tn)
	bindCount := 0
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			bindCount += 1
			if bindCount == 1 {
				query += fmt.Sprintf("\t\t%s", c.Name)
			} else {
				query += fmt.Sprintf("\n\t\t,%s", c.Name)
			}
		}	
	}
	query += fmt.Sprintf("\n\t ) VALUES(%s)", gen.concatBindVariableWithCommas(bindCount))
	query += fmt.Sprintf("\n\t RETURNING %s`\n", aicn)

	binds := "\n"
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
		}
	}
	binds += "\t"

	return fmt.Sprintf(
		TEMPLATE_INSERT_AI,
		tnc, tni, tnp,
		query,
		binds,
		aicnc, aicnc, aicnc, aicnc,
	) 
}


func (gen *generator)codeRepositoryInsertAIMySQL(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)
	aiColumn, _ := gen.getAutoIncrementColumn(table)
	aicn := strings.ToLower(aiColumn.Name)
	aicnc := SnakeToCamel(aicn)

	query := fmt.Sprintf("\n\t`INSERT INTO %s (\n", tn)
	bindCount := 0
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			bindCount += 1
			if bindCount == 1 {
				query += fmt.Sprintf("\t\t%s", c.Name)
			} else {
				query += fmt.Sprintf("\n\t\t,%s", c.Name)
			}
		}	
	}
	query += fmt.Sprintf("\n\t ) VALUES(%s)`\n", gen.concatBindVariableWithCommas(bindCount))

	binds := "\n"
	for _, c := range table.Columns {
		if gen.isInsertColumn(c) {
			binds += fmt.Sprintf("\t\t%s.%s,\n", tni, gen.getFieldName(c.Name ,tn))
		}
	}
	binds += "\t"

	return fmt.Sprintf(
		TEMPLATE_INSERT_AI_MYSQL,
		tnc, tni, tnp,
		query,
		binds,
		aicnc, aicnc, aicnc, aicnc,
	) 
}

// service.go生成
func (gen *generator) generateServiceFile(table ddlparse.Table, path string) error {
	code := codeService(table)
	err := WriteFile(fmt.Sprintf("%sservice.go", path), code)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

func (gen *generator) codeService(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	return fmt.Sprintf(
		SERVICE_FORMAT, 
		tn, tnp, tnp, tnp, tnp, tnp, tnp,
		gen.codeServiceCreate(table),
		gen.codeServiceUpdate(table),
		tnp,
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

	fields := ""
	for i, column := range gen.getPKColumns(table) {
		if i != 0 {
			fields += ", "
		} 
		fn := gen.getFieldName(c.Name ,tn)
		fields += fmt.Sprintf("%s: input.%s", fn, fn)
	}

	return fmt.Sprintf(
		SERVICE_FORMAT_CREATE,
		tnp, tnp, tnp, tnp, fields,
	) 
}


func (gen *generator)codeServiceCreateAI(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)
	aiColumn, _ := gen.getAutoIncrementColumn(table)
	aicn := strings.ToLower(aiColumn.Name)
	aicnc := SnakeToCamel(aicn)
	fn := gen.getFieldName(aicn ,tn)

	return fmt.Sprintf(
		SERVICE_FORMAT_CREATE_AI,
		tnp, tnp, aicnc, tnp, tnp, fmt.Sprintf("%s: %s", fn, aicnc),
	) 
}

func (gen *generator)codeServiceUpdate(table ddlparse.Table) string {
	tn := strings.ToLower(table.Name)
	tnp := SnakeToPascal(tn)

	fields := ""
	for i, column := range gen.getPKColumns(table) {
		if i != 0 {
			fields += ", "
		} 
		fn := gen.getFieldName(c.Name ,tn)
		fields += fmt.Sprintf("%s: input.%s", fn, fn)
	}

	return fmt.Sprintf(
		SERVICE_FORMAT_UPDATE,
		tnp, tnp, tnp, tnp, fields,
	) 
}
/*
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
*/
//var JS_COMMON_CODE = ""
//const JS_FORMAT =
//`
///* <tr></tr>を作成 （tbody末尾の新規登録用レコード）*/
//const createTrNew = (elem) => {
//	return %s
//} 
//
///* <tr></tr>を作成 */
//const createTr = (elem) => {
//	return %s
//} 
//
//
///* セットアップ */
//const setUp = () => {
//	fetch('api/%s')
//	.then(response => response.json())
//	.then(data  => renderTbody(data))
//	.then(() => {%s
//	});
//}
//
//
///* 一括更新 */
//const doPutAll = async () => {
//	let successCount = 0;
//	let errorCount = 0;
//	%s
//
//	for (let i = 0; i < %s.length; i++) {
//		if (%s) {
//
//			let requestBody = {%s
//			}
//
//			await fetch('api/%s', {
//				method: 'PUT',
//				headers: {'Content-Type': 'application/json'},
//				body: JSON.stringify(requestBody)
//			})
//			.then(response => {
//				if (!response.ok){
//					throw new Error(response.statusText);
//				}
//  				return response.json();
//  			})
//			.then(data => {%s
//
//				successCount += 1;
//			}).catch(error => {
//				errorCount += 1;				
//			})
//		}
//	}
//
//	renderMessage('更新', successCount, true);
//	renderMessage('更新', errorCount, false);
//} 
//
//
///* 新規登録 */
//const doPost = () => {%s
//
//	if (%s)
//	{
//		let requestBody = {%s
//		}
//
//		fetch('api/%s', {
//			method: 'POST',
//			headers: {'Content-Type': 'application/json'},
//			body: JSON.stringify(requestBody)
//		})
//		.then(response => {
//			if (!response.ok){
//				throw new Error(response.statusText);
//			}
//  			return response.json();
//  		})
//		.then(data => {
//			document.getElementById('new').remove();
//
//			let tmpElem = document.createElement('tbody');
//			tmpElem.innerHTML = createTr(data);
//			tmpElem.firstChild.addEventListener('change', changeAction);
//			document.getElementById('records').appendChild(tmpElem.firstChild);
//
//			tmpElem = document.createElement('tbody');
//			tmpElem.innerHTML = createTrNew();
//			document.getElementById('records').appendChild(tmpElem.firstChild);
//
//			renderMessage('登録', 1, true);
//		}).catch(error => {
//			renderMessage('登録', 1, false);
//		})
//	}
//}
//
//
///* 一括削除 */
//const doDeleteAll = async () => {
//	let ls = getDeleteTarget();
//	let successCount = 0;
//	let errorCount = 0;
//
//	for (let x of ls) {
//		await fetch('api/%s', {
//			method: 'DELETE',
//			headers: {'Content-Type': 'application/json'},
//			body: x
//		})
//		.then(response => {
//			if (!response.ok){
//				throw new Error(response.statusText);
//			}
//			successCount += 1;
//  		}).catch(error => {
//			errorCount += 1;
//		});
//	}
//
//	setUp();
//
//	renderMessage('削除', successCount, true);
//	renderMessage('削除', errorCount, false);
//}
//`
//
//func generateHtmlCode_h2(table *dto.Table) string {
//	if table.TableNameJp != "" {
//		return fmt.Sprintf("%s（%s）", table.TableName, table.TableNameJp)
//	} else {
//		return table.TableName
//	}
//}
//
//func generateHtmlCode_tr(table *dto.Table) string {
//	code := ""
//	for _, col := range table.Columns {
//		if (col.IsNotNull || col.IsPrimaryKey) && (col.IsUpdAble || col.IsInsAble) {
//			code += fmt.Sprintf("\n\t\t\t\t<th>%s<spnn style='color:red;'>*</spnn></th>", col.ColumnName)
//		} else {
//			code += fmt.Sprintf("\n\t\t\t\t<th>%s</th>", col.ColumnName)
//		}
//	}
//	return code
//}
//
//const HTML_FORMAT =
//`
//<h2 class="ps-2 my-2">%s</h2>
//<hr class="mt-0 mb-2">
//<div class="w-100 vh-100 px-3">
//	<div id=message></div>
//	<button type="button" class="btn btn-danger" data-bs-toggle="modal" data-bs-target="#ModalDeleteAll">削除</button>
//	<button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#ModalSaveAll">保存</button>
//	<button type="button" class="btn btn-secondary" id="reload">リロード</button>
//	<div class="table-responsive mt-2" style="height:70%s">
//	<table class="table table-hover table-bordered table-sm">
//		<thead>
//			<tr class="fixed-table-header bg-light">
//				<th>削除</th>%s
//			</tr>
//		</thead>
//		<tbody id="records">
//		</tbody>
//	</table>
//	</div>
//</div>
//<script src="/static/js/%s.js"></script>
//`
//
//func generateHtmlCodeHeader_ul(tables *[]dto.Table) string {
//	code := ""
//	for _, table := range *tables {
//		tn := table.TableName
//		code += fmt.Sprintf("\n\t\t\t<li class='nav-item'><a href='/mastertables/%s' class='nav-link text-white'>%s</a></li>", tn, tn)
//	}
//	return code
//}
//
//const HTML_HEADER_FORMAT = 
//`<!DOCTYPE html>
//<html>
//<head>
//	<meta charset="utf-8">
//	<meta name=”description“ content=““ />
//	<meta name="viewport" content="width=device-width,initial-scale=1">
//	<link rel="stylesheet" href="/static/css/style.css">
//	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">
//	<title>マスタメンテナンス</title>
//</head>
//<body>
//
//<!-- 削除確認モーダル -->
//<div class="modal" tabindex="-1" id="ModalDeleteAll">
//<div class="modal-dialog">
//	<div class="modal-content">
//		<div class="modal-header">
//			<h4 class="modal-title">削除</h4>
//			<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
//		</div>
//		<div class="modal-body">
//			<p>この操作は元には戻せません。よろしいですか？</p>
//		</div>
//		<div class="modal-footer">
//			<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">キャンセル</button>
//			<button type="button" class="btn btn-danger" data-bs-dismiss="modal" id="ModalDeleteAllOk">削除</button>
//		</div>
//	</div>
//</div>
//</div>
//
//<!-- 保存確認モーダル -->
//<div class="modal" tabindex="-1" id="ModalSaveAll">
//<div class="modal-dialog">
//	<div class="modal-content">
//		<div class="modal-header">
//			<h4 class="modal-title">保存</h4>
//			<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
//		</div>
//		<div class="modal-body">
//			<p>この操作は元には戻せません。よろしいですか？</p>
//		</div>
//		<div class="modal-footer">
//			<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">キャンセル</button>
//			<button type="button" class="btn btn-primary" data-bs-dismiss="modal" id="ModalSaveAllOk">保存</button>
//		</div>
//	</div>
//</div>
//</div>
//
//<main>
//	<!-- サイドバー -->
//	<div class="d-flex flex-column flex-shrink-0 p-3 text-white bg-secondary" style="width: 280px;">
//		<span class="fs-4">テーブル一覧</span>
//		<hr class="mt-0 mb-2">
//		<ul class="nav nav-pills flex-column mb-auto">%s
//		</ul>
//	</div>
//	<!-- メインコンテンツ -->
//	<div class="w-100 vh-100">
//` 
//
//const HTML_FOOTER_CODE = 
//`
//	</div>
//	</div>
//</main>
//<footer>
//	Copyright &copy; kodaimurakami. 2023. 
//</footer>
//
//<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha1/dist/js/bootstrap.bundle.min.js" integrity="sha384-w76AqPfDkMBDXo30jS1Sgez6pr3x5MlQ1ZAGC+nuZB+EYdgRZgiwxhTBTkF7CXvN" crossorigin="anonymous"></script>
//</body>
//</html>
//`
//
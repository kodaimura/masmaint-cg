package generator

import (
	"io"
	"os"
	"fmt"
	"strings"
	"path/filepath"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
)


func WriteFile(path, content string) error {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logger.LogError(err.Error())
		return err
	}
	if _, err = f.WriteString(content); err != nil {
		logger.LogError(err.Error())
		return err
	}
	return nil
}


func CopyFile(source string, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}


func CopyDir(source string, destination string) error {
	err := os.MkdirAll(destination, 0755)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destinationPath := filepath.Join(destination, entry.Name())

		if entry.IsDir() {
			err := CopyDir(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		} else {
			err := CopyFile(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//xxx -> Xxx / xxx_yyy -> XxxYyy
func SnakeToPascal(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}

//xxx -> xxx / xxx_yyy -> xxxYyy
func SnakeToCamel(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		if i != 0 {
			ls[i] = strings.ToUpper(s[0:1]) + s[1:]
		}
	}
	return strings.Join(ls, "")
}

//xxx -> x / xxx_yyy -> xy
func GetSnakeInitial(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	ret := ""
	for _, s := range ls {
		ret = s[0:1]
	}
	return ret
}

func readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		logger.LogFatal(err.Error())
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.LogFatal(err.Error())
	}
	fileSize := fileInfo.Size()

	content := make([]byte, fileSize)

	_, err = file.Read(content)
	if err != nil {
		logger.LogFatal(err.Error())
	}

	return string(content)
}


// Jsコード生成
func GenerateJsCode(table *dto.Table) string {
	tn := table.TableName
	code := JS_COMMON_CODE + "\n"
	code += fmt.Sprintf(
		JS_FORMAT_CODE, 
		generateJsCode_createTrNew(table),
		generateJsCode_createTr(table),
		tn,
		generateJsCode_setUp(table),
		generateJsCode_doPutAll(table),
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

// Js共通部分
var JS_COMMON_CODE = readFile("_originalcopy_/js_common_code.js")
const JS_FORMAT_CODE =
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

	for (let i = 0; i < first_name.length; i++) {
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
	code := "`<tr><td><input class='form-check-input' type='checkbox' name='del' value=${JSON.stringify(elem)}></td>`"
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
		if col.IsUpdAble || col.IsPrimaryKey {
			cn := col.ColumnName
			code += fmt.Sprintf("\n\t\t\t\t%s: %s[i].value,", cn, cn)
		}
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

package generator

import (
	"fmt"
	"masmaint-cg/internal/shared/dto"
)


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

var JS_COMMON_CODE = ReadFile("_originalcopy_/js_common_code.js")
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

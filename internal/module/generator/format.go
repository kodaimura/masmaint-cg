package generator

const CONTROLLER_FORMAT =
`package %s

import (
	"github.com/gin-gonic/gin"
	"masmaint/internal/core/errs"
)

type controller struct {
	service Service
}

func NewController() *controller {
	service := NewService()
	return &controller{service}
}


//GET /%s
func (ctr *controller) GetPage(c *gin.Context) {
	c.HTML(200, "%s.html", gin.H{})
}


//GET /api/%s
func (ctr *controller) Get(c *gin.Context) {
	ret, err := ctr.service.Get()
	if err != nil {
		c.Error(errs.NewServiceError(err))
		return
	}

	c.JSON(200, ret)
}


//POST /api/%s
func (ctr *controller) Post(c *gin.Context) {
	var req PostBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBindError(err, &req))
		return
	}

	ret, err := ctr.service.Create(req)
	if err != nil {
		c.Error(errs.NewServiceError(err))
		return
	}

	c.JSON(200, ret)
}


//PUT /api/%s
func (ctr *controller) Put(c *gin.Context) {
	var req PutBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBindError(err, &req))
		return
	}

	ret, err := ctr.service.Update(req)
	if err != nil {
		c.Error(errs.NewServiceError(err))
		return
	}

	c.JSON(200, ret)
}


//DELETE /api/%s
func (ctr *controller) Delete(c *gin.Context) {
	var req DeleteBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewBindError(err, &req))
		return
	}

	if err := ctr.service.Delete(req); err != nil {
		c.Error(errs.NewServiceError(err))
		return
	}

	c.JSON(200, gin.H{})
}`


const MODEL_FORMAT = `
package %s

type %s struct {
%s
}
`

const REQUEST_FORMAT = `
package %s

type PostBody struct {
%s
}

type PutBody struct {
%s
}

type DeleteBody struct {
%s
}
`


const REQPOSITORY_FORMAT = 
`package %s

import (
	"database/sql"
	"masmaint/internal/core/db"
)


type %sRepository interface {
	Get(%s *model.%s) ([]model.%s, error)
	GetOne(%s *model.%s) (model.%s, error)
	Insert(%s *model.%s, tx *sql.Tx) %s
	Update(%s *model.%s, tx *sql.Tx) error
	Delete(%s *model.%s, tx *sql.Tx) error
}


type %sRepository struct {
	db *sql.DB
}

func New%sRepository() %sRepository {
	db := db.GetDB()
	return &%sRepository{db}
}


%s


%s


%s


%s


%s`

const REQPOSITORY_FORMAT_GET =
`func (rep *%sRepository) Get(%s *model.%s) ([]model.%s, error) {
	where, binds := db.BuildWhereClause(%s)
	query := %s + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.%s{}, err
	}

	ret := []model.%s{}
	for rows.Next() {
		%s := model.%s{}
		err = rows.Scan(%s)
		if err != nil {
			return []model.%s{}, err
		}
		ret = append(ret, %s)
	}

	return ret, nil
}`

const REQPOSITORY_FORMAT_GETONE =
`func (rep *%sRepository) GetOne(%s *model.%s) (model.%s, error) {
	var ret model.%s
	where, binds := db.BuildWhereClause(%s)
	query := %s + where

	err := rep.db.QueryRow(query, binds...).Scan(%s)

	return ret, err
}`

const REQPOSITORY_FORMAT_INSERT =
`func (rep *%sRepository) Insert(%s *model.%s, tx *sql.Tx) error {
	cmd := %s
	binds := []interface{}{%s}

	var err error
	if tx != nil {
		_, err = tx.Exec(cmd, binds...)
	} else {
		_, err = rep.db.Exec(cmd, binds...)
	}

	return err
}`

const REQPOSITORY_FORMAT_INSERT_AI =
`func (rep *%sRepository) Insert(%s *model.%s, tx *sql.Tx) (int, error) {
	cmd := %s
	binds := []interface{}{%s}

	var %s int
	var err error
	if tx != nil {
		err = tx.QueryRow(cmd, binds...).Scan(&%s)
	} else {
		err = rep.db.QueryRow(cmd, binds...).Scan(&%s)
	}

	return %s, err
}`

const REQPOSITORY_FORMAT_INSERT_AI_MYSQL =
`func (rep *%sRepository) Insert(%s *model.%s, tx *sql.Tx) (int, error) {
	cmd := %s
	binds := []interface{}{%s}

	var err error
	if tx != nil {
		_, err = tx.Exec(cmd, binds...)
	} else {
		_, err = rep.db.Exec(cmd, binds...)
	}

	if err != nil {
		return 0, err
	}

	var %s int
	if tx != nil {
		err = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&%s)
	} else {
		err = rep.db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&%s)
	}

	return %s, err
}`

const REQPOSITORY_FORMAT_UPDATE =
`func (rep *%sRepository) Update(%s *model.%s, tx *sql.Tx) error {
	cmd := %s
	binds := []interface{}{%s}

	var err error
	if tx != nil {
		_, err = tx.Exec(cmd, binds...)
	} else {
		_, err = rep.db.Exec(cmd, binds...)
	}

	return err
}`

const REQPOSITORY_FORMAT_DELETE =
`func (rep *%sRepository) Delete(%s *model.%s, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(%s)
	cmd := "DELETE FROM %s " + where

	var err error
	if tx != nil {
		_, err = tx.Exec(cmd, binds...)
	} else {
		_, err = rep.db.Exec(cmd, binds...)
	}

	return err
}`

const SERVICE_FORMAT =
`package %s

import (
	"masmaint/internal/core/logger"
	"masmaint/internal/core/utils"
)

type Service interface {
	Get() ([]%s, error)
	Create(input PostBody) (%s, error)
	Update(input PutBody) (%s, error)
	Delete(input DeleteBody) error
}

type service struct {
	repository Repository
}

func NewService() Service {
	return &service{
		repository: NewRepository(),
	}
}


%s


%s


%s


%s`

const SERVICE_FORMAT_GET =
`func (srv *service) Get() ([]%s, error) {
	rows, err := srv.repository.Get(&%s{})
	if err != nil {
		logger.Error(err.Error())
		return []%s{}, err
	}
	return rows, nil
}`

const SERVICE_FORMAT_CREATE =
`func (srv *service) Create(input PostBody) (%s, error) {
	var model %s
	utils.MapFields(&model, input)

	err := srv.repository.Insert(&model, nil)
	if err != nil {
		logger.Error(err.Error())
		return %s{}, err
	}

	return srv.repository.GetOne(&%s{ %s })
}`

const SERVICE_FORMAT_CREATE_AI =
`func (srv *service) Create(input PostBody) (%s, error) {
	var model %s
	utils.MapFields(&model, input)

	%s, err := srv.repository.Insert(&model, nil)
	if err != nil {
		logger.Error(err.Error())
		return %s{}, err
	}

	return srv.repository.GetOne(&%s{ %s })
}`

const SERVICE_FORMAT_UPDATE =
`func (srv *service) Update(input PutBody) (%s, error) {
	var model %s
	utils.MapFields(&model, input)

	err := srv.repository.Update(&model, nil)
	if err != nil {
		logger.Error(err.Error())
		return %s{}, err
	}

	return srv.repository.GetOne(&%s{ %s })
}`

const SERVICE_FORMAT_DELETE =
`func (srv *service) Delete(input DeleteBody) error {
	var model %s
	utils.MapFields(&model, input)

	err := srv.repository.Delete(&model, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}`

var JS_FORMAT = ReadFile("_template/js_format.txt")

const JS_FORMAT_CREATETRNEW =
`const createTrNew = (elem) => {
	const tr = document.createElement('tr');
	tr.id = 'new';
	tr.innerHTML = `+"`%s`"+`;
	return tr;
}`

const JS_FORMAT_CREATETR =
`const createTr = (elem) => {
	const tr = document.createElement('tr');
	tr.innerHTML = `+"`%s`"+`;
	return tr;
}`

const JS_FORMAT_GETROWS =
`const getRows = async () => {
	document.getElementById('records').innerHTML = '';
	const rows = await api.get('%s');
	renderTbody(rows);
%s
}`

const JS_FORMAT_PUTROWS =
`const putRows = async () => {
	let successCount = 0;
	let errorCount = 0;

%s

%s

	for (let i = 0; i < code.length; i++) {
		const rowMap = {
%s
		}

		const rowBkMap = {
%s
		}

		//差分がある行のみ更新
		if (Object.keys(rowMap).some(key => rowMap[key].value !== rowBkMap[key].value)) {
			const requestBody = {
%s
			}

			try {
				const data = await api.put('%s', requestBody);

%s

				Object.values(rowMap).forEach(element => {
					element.classList.remove('changed');
					element.classList.remove('error');
				});

				successCount += 1;
			} catch (e) {
				Object.keys(rowMap).forEach(key => {
					rowMap[key].classList.toggle('error', key === e.details.field);
				});
				errorCount += 1;
			}
		}
	}

	renderMessage('更新', successCount, true);
	renderMessage('更新', errorCount, false);
}`

const JS_FORMAT_POSTROW =
`const postRow = async () => {
	const rowMap = {
%s
	}

	if (Object.keys(rowMap).some(key => rowMap[key].value !== '')) {
		const requestBody = {
%s
		}

		try {
			const data = await api.post('%s', requestBody);

			document.getElementById('new').remove();
			const tr = createTr(data);
			tr.addEventListener('change', handleChange);
			document.getElementById('records').appendChild(tr);
			document.getElementById('records').appendChild(createTrNew());

			renderMessage('登録', 1, true);
		} catch (e) {
			Object.keys(rowMap).forEach(key => {
				rowMap[key].classList.toggle('error', key === e.details.field || `+"%s.${key}"+` === e.details.column);
			});
			renderMessage('登録', 1, false);
		}
	}
}`

const JS_FORMAT_DELETEROWS =
`const deleteRows = async () => {
	const rows = getDeleteTargetRows();
	let successCount = 0;
	let errorCount = 0;

	for (let row of rows) {
		try {
			await api.delete('%s', row);
			successCount += 1;
		} catch (e) {
			errorCount += 1;
		}
	}

	getRows();

	renderMessage('削除', successCount, true);
	renderMessage('削除', errorCount, false);
}`

const TEMPLATE_FORMAT =
`<!DOCTYPE html>
<html>

<head>
	{{template "head" .}}
</head>

<body>
	{{template "header" .}}
	<div class="container-fluid">
		{{template "menu" .}}
		<main>
			<div class="w-100 px-3 py-3">
				<h1 class="h4">%s</h1>
				<div id="message"></div>
				<button type="button" class="btn btn-danger" data-bs-toggle="modal"
					data-bs-target="#modal-delete">削除</button>
				<button type="button" class="btn btn-primary" data-bs-toggle="modal"
					data-bs-target="#modal-save">保存</button>
				<button type="button" class="btn btn-secondary" id="reload">リロード</button>
				<div class="table-responsive mt-2" style="max-height: calc(100vh - 190px);">
					<table class="table table-hover table-bordered table-sm">
						<thead class="fixed-table-header bg-light">
							<tr>
								<th>削除</th>
%s
							</tr>
						</thead>
						<tbody id="records">
						</tbody>
					</table>
				</div>
			</div>
		</main>
	</div>
	{{template "modal" .}}
	<script type="module" src="js/%s.js"></script>
	{{template "footer" .}}
</body>

</html>`

const TEMPLATE_FORMAT_MENU =
`{{define "menu"}}
<div class="sidemenu vh-100" style="overflow-y: auto;">
	<ul class="nav flex-column mb-5">
%s
	</ul>
</div>
{{end}}`
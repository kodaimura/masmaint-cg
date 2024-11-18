package generator

const FORMAT_CONTROLLER =
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


const FORMAT_MODEL = `
package %s

type %s struct {
%s
}
`

const FORMAT_REQUEST = `
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


const FORMAT_REPOSITORY = 
`package %s

import (
	"database/sql"
	"masmaint/internal/core/db"
)


%s


type repository struct {
	db *sql.DB
}

func NewRepository() Repository {
	db := db.GetDB()
	return &repository{db}
}


%s


%s


%s


%s


%s`

const FORMAT_REPOSITORY_INTERFACE =
`type Repository interface {
	Get(%s *%s) ([]%s, error)
	GetOne(%s *%s) (%s, error)
	Insert(%s *%s, tx *sql.Tx) %s
	Update(%s *%s, tx *sql.Tx) error
	Delete(%s *%s, tx *sql.Tx) error
}`

const FORMAT_REPOSITORY_GET =
`func (rep *repository) Get(%s *%s) ([]%s, error) {
	where, binds := db.BuildWhereClause(%s)
	query := %s + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []%s{}, err
	}

	ret := []%s{}
	for rows.Next() {
		%s := %s{}
		err = rows.Scan(%s)
		if err != nil {
			return []%s{}, err
		}
		ret = append(ret, %s)
	}

	return ret, nil
}`

const FORMAT_REPOSITORY_GETONE =
`func (rep *repository) GetOne(%s *%s) (%s, error) {
	var ret %s
	where, binds := db.BuildWhereClause(%s)
	query := %s + where

	err := rep.db.QueryRow(query, binds...).Scan(%s)

	return ret, err
}`

const FORMAT_REPOSITORY_INSERT =
`func (rep *repository) Insert(%s *%s, tx *sql.Tx) error {
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

const FORMAT_REPOSITORY_INSERT_AI =
`func (rep *repository) Insert(%s *%s, tx *sql.Tx) (int, error) {
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

const FORMAT_REPOSITORY_INSERT_AI_MYSQL =
`func (rep *repository) Insert(%s *%s, tx *sql.Tx) (int, error) {
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

const FORMAT_REPOSITORY_UPDATE =
`func (rep *repository) Update(%s *%s, tx *sql.Tx) error {
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

const FORMAT_REPOSITORY_DELETE =
`func (rep *repository) Delete(%s *%s, tx *sql.Tx) error {
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

const FORMAT_SERVICE =
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

const FORMAT_SERVICE_GET =
`func (srv *service) Get() ([]%s, error) {
	rows, err := srv.repository.Get(&%s{})
	if err != nil {
		logger.Error(err.Error())
		return []%s{}, err
	}
	return rows, nil
}`

const FORMAT_SERVICE_CREATE =
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

const FORMAT_SERVICE_CREATE_AI =
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

const FORMAT_SERVICE_UPDATE =
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

const FORMAT_SERVICE_DELETE =
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

const FORMAT_ROUTER =
`package server

import (
	"github.com/gin-gonic/gin"
	"masmaint/config"
	"masmaint/internal/core/jwt"
	"masmaint/internal/middleware"

%s
)

/*
 Routing for "/" 
*/
%s


%s`

const FORMAT_ROUTER_SETWEB =
`func SetWebRouter(r *gin.RouterGroup) {
%s

	r.GET("/login", func(c *gin.Context) { c.HTML(200, "login.html", gin.H{}) })

	auth := r.Group("", middleware.JwtAuthMiddleware())
	{
		auth.GET("/", func(c *gin.Context) { c.HTML(200, "index.html", gin.H{}) })
%s
	}
}`

const FORMAT_ROUTER_SETAPI =
`func SetApiRouter(r *gin.RouterGroup) {
	r.Use(middleware.ApiResponseMiddleware())

%s

	//カスタム推奨
	r.POST("/login", func(c *gin.Context) { 
		var body map[string]string
		c.ShouldBindJSON(&body)
		name := body["username"]
		pass := body["password"]

		cf := config.GetConfig()
		if name == cf.AuthUser && pass == cf.AuthPass {
			cc := jwt.CustomClaims{ AccountId: 1, AccountName: name}
			jwt.SetTokenToCookie(c, jwt.NewPayload(cc))
		} else {
			c.JSON(401, gin.H{"error": "ユーザ名またはパスワードが異なります。"})
		}
	})

	auth := r.Group("", middleware.JwtAuthApiMiddleware())
	{
%s
	}
}`

var FORMAT_JS = ReadFile("_template/js_format.txt")

const FORMAT_JS_CREATETRNEW =
`const createTrNew = (elem) => {
	const tr = document.createElement('tr');
	tr.id = 'new';
	tr.innerHTML = `+"`%s`"+`;
	return tr;
}`

const FORMAT_JS_CREATETR =
`const createTr = (elem) => {
	const tr = document.createElement('tr');
	tr.innerHTML = `+"`%s`"+`;
	return tr;
}`

const FORMAT_JS_GETROWS =
`const getRows = async () => {
	document.getElementById('records').innerHTML = '';
	const rows = await api.get('%s');
	renderTbody(rows);
%s
}`

const FORMAT_JS_PUTROWS =
`const putRows = async () => {
	let successCount = 0;
	let errorCount = 0;

%s

%s

	for (let i = 0; i < %s.length; i++) {
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

const FORMAT_JS_POSTROW =
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
				rowMap[key].classList.toggle('error', key === e.details.field || `+"`%s.${key}`"+` === e.details.column);
			});
			renderMessage('登録', 1, false);
		}
	}
}`

const FORMAT_JS_DELETEROWS =
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

const FORMAT_TEMPLATE =
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

const FORMAT_TEMPLATE_MENU =
`{{define "menu"}}
<div class="sidemenu vh-100" style="overflow-y: auto;">
	<ul class="nav flex-column mb-5">
%s
	</ul>
</div>
{{end}}`
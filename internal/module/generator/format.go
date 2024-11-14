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
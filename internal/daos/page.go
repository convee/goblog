package daos

import (
	"database/sql"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/pkg/logger"
	"go.uber.org/zap"
	"strconv"
)

type PageParams struct {
	PerPage int
	Page    int
}

func GetPages(params PageParams) (pages []model.Page, err error) {
	querySql := "select id,title,ident,content from page"

	querySql += " order by id asc"

	if params.PerPage > 0 && params.Page > 0 {
		offset := (params.Page - 1) * params.PerPage
		querySql += " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(params.PerPage)
	}
	rows, err := db.Query(querySql)
	if err != nil {
		logger.Error("GetPagesError", zap.Error(err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var page model.Page
		if err = rows.Scan(&page.Id, &page.Title, &page.Ident, &page.Content); err != nil {
			logger.Error("GetPagesError", zap.Error(err))
			return nil, err
		}
		pages = append(pages, page)
	}
	if err = rows.Err(); err != nil {
		logger.Error("GetPagesError", zap.Error(err))
	}
	return pages, err
}

func GetPage(id string) (page model.Page, err error) {
	row := db.QueryRow("select id,title,ident,content from page  where id=? limit 1", id)
	err = row.Scan(&page.Id, &page.Title, &page.Ident, &page.Content)
	if err != nil {
		logger.Error("GetPageError", zap.Error(err))
		return page, err
	}
	return page, nil
}

func GetPageByIdent(ident string) (page model.Page, err error) {
	row := db.QueryRow("select id,title,content from page  where ident=? limit 1", ident)
	err = row.Scan(&page.Id, &page.Title, &page.Content)
	if err != nil {
		logger.Error("GetPageError", zap.Error(err))
		return page, err
	}
	return page, nil
}

func PageDelete(page model.Page) (id int, err error) {
	id = page.Id
	_, err = db.Exec("delete from page where id=?", id)
	if err != nil {
		logger.Error("PageDeleteError", zap.Error(err))
		return
	}
	return
}
func PageSave(page model.Page) (id int, err error) {

	var rs sql.Result
	if page.Id > 0 {
		id = page.Id
		_, err = db.Exec("update page set title=?,content=?,ident=? where id=?", page.Title, page.Content, page.Ident, id)
		if err != nil {
			logger.Error("PageSaveError", zap.Error(err))
			return
		}
	} else {
		rs, err = db.Exec("insert into page (title, ident, content) values (?,?,?)", page.Title, page.Ident, page.Content)
		if err != nil {
			logger.Error("PageSaveError", zap.Error(err))
			return
		}
		id64, lastIDErr := rs.LastInsertId()
		if lastIDErr != nil {
			return 0, lastIDErr
		}
		id = int(id64)

	}
	return
}

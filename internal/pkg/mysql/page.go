package mysql

import (
	"database/sql"
	"github.com/convee/goblog/internal/pkg/model"
	"log"
	"strconv"
)

type PageParams struct {
	PerPage int
	Page    int
}

func GetPages(params PageParams) (pages []model.Page, err error) {
	querySql := "select id,title,content from page"

	querySql += " order by id asc"

	if params.PerPage > 0 && params.Page > 0 {
		offset := (params.Page - 1) * params.PerPage
		querySql += " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(params.PerPage)
	}
	rows, err := db.Query(querySql)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var page model.Page
		rows.Scan(&page.Id, &page.Title, &page.Content)
		pages = append(pages, page)
	}
	return
}

func GetPage(id string) (page model.Page) {
	row := db.QueryRow("select id,title,content from page  where id=? limit 1", id)
	err := row.Scan(&page.Id, &page.Title, &page.Content)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func PageDelete(page model.Page) (id int, err error) {
	var rs sql.Result
	id = page.Id
	log.Println(page)
	rs, err = db.Exec("delete from page where id=?", id)
	if err != nil {
		log.Printf("page %d delete err %v", id, err)
		return
	}
	affected, _ := rs.RowsAffected()
	log.Printf("page %d delete success affected:%d", id, affected)
	return
}
func PageSave(page model.Page) (id int, err error) {

	var rs sql.Result
	if page.Id > 0 {
		id = page.Id
		log.Println(page)
		rs, err = db.Exec("update page set title=?,content=? where id=?", page.Title, page.Content, id)
		if err != nil {
			log.Printf("page %d update err %v", id, err)
			return
		}
		affected, _ := rs.RowsAffected()
		log.Printf("page %d save success affected:%d", page.Id, affected)
	} else {
		rs, err = db.Exec("insert into page (title, content) values (?,?)", page.Title, page.Content)
		if err != nil {
			log.Printf("page %d insert err %v", id, err)
			return
		}
		id64, _ := rs.LastInsertId()
		id = int(id64)

	}
	return
}

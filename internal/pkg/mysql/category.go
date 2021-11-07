package mysql

import (
	"database/sql"
	"log"

	"github.com/convee/goblog/internal/pkg/model"
)

func GetCategories() (categories []model.Category, err error) {
	rs, err := db.Query("select id,name from category")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var category model.Category
		rs.Scan(&category.Id, &category.Name)
		categories = append(categories, category)
	}
	return
}

func GetCategory(id int) (category model.Category) {
	row := db.QueryRow("select id,name from category where id=?", id)
	row.Scan(&category.Id, &category.Name)
	return
}

func CategoryDelete(category model.Category) (id int, err error) {
	var rs sql.Result
	id = category.Id
	log.Println(category)
	rs, err = db.Exec("delete from category where id=?", id)
	if err != nil {
		log.Printf("category %d delete err %v", id, err)
		return
	}
	affected, _ := rs.RowsAffected()
	log.Printf("category %d delete success affected:%d", id, affected)
	return
}

func CategorySave(category model.Category) (id int, err error) {

	var rs sql.Result
	if category.Id > 0 {
		id = category.Id
		log.Println(category)
		rs, err = db.Exec("update category set name=? where id=?", category.Name, id)
		if err != nil {
			log.Printf("category %d update err %v", id, err)
			return
		}
		affected, _ := rs.RowsAffected()
		log.Printf("category %d save success affected:%d", category.Id, affected)
	} else {
		rs, err = db.Exec("insert into category (`name`) values (?)", category.Name)
		if err != nil {
			log.Printf("category %d insert err %v", id, err)
			return
		}
		id64, _ := rs.LastInsertId()
		id = int(id64)

	}
	return
}

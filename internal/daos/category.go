package daos

import (
	"database/sql"
	"errors"
	"github.com/convee/goblog/internal/model"
	"log"
)

var ErrCategoryInUse = errors.New("category is used by posts")

func GetCategories() (categories []model.Category, err error) {
	rs, err := db.Query("select id,name from category")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var category model.Category
		if err := rs.Scan(&category.Id, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rs.Err()
}

func GetCategoryIdsByName(name string) (categoryIds []string, err error) {
	rs, err := db.Query("select id from category where name like ?", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var categoryId string
		if err := rs.Scan(&categoryId); err != nil {
			return nil, err
		}
		categoryIds = append(categoryIds, categoryId)
	}
	return categoryIds, rs.Err()
}

func GetCategory(id int) (category model.Category, err error) {
	row := db.QueryRow("select id,name from category where id=?", id)
	err = row.Scan(&category.Id, &category.Name)
	return category, err
}

func CategoryDelete(category model.Category) (id int, err error) {
	id = category.Id
	var postCount int
	if err = db.QueryRow("select count(*) from post where category_id=?", id).Scan(&postCount); err != nil {
		return 0, err
	}
	if postCount > 0 {
		return 0, ErrCategoryInUse
	}
	_, err = db.Exec("delete from category where id=?", id)
	if err != nil {
		log.Printf("category %d delete err %v", id, err)
		return
	}
	return
}

func CategorySave(category model.Category) (id int, err error) {

	var rs sql.Result
	if category.Id > 0 {
		id = category.Id
		_, err = db.Exec("update category set name=? where id=?", category.Name, id)
		if err != nil {
			log.Printf("category %d update err %v", id, err)
			return
		}
	} else {
		rs, err = db.Exec("insert into category (`name`) values (?)", category.Name)
		if err != nil {
			log.Printf("category %d insert err %v", id, err)
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

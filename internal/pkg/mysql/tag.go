package mysql

import "github.com/convee/goblog/internal/pkg/model"

func GetTags() (tags []model.Tag, err error) {
	rows, err := db.Query("select id,name,count from tag order by count desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag model.Tag
		rows.Scan(&tag.Id, &tag.Name, &tag.Count)
		tags = append(tags, tag)
	}
	return
}

func AddTag(tag model.Tag) (id int, err error) {
	rs, err := db.Exec("insert into tag (name) values (?)", tag.Name)
	if err != nil {
		return
	}
	id64, err := rs.LastInsertId()
	return int(id64), err
}

func IncrTagCount(id string) {
	db.Exec("update tag set count=count+1 where id = ?", id)
}

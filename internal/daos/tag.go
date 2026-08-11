package daos

import (
	"github.com/convee/goblog/internal/model"
	"strconv"
)

func GetTags() (tags []model.Tag, err error) {
	rows, err := db.Query("select id,name,count from tag order by count desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.Id, &tag.Name, &tag.Count); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func GetTagIdsByName(name string) (tagIds []string, err error) {
	rs, err := db.Query("select id from tag where name like ?", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var tagId string
		if err := rs.Scan(&tagId); err != nil {
			return nil, err
		}
		tagIds = append(tagIds, tagId)
	}
	return tagIds, rs.Err()
}

func GetTagByName(name string) (tag model.Tag, err error) {
	err = db.QueryRow("select id,name,count from tag where name=? limit 1", name).Scan(&tag.Id, &tag.Name, &tag.Count)
	return tag, err
}

func AddTag(tag model.Tag) (id int, err error) {
	rs, err := db.Exec("insert into tag (name) values (?)", tag.Name)
	if err != nil {
		return
	}
	id64, err := rs.LastInsertId()
	return int(id64), err
}

func RefreshTagCount(id int) error {
	count, err := GetPostCountByTagID(strconv.Itoa(id))
	if err != nil {
		return err
	}
	_, err = db.Exec("update tag set count=? where id=?", count, id)
	return err
}

package mysql

import (
	"github.com/convee/goblog/internal/pkg/model"
	logger "github.com/convee/goblog/pkg/log"
	"go.uber.org/zap"
)

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

func GetTagIdsByName(name string) (tagIds []string, err error) {
	rs, err := db.Query("select id from tag where name like ?", "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var tagId string
		rs.Scan(&tagId)
		tagIds = append(tagIds, tagId)
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
	count, err := GetPostCountByTagId(id)
	if err != nil {
		logger.Error("incr_tag_count_err", zap.Error(err))
		return
	}
	db.Exec("update tag set count=? where id = ?", count, id)
}

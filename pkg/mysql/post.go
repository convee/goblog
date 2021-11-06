package mysql

import (
	"blog/pkg/model"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type PostParams struct {
	Ids []string
	CategoryId string
	TagId      string
	PerPage    int
	Page       int
}

func GetPosts(params PostParams) (posts []model.Post, err error) {
	var condition []string

	if len(params.CategoryId) > 0 {
		condition = append(condition, "category_id="+params.CategoryId)
	}
	if len(params.Ids) > 0 {
		condition = append(condition, "id in"+ strings.Join(params.Ids, ","))
	}
	if len(params.TagId) > 0 {
		condition = append(condition, "JSON_CONTAINS(tag_ids,'"+params.TagId+"')")
	}

	querySql := "select id,title,created_at,updated_at,category_id,tag_ids,views,description from post"

	condition = append(condition, "status=1")
	if len(condition) > 0 {
		querySql += ` where ` + strings.Join(condition, " and ")
	}

	querySql += " order by updated_at desc"

	if params.PerPage > 0 && params.Page > 0 {
		offset := (params.Page - 1) * params.PerPage
		querySql += " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(params.PerPage)
	}
	rows, err := db.Query(querySql)
	if err != nil {
		return
	}
	var tags []byte
	defer rows.Close()
	for rows.Next() {
		var post model.Post
		rows.Scan(&post.Id, &post.Title, &post.CreatedAt, &post.UpdatedAt, &post.CategoryId, &tags, &post.Views, &post.Description)
		json.Unmarshal(tags, &post.TagIds)
		posts = append(posts, post)
	}
	return
}

func GetPost(id string) (post model.Post) {
	var tags []byte
	row := db.QueryRow("select post.id,post.title,post.created_at,post.updated_at,post.category_id,post.tag_ids,post.views,post_content.content,post.description from post left join post_content on post.id=post_content.post_id where post.id=? limit 1", id)
	err := row.Scan(&post.Id, &post.Title, &post.CreatedAt, &post.UpdatedAt, &post.CategoryId, &tags, &post.Views, &post.Content, &post.Description)
	if err != nil {
		log.Println("get post err ", err)
		return
	}
	json.Unmarshal(tags, &post.TagIds)
	return
}

func IncrView(id string) {
	db.Exec("update post set views=views+1 where id = ?", id)
}

func PostDelete(post model.Post) (id int, err error) {
	conn, err := db.Begin()
	if err != nil {
		log.Fatalln("post db conn err ", err)
	}
	var rs sql.Result
	id = post.Id
	log.Println(post)
	rs, err = conn.Exec("delete from post where id=?", id)
	if err != nil {
		log.Printf("post %d delete err %v", id, err)
		conn.Rollback()
		return
	}
	affected ,_ := rs.RowsAffected()
	log.Printf("post %d delete success affected:%d", id, affected)
	rs, err = conn.Exec("delete from post_content where post_id=?", id)
	if err != nil {
		log.Printf("post_content %d delete err %v content:", id, err)
		conn.Rollback()
		return
	}
	affected ,_ = rs.RowsAffected()
	log.Printf("post_content %d delete success affected:%d", id, affected)

	conn.Commit()
	return
}
func PostSave(post model.Post) (id int, err error) {
	conn, err := db.Begin()
	if err != nil {
		log.Fatalln("post db conn err ", err)
	}
	var tagIds []byte
	tagIds, _ = json.Marshal(post.TagIds)
	var rs sql.Result
	if post.Id > 0 {
		id = post.Id
		log.Println(post)
		rs, err = conn.Exec("update post set title=?,description=?,category_id=?,tag_ids=? where id=?", post.Title, post.Description, post.CategoryId, tagIds, id)
		if err != nil {
			log.Printf("post %d update err %v", id, err)
			conn.Rollback()
			return
		}
		affected ,_ := rs.RowsAffected()
		log.Printf("post %d save success affected:%d", post.Id, affected)
		rs, err = conn.Exec("update post_content set content=? where post_id=?", post.Content, id)
		if err != nil {
			log.Printf("post_content %d update err %v content:", id, err)
			conn.Rollback()
			return
		}
		affected ,_ = rs.RowsAffected()
		log.Printf("post_content %d save success affected:%d", post.Id, affected)
	} else {
		rs, err = conn.Exec("insert into post (title, description,category_id,tag_ids,status) values (?,?,?,?)", post.Title, post.Description, post.CategoryId, tagIds, post.Status)
		if err != nil {
			log.Printf("post %d insert err %v", id, err)
			conn.Rollback()
			return
		}
		id64, _ := rs.LastInsertId()
		id = int(id64)
		_, err = conn.Exec("insert into post_content (post_id,content) values (?,?)", id, post.Content)
		if err != nil {
			log.Printf("post_content %d insert err %v content:", id, err)
			conn.Rollback()
			return
		}
	}
	conn.Commit()
	return
}

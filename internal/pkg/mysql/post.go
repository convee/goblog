package mysql

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/convee/goblog/internal/pkg/model"
)

type PostParams struct {
	Ids        map[string][]string
	CategoryId string
	TagId      string
	PerPage    int
	Page       int
	Keyword    string
}

func GetPosts(params PostParams) (posts []model.Post, err error) {
	var condition []string
	var args []interface{}

	if len(params.CategoryId) > 0 {
		condition = append(condition, "category_id=?")
		args = append(args, params.CategoryId)
	}

	if len(params.TagId) > 0 {
		condition = append(condition, "JSON_CONTAINS(tag_ids,?)")
		args = append(args, params.TagId)
	}
	if len(params.Keyword) > 0 {
		keywordSql := "title like ? or description like ?"
		args = append(args, "%"+params.Keyword+"%", "%"+params.Keyword+"%")
		if len(params.Ids["ids"]) > 0 {
			keywordSql += " or id in (?)"
			args = append(args, strings.Join(params.Ids["ids"], ","))
		}
		if len(params.Ids["category_ids"]) > 0 {
			keywordSql += " or category_id in (?)"
			args = append(args, strings.Join(params.Ids["category_ids"], ","))
		}
		if len(params.Ids["tag_ids"]) > 0 {
			for _, tagId := range params.Ids["tag_ids"] {
				keywordSql += " or JSON_CONTAINS(tag_ids,?)"
				args = append(args, tagId)
			}
		}
		condition = append(condition, "("+keywordSql+")")
	}

	querySql := "select id,title,created_at,updated_at,category_id,tag_ids,views,description from post"

	condition = append(condition, "status=1")
	if len(condition) > 0 {
		querySql += ` where ` + strings.Join(condition, " and ")
	}

	querySql += " order by is_top desc,updated_at desc"

	if params.PerPage > 0 && params.Page > 0 {
		offset := (params.Page - 1) * params.PerPage
		querySql += " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(params.PerPage)
	}
	rows, err := db.Query(querySql, args...)
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
	_, err := db.Exec("update post set views=views+1 where id = ?", id)
	if err != nil {
		log.Printf("update post views err id:%s, err%v", id, err)
		return
	}
}

func PostDelete(post model.Post) (id int, err error) {
	conn, err := db.Begin()
	if err != nil {
		log.Fatalln("post db conn err ", err)
		return
	}
	id = post.Id
	_, err = conn.Exec("delete from post where id=?", id)
	if err != nil {
		log.Printf("post %d delete err %v", id, err)
		conn.Rollback()
		return
	}
	_, err = conn.Exec("delete from post_content where post_id=?", id)
	if err != nil {
		log.Printf("post_content %d delete err %v content:", id, err)
		conn.Rollback()
		return
	}

	conn.Commit()
	return
}
func PostSave(post model.Post) (id int, err error) {
	conn, err := db.Begin()
	if err != nil {
		log.Fatalln("post db conn err ", err)
		return
	}
	var tagIds []byte
	tagIds, _ = json.Marshal(post.TagIds)
	var rs sql.Result
	if post.Id > 0 {
		id = post.Id
		_, err = conn.Exec("update post set title=?,description=?,category_id=?,tag_ids=? where id=?", post.Title, post.Description, post.CategoryId, tagIds, id)
		if err != nil {
			log.Printf("post %d update err %v", id, err)
			conn.Rollback()
			return
		}

		_, err = conn.Exec("update post_content set content=? where post_id=?", post.Content, id)
		if err != nil {
			log.Printf("post_content %d update err %v content:", id, err)
			conn.Rollback()
			return
		}

	} else {
		rs, err = conn.Exec("insert into post (title, description,category_id,tag_ids,status) values (?,?,?,?,?)", post.Title, post.Description, post.CategoryId, tagIds, post.Status)
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

func GetPostCountByTagId(id string) (int, error) {
	var count *int
	row := db.QueryRow("select count(*) from post where JSON_CONTAINS(tag_ids,'" + id + "')")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return *count, err
}

func GetPostIdsByContent(content string) ([]string, error) {
	rows, err := db.Query("select post_id from post_content where content like ?", "%"+content+"%")
	if err != nil {
		return nil, err
	}
	var postIds []string
	defer rows.Close()
	for rows.Next() {
		var postId string
		rows.Scan(&postId)
		postIds = append(postIds, postId)
	}
	return postIds, nil
}

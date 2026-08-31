package daos

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/convee/goblog/internal/model"
)

type PostParams struct {
	Ids        map[string][]string
	CategoryId string
	TagId      string
	PerPage    int
	Page       int
	Keyword    string
}

// PostQueryBuilder builds the bounded set of post-list queries used by the
// application. Field names remain internal constants; all user values are
// bound as parameters.
type PostQueryBuilder struct {
	conditions []string
	args       []interface{}
	orderBy    string
	limit      string
}

func NewPostQueryBuilder() *PostQueryBuilder {
	return &PostQueryBuilder{
		conditions: make([]string, 0, 4),
		args:       make([]interface{}, 0, 8),
		orderBy:    "order by is_top desc,updated_at desc",
	}
}

func (qb *PostQueryBuilder) WithCategory(categoryID string) *PostQueryBuilder {
	if categoryID != "" {
		qb.conditions = append(qb.conditions, "category_id=?")
		qb.args = append(qb.args, categoryID)
	}
	return qb
}

func (qb *PostQueryBuilder) WithTag(tagID string) *PostQueryBuilder {
	if tagID != "" {
		qb.conditions = append(qb.conditions, "JSON_CONTAINS(tag_ids,?)")
		qb.args = append(qb.args, tagID)
	}
	return qb
}

func (qb *PostQueryBuilder) WithKeyword(keyword string, ids map[string][]string) *PostQueryBuilder {
	if keyword == "" {
		return qb
	}

	keywordSQL := "title like ? escape '|' or description like ? escape '|'"
	qb.args = append(qb.args, likePattern(keyword), likePattern(keyword))
	if ids != nil {
		keywordSQL = qb.appendIn(keywordSQL, "id", ids["ids"])
		keywordSQL = qb.appendIn(keywordSQL, "category_id", ids["category_ids"])
		for _, tagID := range ids["tag_ids"] {
			if tagID == "" {
				continue
			}
			keywordSQL += " or JSON_CONTAINS(tag_ids,?)"
			qb.args = append(qb.args, tagID)
		}
	}
	qb.conditions = append(qb.conditions, "("+keywordSQL+")")
	return qb
}

func (qb *PostQueryBuilder) appendIn(query, field string, values []string) string {
	placeholders := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		placeholders = append(placeholders, "?")
		qb.args = append(qb.args, value)
	}
	if len(placeholders) == 0 {
		return query
	}
	return query + " or " + field + " in (" + strings.Join(placeholders, ",") + ")"
}

func likePattern(value string) string {
	replacer := strings.NewReplacer("|", "||", "%", "|%", "_", "|_")
	return "%" + replacer.Replace(value) + "%"
}

func (qb *PostQueryBuilder) WithStatus(status int) *PostQueryBuilder {
	qb.conditions = append(qb.conditions, "status=?")
	qb.args = append(qb.args, status)
	return qb
}

func (qb *PostQueryBuilder) WithPagination(page, perPage int) *PostQueryBuilder {
	if perPage > 0 && page > 0 {
		offset := (page - 1) * perPage
		qb.limit = " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(perPage)
	}
	return qb
}

func (qb *PostQueryBuilder) BuildQuery(fields string) string {
	querySQL := "select " + fields + " from post"
	if len(qb.conditions) > 0 {
		querySQL += " where " + strings.Join(qb.conditions, " and ")
	}
	if qb.orderBy != "" {
		querySQL += " " + qb.orderBy
	}
	if qb.limit != "" {
		querySQL += qb.limit
	}
	return querySQL
}

func (qb *PostQueryBuilder) GetArgs() []interface{} {
	return append([]interface{}(nil), qb.args...)
}

func GetPosts(params PostParams) ([]model.Post, error) {
	qb := NewPostQueryBuilder().
		WithCategory(params.CategoryId).
		WithTag(params.TagId).
		WithKeyword(params.Keyword, params.Ids).
		WithStatus(1).
		WithPagination(params.Page, params.PerPage)

	querySQL := qb.BuildQuery("id,title,created_at,updated_at,category_id,tag_ids,views,description")
	rows, err := db.Query(querySQL, qb.GetArgs()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		var tags []byte
		if err := rows.Scan(&post.Id, &post.Title, &post.CreatedAt, &post.UpdatedAt, &post.CategoryId, &tags, &post.Views, &post.Description); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(tags, &post.TagIds); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func GetPost(id string) (model.Post, error) {
	return getPost(id, false)
}

func GetPublishedPost(id string) (model.Post, error) {
	return getPost(id, true)
}

func getPost(id string, publishedOnly bool) (model.Post, error) {
	var post model.Post
	var tags []byte
	query := "select post.id,post.title,post.created_at,post.updated_at,post.category_id,post.tag_ids,post.views,post_content.content,post.description from post left join post_content on post.id=post_content.post_id where post.id=?"
	if publishedOnly {
		query += " and post.status=1"
	}
	row := db.QueryRow(query+" limit 1", id)
	if err := row.Scan(&post.Id, &post.Title, &post.CreatedAt, &post.UpdatedAt, &post.CategoryId, &tags, &post.Views, &post.Content, &post.Description); err != nil {
		return post, err
	}
	if err := json.Unmarshal(tags, &post.TagIds); err != nil {
		return post, err
	}
	return post, nil
}

func IncrView(id string) error {
	_, err := db.Exec("update post set views=views+1 where id=? and status=1", id)
	return err
}

func PostDelete(post model.Post) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec("delete from post where id=?", post.Id); err != nil {
		return 0, err
	}
	if _, err := tx.Exec("delete from post_content where post_id=?", post.Id); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return post.Id, nil
}

func PostSave(post model.Post) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	tagIDs, err := json.Marshal(post.TagIds)
	if err != nil {
		return 0, err
	}

	id := post.Id
	if id > 0 {
		if _, err := tx.Exec("update post set title=?,description=?,category_id=?,tag_ids=?,status=? where id=?", post.Title, post.Description, post.CategoryId, tagIDs, post.Status, id); err != nil {
			return 0, err
		}
		if _, err := tx.Exec("update post_content set content=? where post_id=?", post.Content, id); err != nil {
			return 0, err
		}
	} else {
		result, err := tx.Exec("insert into post (title,description,category_id,tag_ids,status) values (?,?,?,?,?)", post.Title, post.Description, post.CategoryId, tagIDs, post.Status)
		if err != nil {
			return 0, err
		}
		insertedID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		id = int(insertedID)
		if _, err := tx.Exec("insert into post_content (post_id,content) values (?,?)", id, post.Content); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func GetPostCountByTagID(id string) (int, error) {
	var count int
	if err := db.QueryRow("select count(*) from post where status=1 and JSON_CONTAINS(tag_ids,?)", id).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func GetPostIdsByContent(content string) ([]string, error) {
	rows, err := db.Query("select post_id from post_content where content like ? escape '|'", likePattern(content))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postIDs := make([]string, 0)
	for rows.Next() {
		var postID string
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}
	return postIDs, rows.Err()
}

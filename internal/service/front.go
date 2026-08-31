package service

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
)

type FrontService struct{}

func NewFrontService() *FrontService {
	return &FrontService{}
}

const (
	defaultPerPage = 20
	maxPerPage     = 50
	maxPage        = 10000
	maxKeywordSize = 100
)

// PostListQuery is the validated, transport-independent input for the public
// post list. HTTP handlers normalize request data before crossing this boundary.
type PostListQuery struct {
	CategoryID string
	TagID      string
	Keyword    string
	Page       int
	PerPage    int
}

func NewPostListQuery(categoryID, tagID, keyword string, page, perPage int) PostListQuery {
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		if size := utf8.RuneCountInString(keyword); size > maxKeywordSize {
			keyword = string([]rune(keyword)[:maxKeywordSize])
		}
	}

	if page < 1 {
		page = 1
	} else if page > maxPage {
		page = maxPage
	}
	if perPage < 1 {
		perPage = defaultPerPage
	} else if perPage > maxPerPage {
		perPage = maxPerPage
	}

	return PostListQuery{
		CategoryID: normalizeID(categoryID),
		TagID:      normalizeID(tagID),
		Keyword:    keyword,
		Page:       page,
		PerPage:    perPage,
	}
}

func normalizeID(value string) string {
	value = strings.TrimSpace(value)
	id, err := strconv.Atoi(value)
	if err != nil || id < 1 || id > math.MaxInt32 {
		return ""
	}
	return strconv.Itoa(id)
}

// PostList contains the current page and enough pagination metadata for a
// transport layer to render links without knowing the database fetch strategy.
type PostList struct {
	Posts       []model.Post
	Page        int
	PerPage     int
	HasPrevious bool
	HasNext     bool
}

func (s *FrontService) GetPostList(query PostListQuery) (PostList, error) {
	query = NewPostListQuery(query.CategoryID, query.TagID, query.Keyword, query.Page, query.PerPage)

	// Fetch one extra row so pagination does not require a second COUNT query
	// for every keyword/category combination.
	fetchQuery := query
	fetchQuery.PerPage++
	posts, err := s.getHomePosts(fetchQuery)
	if err != nil {
		return PostList{}, err
	}

	hasNext := len(posts) > query.PerPage
	if hasNext {
		posts = posts[:query.PerPage]
	}
	if posts == nil {
		posts = make([]model.Post, 0)
	}

	return PostList{
		Posts:       posts,
		Page:        query.Page,
		PerPage:     query.PerPage,
		HasPrevious: query.Page > 1,
		HasNext:     hasNext,
	}, nil
}

func (s *FrontService) getHomePosts(query PostListQuery) ([]model.Post, error) {
	params := daos.PostParams{
		CategoryId: query.CategoryID,
		TagId:      query.TagID,
		PerPage:    query.PerPage,
		Page:       query.Page,
	}
	if query.Keyword != "" {
		params.Keyword = query.Keyword
		params.Ids = make(map[string][]string)
		postIDs, err := daos.GetPostIdsByContent(query.Keyword)
		if err != nil {
			return nil, fmt.Errorf("search post content: %w", err)
		}
		params.Ids["ids"] = postIDs

		categoryIDs, err := daos.GetCategoryIdsByName(query.Keyword)
		if err != nil {
			return nil, fmt.Errorf("search categories: %w", err)
		}
		params.Ids["category_ids"] = categoryIDs

		tagIDs, err := daos.GetTagIdsByName(query.Keyword)
		if err != nil {
			return nil, fmt.Errorf("search tags: %w", err)
		}
		params.Ids["tag_ids"] = tagIDs
	}
	return NewPostService().GetPosts(params)
}

func (s *FrontService) GetPostDetail(id string) (model.Post, []model.Tag, error) {
	post, err := daos.GetPublishedPost(id)
	if err != nil {
		return post, nil, err
	}
	if err := daos.IncrView(id); err != nil {
		return post, nil, err
	}
	post.Views++
	category, err := daos.GetCategory(post.CategoryId)
	if err != nil {
		return post, nil, err
	}
	post.CategoryName = category.Name

	allTags, err := daos.GetTags()
	if err != nil {
		return post, nil, err
	}
	tagsByID := make(map[int]model.Tag, len(allTags))
	for _, tag := range allTags {
		tagsByID[tag.Id] = tag
	}
	tags := make([]model.Tag, 0, len(post.TagIds))
	for _, tagID := range post.TagIds {
		if tag, exists := tagsByID[tagID]; exists {
			tags = append(tags, tag)
		}
	}
	return post, tags, nil
}

func (s *FrontService) GetCategories() ([]model.Category, error) {
	return daos.GetCategories()
}

func (s *FrontService) GetTags() ([]model.Tag, error) {
	return daos.GetTags()
}

func (s *FrontService) GetPage(ident string) (model.Page, error) {
	return daos.GetPageByIdent(ident)
}

package service

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
)

type FrontService struct{}

func NewFrontService() *FrontService {
	return &FrontService{}
}

func (s *FrontService) GetHomePosts(categoryID, tagID, keyword string, page, perPage int) ([]model.Post, error) {
	params := daos.PostParams{
		CategoryId: categoryID,
		TagId:      tagID,
		PerPage:    perPage,
		Page:       page,
	}
	if keyword != "" {
		params.Keyword = keyword
		params.Ids = make(map[string][]string)
		postIDs, err := daos.GetPostIdsByContent(keyword)
		if err != nil {
			return nil, fmt.Errorf("search post content: %w", err)
		}
		params.Ids["ids"] = postIDs
		categoryIDs, err := daos.GetCategoryIdsByName(keyword)
		if err != nil {
			return nil, fmt.Errorf("search categories: %w", err)
		}
		params.Ids["category_ids"] = categoryIDs
		tagIDs, err := daos.GetTagIdsByName(keyword)
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

func (s *FrontService) GetPageURL(categoryID, tagID, keyword string, page int) string {
	query := url.Values{}
	if categoryID != "" {
		query.Set("category_id", categoryID)
	}
	if tagID != "" {
		query.Set("tag_id", tagID)
	}
	if keyword != "" {
		query.Set("keyword", keyword)
	}
	query.Set("page", strconv.Itoa(page))
	return "/?" + query.Encode()
}

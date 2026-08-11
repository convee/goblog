package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/convee/goblog/internal/daos"
	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/pkg/logger"
	"go.uber.org/zap"
)

const maxTagsPerPost = 20

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) GetPosts(params daos.PostParams) ([]model.Post, error) {
	posts, err := daos.GetPosts(params)
	if err != nil {
		return nil, err
	}
	categories, err := daos.GetCategories()
	if err != nil {
		return nil, err
	}

	categoryNames := make(map[int]string, len(categories))
	for _, category := range categories {
		categoryNames[category.Id] = category.Name
	}
	for i := range posts {
		posts[i].CategoryName = categoryNames[posts[i].CategoryId]
	}
	return posts, nil
}

func (s *PostService) GetPost(id string) (model.Post, error) {
	return daos.GetPost(id)
}

func (s *PostService) DeletePost(postID int) (int, error) {
	post, err := daos.GetPost(strconv.Itoa(postID))
	if err != nil {
		return 0, err
	}
	id, err := daos.PostDelete(model.Post{Id: postID})
	if err != nil {
		return 0, err
	}

	s.refreshTagCounts(post.TagIds)
	return id, nil
}

func (s *PostService) SavePost(post model.Post, tags string) (int, error) {
	oldTagIDs := make([]int, 0)
	if post.Id > 0 {
		oldPost, err := daos.GetPost(strconv.Itoa(post.Id))
		if err != nil {
			return 0, err
		}
		oldTagIDs = oldPost.TagIds
	}

	tagIDs, err := s.processTagIDs(tags)
	if err != nil {
		return 0, err
	}
	post.TagIds = tagIDs
	id, err := daos.PostSave(post)
	if err != nil {
		return 0, err
	}
	post.Id = id

	affectedTagIDs := append(oldTagIDs, post.TagIds...)
	s.refreshTagCounts(affectedTagIDs)
	return id, nil
}

func (s *PostService) GetPostTags(post model.Post) (string, error) {
	if len(post.TagIds) == 0 {
		return "", nil
	}
	allTags, err := daos.GetTags()
	if err != nil {
		return "", err
	}
	tagsByID := make(map[int]string, len(allTags))
	for _, tag := range allTags {
		tagsByID[tag.Id] = tag.Name
	}
	names := make([]string, 0, len(post.TagIds))
	for _, tagID := range post.TagIds {
		if name, ok := tagsByID[tagID]; ok {
			names = append(names, name)
		}
	}
	return strings.Join(names, ","), nil
}

func (s *PostService) processTagIDs(tags string) ([]int, error) {
	if strings.TrimSpace(tags) == "" {
		return []int{}, nil
	}
	allTags, err := daos.GetTags()
	if err != nil {
		return nil, err
	}
	tagsByName := make(map[string]model.Tag, len(allTags))
	for _, tag := range allTags {
		tagsByName[tag.Name] = tag
	}

	seen := make(map[string]struct{})
	tagIDs := make([]int, 0)
	for _, rawName := range strings.Split(tags, ",") {
		name := strings.TrimSpace(rawName)
		if name == "" {
			continue
		}
		if len([]rune(name)) > 100 {
			return nil, fmt.Errorf("标签 %q 超过 100 个字符", name)
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		if len(seen) > maxTagsPerPost {
			return nil, fmt.Errorf("每篇文章最多使用 %d 个标签", maxTagsPerPost)
		}

		if tag, exists := tagsByName[name]; exists {
			tagIDs = append(tagIDs, tag.Id)
			continue
		}
		newTagID, err := daos.AddTag(model.Tag{Name: name})
		if err != nil {
			// A concurrent request may have inserted the same unique tag.
			existing, getErr := daos.GetTagByName(name)
			if getErr != nil {
				return nil, err
			}
			newTagID = existing.Id
		}
		tagIDs = append(tagIDs, newTagID)
	}
	return tagIDs, nil
}

func (s *PostService) refreshTagCounts(tagIDs []int) {
	seen := make(map[int]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		if _, ok := seen[tagID]; ok {
			continue
		}
		seen[tagID] = struct{}{}
		if err := daos.RefreshTagCount(tagID); err != nil {
			logger.Error("refresh_tag_count", zap.Int("tag_id", tagID), zap.Error(err))
		}
	}
}

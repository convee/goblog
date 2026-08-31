package front

import (
	"testing"

	"github.com/convee/goblog/internal/model"
	"github.com/convee/goblog/internal/service"
)

func TestIndexViewDataPreparesFiltersAndPagination(t *testing.T) {
	query := service.NewPostListQuery("2", "7", "Go 安全", 3, 10)
	postList := service.PostList{
		Posts:       []model.Post{{Id: 11, Title: "Go 架构"}},
		Page:        3,
		PerPage:     10,
		HasPrevious: true,
		HasNext:     true,
	}
	categories := []model.Category{
		{Id: 1, Name: "默认"},
		{Id: 2, Name: "架构", Cur: 1},
	}

	data := indexViewData(query, postList, categories)

	if categories[1].Cur != 1 || categories[0].Cur != 0 {
		t.Fatalf("category selection = %+v, want only category 2 selected", categories)
	}
	if data["has_filters"] != true {
		t.Fatal("has_filters = false, want true for combined filters")
	}
	if got := data["pre_url"]; got != "/?category_id=2&keyword=Go+%E5%AE%89%E5%85%A8&page=2&per_page=10&tag_id=7" {
		t.Fatalf("pre_url = %v, want filters and per_page retained", got)
	}
	if got := data["next_url"]; got != "/?category_id=2&keyword=Go+%E5%AE%89%E5%85%A8&page=4&per_page=10&tag_id=7" {
		t.Fatalf("next_url = %v, want filters and per_page retained", got)
	}
	if got := data["clear_url"]; got != "/" {
		t.Fatalf("clear_url = %v, want a clean index URL", got)
	}
	if got := data["category_name"]; got != "架构" {
		t.Fatalf("category_name = %v, want selected category label", got)
	}
	if data["has_active_filters"] != true {
		t.Fatal("has_active_filters = false, want true for non-default filters")
	}
}

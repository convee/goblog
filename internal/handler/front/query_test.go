package front

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/convee/artgo"
	"github.com/convee/goblog/internal/service"
)

func TestPostListQueryFromRequestNormalizesValues(t *testing.T) {
	request := httptest.NewRequest("GET", "/?page=-2&per_page=999&keyword=++Go+安全++&category_id=+2+&tag_id=nope", nil)
	context := &artgo.Context{Writer: httptest.NewRecorder(), Req: request}

	got := postListQueryFromRequest(context)
	if got.Page != 1 || got.PerPage != 50 {
		t.Fatalf("pagination = page %d, per page %d; want 1, 50", got.Page, got.PerPage)
	}
	if got.Keyword != "Go 安全" {
		t.Fatalf("keyword = %q, want trimmed keyword", got.Keyword)
	}
	if got.CategoryID != "2" || got.TagID != "" {
		t.Fatalf("filters = category %q, tag %q; want 2 and an ignored invalid tag", got.CategoryID, got.TagID)
	}
}

func TestPageURLPreservesFiltersAndPerPage(t *testing.T) {
	query := service.NewPostListQuery("2", "7", "Go 安全", 3, 10)

	if got, want := pageURL(query, 4), "/?category_id=2&keyword=Go+%E5%AE%89%E5%85%A8&page=4&per_page=10&tag_id=7"; got != want {
		t.Fatalf("pageURL()\n got: %s\nwant: %s", got, want)
	}
}

func TestClearFilterURLRemovesAllViewState(t *testing.T) {
	got := clearFilterURL()
	if got != "/" {
		t.Fatalf("clearFilterURL() = %q, want a clean index URL", got)
	}
	if strings.Contains(got, "?") {
		t.Fatalf("clearFilterURL() still contains query state: %q", got)
	}
}

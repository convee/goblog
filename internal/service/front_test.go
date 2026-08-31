package service

import (
	"strings"
	"testing"
)

func TestNewPostListQueryNormalizesUserInput(t *testing.T) {
	got := NewPostListQuery(" 2 ", "-7", "  Go 安全  ", -3, 999)

	if got.CategoryID != "2" || got.TagID != "" {
		t.Fatalf("IDs = category %q, tag %q; want normalized category and ignored invalid tag", got.CategoryID, got.TagID)
	}
	if got.Keyword != "Go 安全" {
		t.Fatalf("keyword = %q, want trimmed keyword", got.Keyword)
	}
	if got.Page != 1 || got.PerPage != 50 {
		t.Fatalf("pagination = page %d, per page %d; want 1, 50", got.Page, got.PerPage)
	}
}

func TestNewPostListQueryBoundsLargeValues(t *testing.T) {
	got := NewPostListQuery("", "", strings.Repeat("安", 101), 999999, -1)

	if got.Page != 10000 {
		t.Fatalf("page = %d, want 10000", got.Page)
	}
	if got.PerPage != 20 {
		t.Fatalf("per page = %d, want 20", got.PerPage)
	}
	if size := len([]rune(got.Keyword)); size != 100 {
		t.Fatalf("keyword length = %d, want 100", size)
	}
}

package service

import "testing"

func TestGetPageURLPreservesFilters(t *testing.T) {
	service := NewFrontService()
	got := service.GetPageURL("2", "7", "Go 安全", 3)
	want := "/?category_id=2&keyword=Go+%E5%AE%89%E5%85%A8&page=3&tag_id=7"
	if got != want {
		t.Fatalf("GetPageURL() = %q, want %q", got, want)
	}
}

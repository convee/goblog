package pkg

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFrontTemplatesCompile(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file")
	}
	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	if err := loadFrontTemplates(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"index", "post", "page", "tag", "404"} {
		if frontTemplates[name] == nil {
			t.Fatalf("template %q was not loaded", name)
		}
	}
}

func TestIndexTemplateRendersUsableListInteractions(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test file")
	}
	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	if err := loadFrontTemplates(); err != nil {
		t.Fatal(err)
	}
	data := map[string]interface{}{
		"name":               "测试博客",
		"cdn":                "/static",
		"asset_version":      "test",
		"title":              "文章",
		"description":        "文章列表",
		"keyword":            "不存在",
		"category_id":        "2",
		"tag_id":             "7",
		"per_page":           20,
		"page":               3,
		"has_filters":        true,
		"has_active_filters": true,
		"has_pre":            true,
		"has_next":           false,
		"pre_url":            "/?category_id=2&keyword=%E4%B8%8D%E5%AD%98%E5%9C%A8&page=2&per_page=20&tag_id=7",
		"next_url":           "/?category_id=2&keyword=%E4%B8%8D%E5%AD%98%E5%9C%A8&page=4&per_page=20&tag_id=7",
		"clear_url":          "/",
		"category_name":      "架构",
		"categories": []map[string]interface{}{
			{"Id": 1, "Name": "默认", "Cur": false},
			{"Id": 2, "Name": "架构", "Cur": true},
		},
		"posts": []map[string]interface{}{},
	}

	var output bytes.Buffer
	if err := frontTemplates["index"].ExecuteTemplate(&output, "layout.html", data); err != nil {
		t.Fatal(err)
	}
	page := output.String()
	for _, fragment := range []string{
		`<form class="filter-form" action="/" method="GET">`,
		`name="per_page"`,
		`name="tag_id" value="7"`,
		"没有找到符合条件的文章",
		"重置全部",
		"本页显示 0 篇",
		"分类：架构",
		`href="/?category_id=2&amp;keyword=%E4%B8%8D%E5%AD%98%E5%9C%A8&amp;page=2&amp;per_page=20&amp;tag_id=7"`,
		`<li class="disabled"><span aria-disabled="true">下一页</span></li>`,
		`src="/static/js/site.js?v=test"`,
	} {
		if !strings.Contains(page, fragment) {
			t.Fatalf("index output is missing %q\noutput:\n%s", fragment, page)
		}
	}
	if strings.Contains(page, `role="search"`) {
		t.Fatal("index output contains the duplicate navigation search form")
	}
}

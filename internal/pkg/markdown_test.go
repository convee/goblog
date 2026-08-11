package pkg

import (
	"strings"
	"testing"
)

func TestRenderMarkdownRejectsExecutableHTML(t *testing.T) {
	got, err := RenderMarkdown("# Safe\n\n<script>alert(1)</script>\n\n[x](javascript:alert(1))")
	if err != nil {
		t.Fatal(err)
	}
	html := string(got)
	if strings.Contains(html, "<script") || strings.Contains(html, "javascript:") {
		t.Fatalf("unsafe content was rendered: %s", html)
	}
	if !strings.Contains(html, "<h1>Safe</h1>") {
		t.Fatalf("expected Markdown heading, got: %s", html)
	}
}

func TestRenderMarkdownSupportsGFM(t *testing.T) {
	got, err := RenderMarkdown("| a | b |\n| - | - |\n| 1 | 2 |")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "<table>") {
		t.Fatalf("expected GFM table, got: %s", got)
	}
}

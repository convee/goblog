package pkg

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Footnote,
	),
)

// RenderMarkdown converts trusted or untrusted Markdown with Goldmark's safe
// defaults. Raw HTML and dangerous link schemes are not emitted.
func RenderMarkdown(source string) (template.HTML, error) {
	var output bytes.Buffer
	if err := markdownRenderer.Convert([]byte(source), &output); err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}
	return template.HTML(output.String()), nil // #nosec G203 -- Goldmark unsafe HTML rendering is not enabled.
}

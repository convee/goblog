package pkg

import (
	"os"
	"path/filepath"
	"runtime"
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

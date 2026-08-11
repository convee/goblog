package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAppliesEnvironmentOverridesAndDefaults(t *testing.T) {
	previous := Conf
	t.Cleanup(func() { Conf = previous })

	path := filepath.Join(t.TempDir(), "config.yml")
	content := []byte("app:\n  mode: debug\nmysql:\n  Dsn: file-dsn\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOBLOG_APP_NAME", "Environment Blog")
	t.Setenv("GOBLOG_MYSQL_DSN", "environment-dsn")
	t.Setenv("GOBLOG_APP_TRUSTED_PROXIES", "10.0.0.0/8, 192.0.2.10/32")

	cfg := Init(path)
	if cfg.App.Name != "Environment Blog" {
		t.Fatalf("name = %q", cfg.App.Name)
	}
	if cfg.Mysql.DSN != "environment-dsn" {
		t.Fatalf("DSN = %q", cfg.Mysql.DSN)
	}
	if cfg.App.Cdn != "/static" || len(cfg.App.Navigation) == 0 {
		t.Fatalf("defaults were not applied: %#v", cfg.App)
	}
	if len(cfg.App.TrustedProxies) != 2 {
		t.Fatalf("trusted proxies = %#v", cfg.App.TrustedProxies)
	}
}

func TestDefaultPathUsesEnvironment(t *testing.T) {
	t.Setenv("GOBLOG_CONFIG", "/tmp/goblog-config.yml")
	if got := DefaultPath(); got != "/tmp/goblog-config.yml" {
		t.Fatalf("DefaultPath() = %q", got)
	}
}

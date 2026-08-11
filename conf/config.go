package conf

import (
	"github.com/convee/goblog/pkg/logger"
	"github.com/convee/goblog/pkg/storage/mysql"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config global config
type Config struct {
	App    AppConfig
	Logger logger.Config
	Mysql  mysql.Config
}

// AppConfig app config
type AppConfig struct {
	Name            string
	Description     string
	Keywords        string
	Mode            string
	Addr            string
	Cdn             string
	DisableDingDing bool
	TrustedProxies  []string
	Navigation      []NavigationItem
}

type NavigationItem struct {
	Label string
	URL   string
}

var (
	// Conf app global config
	Conf = &Config{}
)

func Init(configPath string) *Config {
	loader := viper.New()
	loader.SetConfigType("yml")
	loader.SetConfigFile(configPath)
	if err := loader.ReadInConfig(); err != nil {
		panic(err)
	}
	loader.SetEnvPrefix("GOBLOG")
	loader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, key := range []string{
		"app.name",
		"app.description",
		"app.keywords",
		"app.mode",
		"app.addr",
		"app.cdn",
		"app.disabledingding",
		"mysql.dsn",
		"mysql.maxopenconn",
		"mysql.maxidleconn",
		"mysql.connmaxlifetime",
	} {
		if err := loader.BindEnv(key); err != nil {
			panic(err)
		}
	}
	loader.AutomaticEnv()
	cfg := &Config{}
	if err := loader.Unmarshal(cfg); err != nil {
		panic(err)
	}
	if raw := os.Getenv("GOBLOG_APP_TRUSTED_PROXIES"); raw != "" {
		cfg.App.TrustedProxies = splitList(raw)
	}
	applyDefaults(cfg)
	Conf = cfg
	return Conf
}

func applyDefaults(cfg *Config) {
	if cfg.App.Name == "" {
		cfg.App.Name = "goblog"
	}
	if cfg.App.Description == "" {
		cfg.App.Description = "Go Markdown 博客系统"
	}
	if cfg.App.Keywords == "" {
		cfg.App.Keywords = "go,golang,markdown,blog"
	}
	if cfg.App.Addr == "" {
		cfg.App.Addr = "localhost:9091"
	}
	if cfg.App.Cdn == "" {
		cfg.App.Cdn = "/static"
	}
	if len(cfg.App.Navigation) == 0 {
		cfg.App.Navigation = []NavigationItem{
			{Label: "首页", URL: "/"},
			{Label: "标签", URL: "/tag"},
			{Label: "关于", URL: "/page/about"},
		}
	}
}

func DefaultPath() string {
	if path := strings.TrimSpace(os.Getenv("GOBLOG_CONFIG")); path != "" {
		return path
	}
	return "./conf/dev.yml"
}

func splitList(raw string) []string {
	values := make([]string, 0)
	for _, value := range strings.Split(raw, ",") {
		if value = strings.TrimSpace(value); value != "" {
			values = append(values, value)
		}
	}
	return values
}

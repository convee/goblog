# goblog

一个使用 Go、MySQL 和服务端模板实现的轻量 Markdown 博客，包含文章、页面、分类、标签、搜索和管理后台。

> 安全提示：公开发行的 v1.2.1 及更早版本存在后台会话认证缺陷，不应继续暴露在公网。请升级到修复版本，并撤销旧 Gitalk OAuth secret、删除或重置默认管理员。

## 能力与边界

- 前台：文章列表与详情、分类/标签筛选、内容搜索、自定义 Markdown 页面
- 后台：文章、页面和分类管理，Markdown 编辑器
- 安全：bcrypt 密码、签名且可过期的会话、CSRF、防滥用登录限流、安全响应头和安全 Markdown
- 运行：MySQL 连接池、健康/就绪探针、HTTP 超时、优雅退出
- 质量：单元测试、race、vet、govulncheck、GitHub Actions 和发行归档

goblog 的默认运行时只依赖 MySQL。当前搜索适合个人博客和中小数据量，使用 MySQL 模糊匹配，不承诺大型全文检索场景。

## 环境要求

- Go 1.25 或更高版本；仓库固定使用 Go 1.26.7 工具链
- MySQL 8.4 LTS
- Docker 与 Docker Compose：可选，用于最快启动

## Docker 快速开始

1. 设置本地开发密码和会话密钥：

   ```bash
   export GOBLOG_MYSQL_PASSWORD='replace-this-password'
   export GOBLOG_MYSQL_ROOT_PASSWORD='replace-this-root-password'
   export GOBLOG_SESSION_SECRET="$(openssl rand -hex 32)"
   ```

2. 构建并启动：

   ```bash
   docker compose up --build -d
   ```

3. 创建管理员：

   ```bash
   docker compose run --rm \
     --entrypoint /app/bin/goblog-admin \
     goblog --email you@example.com
   ```

   命令会在终端中提示输入管理员密码。

4. 访问 <http://localhost:9091>，后台登录地址为 <http://localhost:9091/login>。

删除 Compose 数据卷会永久删除本地博客数据。升级或清理前请先备份 MySQL。

## 从源码启动

1. 创建数据库并导入初始化结构：

   ```bash
   mysql -u root -p -e 'CREATE DATABASE blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'
   mysql -u root -p blog < blog.sql
   ```

2. 创建本地配置并编辑 MySQL DSN、站点名称和导航：

   ```bash
   cp conf/config.example.yml conf/local.yml
   ```

3. 创建管理员并启动：

   ```bash
   export GOBLOG_CONFIG=conf/local.yml
   export GOBLOG_SESSION_SECRET="$(openssl rand -hex 32)"
   go run ./cmd/admin --email you@example.com
   go run .
   ```

   管理员命令会在终端安全读取至少 12 个字符的密码。非交互环境可以用 `GOBLOG_ADMIN_PASSWORD` 提供，但不要把密码写入仓库或 Shell 历史。

进程存活探针为 `/healthz`，包含 MySQL 检查的就绪探针为 `/readyz`。

## 配置

配置示例位于 `conf/config.example.yml`，开发示例位于 `conf/dev.yml`。生产配置应保存为被 Git 忽略的 `conf/config.yml`，并由 Secret 管理系统注入敏感值。

常用环境变量：

| 变量 | 用途 |
|---|---|
| `GOBLOG_CONFIG` | 配置文件路径 |
| `GOBLOG_SESSION_SECRET` | 生产必需的稳定会话密钥，至少 32 字节 |
| `GOBLOG_MYSQL_DSN` | 覆盖 MySQL DSN |
| `GOBLOG_APP_NAME` | 站点名称 |
| `GOBLOG_APP_CDN` | 静态资源前缀，默认 `/static` |
| `GOBLOG_APP_TRUSTED_PROXIES` | 逗号分隔的可信代理 CIDR；只有这些代理的 `X-Forwarded-Proto` 会被接受 |
| `GOBLOG_DINGTALK_TOKEN` | 可选的钉钉告警 token |

非 `debug` 模式未设置 `GOBLOG_SESSION_SECRET` 时进程会拒绝启动。轮换密钥会注销全部现有会话。

## 开发与验证

```bash
make             # 格式、vet、race 测试并构建
make check       # 格式、vet、race 测试
make build       # 为当前平台构建；可用 GOOS/GOARCH 交叉编译
make vuln        # 扫描可达 Go 漏洞
make package GOOS=linux GOARCH=amd64  # 生成完整发行归档和 SHA256
./bin/goblog -v  # 查看版本、提交和构建时间
```

## 数据库升级

全新安装直接导入 `blog.sql`。从 v1.2.x 升级时，先备份数据库并检查重复的分类名、页面标识、标签名和管理员邮箱，再执行：

```bash
mysql -u root -p blog < migrations/001_integrity_and_indexes.sql
```

该迁移只能执行一次。旧部署如果启用了 Redis 或 Elasticsearch，可以在确认没有其他应用使用后移除；新版本不再连接这些组件。

## 生产运行

发行归档包含二进制、管理员命令、模板、静态资源、数据库脚本、配置和许可证。解压并修改 `conf/config.yml` 后：

```bash
export GOBLOG_SESSION_SECRET='at-least-32-bytes-from-a-secret-manager'
./startup.sh start
./startup.sh status
./startup.sh stop
```

生产环境应使用 systemd、容器编排或其他进程管理器，在受信任的反向代理后终止 TLS，并保留网关级登录限流。应用自身的登录限流是第二道防线。

## 目录结构

```text
cmd/admin/          管理员初始化命令
conf/               配置定义与示例
internal/daos/      MySQL 数据访问
internal/service/   业务编排
internal/handler/   前台和后台 HTTP 处理
internal/routers/   路由与中间件
internal/security/  会话、CSRF、代理信任与登录限流
templates/          服务端 HTML 模板
static/             CSS、JavaScript 与 Markdown 编辑器资源
migrations/         存量数据库迁移
```

## 贡献与安全

- 贡献说明：[CONTRIBUTING.md](CONTRIBUTING.md)
- 行为准则：[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- 支持范围：[SUPPORT.md](SUPPORT.md)
- 安全报告：[SECURITY.md](SECURITY.md)
- 更新记录：[CHANGELOG.md](CHANGELOG.md)

## License

[MIT](LICENSE)

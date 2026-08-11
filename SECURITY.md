# Security Policy

## Supported versions

| Version | Supported |
|---|---|
| 默认分支最新版本 | Yes |
| v1.2.1 及更早版本 | No |

v1.2.1 及更早版本存在后台会话认证缺陷，不应继续暴露在公网。升级前请先备份数据库，并在升级后注销旧会话、撤销旧 Gitalk OAuth secret、删除或重置默认管理员。

## Reporting a vulnerability

请通过 GitHub Security Advisory 的私密报告入口联系维护者，不要先创建公开 issue。报告应包含受影响版本、复现条件、影响和建议修复方向；请勿附带真实生产凭据或个人数据。

维护者确认问题后会在私密 Advisory 中协调修复和披露。修复版本发布前，请避免公开复现细节。

## Deployment requirements

- 首次部署必须用 `goblog-admin` 或 `go run ./cmd/admin` 创建独立强密码管理员。
- 非 `debug` 模式必须从 Secret 管理系统注入至少 32 字节的 `GOBLOG_SESSION_SECRET`。
- 生产环境必须使用 HTTPS。只有在 `app.TrustedProxies` 或 `GOBLOG_APP_TRUSTED_PROXIES` 中配置的代理才可传递 HTTPS 状态。
- 应用内置基础登录限流；生产网关仍应限制 `/login` 和 `/signin` 的访问频率。
- `conf/config.yml`、数据库凭据、告警 token、SMTP 密码和 OAuth secret 不得提交到仓库。
- Markdown 原始 HTML 默认不执行；不要为了兼容旧内容重新启用任意 HTML 解码。
- 发布或升级前至少执行：

  ```bash
  make check
  make vuln
  ```

## Historical credential exposure

旧版本曾把 Gitalk OAuth client secret 下发到浏览器。仓库所有者必须先在 OAuth App 中撤销或轮换它，再评估是否需要协调 forks 和 clones 重写 Git 历史。仅从最新提交删除字符串不能使已泄露凭据重新安全。

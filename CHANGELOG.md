# Changelog

本项目遵循 [Semantic Versioning](https://semver.org/)。

## Unreleased

### Security

- 使用带过期时间的 HMAC 会话替代可伪造的明文身份 Cookie。
- 为后台写操作增加 CSRF 防护、请求体上限和安全响应头。
- 移除浏览器端 Gitalk OAuth secret 和默认弱口令管理员。
- 默认安全渲染 Markdown，不执行原始 HTML。

### Changed

- 运行时仅依赖 MySQL，Redis 和 Elasticsearch 不再是部署组件。
- 增加健康检查、就绪检查、HTTP 超时和更可靠的优雅退出。
- 增加管理员初始化命令、数据库完整性迁移、CI 和发行归档。
- 最低数据库基线调整为 MySQL 8.4 LTS。

### Fixed

- 修复后台认证中间件无法可靠短路的问题。
- 修复文章、页面、分类和标签相关的数据一致性与错误处理问题。

## 1.2.1 - 2023-10-25

- 调整路由与专题页标识。

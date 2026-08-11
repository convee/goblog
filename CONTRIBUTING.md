# Contributing

感谢你改进 goblog。提交改动前请先搜索现有 issue，较大的行为变化建议先创建讨论 issue。

## 本地开发

1. 安装 Go 1.25 或更高版本和 MySQL 8.4。
2. 导入 `blog.sql`，复制并调整 `conf/config.example.yml`。
3. 运行 `make check`。
4. 涉及数据库、模板或浏览器行为时，同时执行对应的集成或端到端测试。

## Pull Request

- 一个 PR 聚焦一个问题，说明动机、行为变化和验证结果。
- 安全修复不要先提交公开 issue，请遵循 `SECURITY.md`。
- 不提交凭据、生产配置、日志、构建产物或无关格式化。
- 新功能应补测试和文档；破坏兼容性的变更应说明迁移方法。
- 提交即表示你同意按本项目 MIT License 授权你的贡献。

/ *
Navicat Premium Data Transfer

Source Server : rds
Source Server Type : MySQL
Source Server Version : 80018(8.0.18)
Source Host           : rm-uf6amy1qxilm4k083.mysql.rds.aliyuncs.com:3306
Source Schema : blog

Target Server Type : MySQL
Target Server Version : 80018(8.0.18)
File Encoding         : 65001 Date : 11 / 02 / 2023 22:06:59
* /

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`
(
`id` int(11
)
unsigned NOT NULL AUTO_INCREMENT,
`name` varchar
(
100
)
CHARACTER
SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY(`id`)
)ENGINE = InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of category
-- ----------------------------
BEGIN
;
INSERT INTO `category`(`id`, `name`, `created_at`, `updated_at`)
VALUES (5, '技术', '2023-01-09 23:54:35', '2023-01-09 23:54:35');
INSERT INTO `category`(`id`, `name`, `created_at`, `updated_at`)
VALUES (6, '产品', '2023-01-09 23:55:57', '2023-01-09 23:55:57');
COMMIT
;

-- ----------------------------
-- Table structure for page
-- ----------------------------
DROP TABLE IF EXISTS `page`;
CREATE TABLE `page`
(
`id` int(11
)
unsigned NOT NULL AUTO_INCREMENT,
`ident` varchar
(
20
)
CHARACTER
SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`title` varchar(100)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
PRIMARY KEY(`id`)
)ENGINE = InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of page
-- ----------------------------
BEGIN
;
COMMIT
;

-- ----------------------------
-- Table structure for post
-- ----------------------------
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post`
(
`id` int(11
)
unsigned NOT NULL AUTO_INCREMENT,
`title` varchar
(
200
)
CHARACTER
SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`status` tinyint(1)NOT NULL DEFAULT '0',
`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
`category_id` int(11)NOT NULL DEFAULT '0',
`is_top` tinyint(1)NOT NULL DEFAULT '0',
`tag_ids` json NOT NULL,
`views` int(11)NOT NULL DEFAULT '0',
`description` varchar(2000)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
PRIMARY KEY(`id`)
)ENGINE = InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of post
-- ----------------------------
BEGIN
;
INSERT INTO `post`(`id`, `title`, `status`, `created_at`, `updated_at`, `category_id`, `is_top`, `tag_ids`, `views`,
`description`)
VALUES (1, 'Go Markdown 博客系统', 1, '2023-01-09 23:55:05', '2023-01-09 23:55:05', 5, 0, '[28, 29]', 301,
'基于 Go 语言实现的 Markdown 博客系统');
COMMIT
;

-- ----------------------------
-- Table structure for post_content
-- ----------------------------
DROP TABLE IF EXISTS `post_content`;
CREATE TABLE `post_content`
(
`id` int(11
)
unsigned NOT NULL AUTO_INCREMENT,
`content` longtext CHARACTER
SET utf8mb4 COLLATE utf8mb4_general_ci,
`post_id` int(11)NOT NULL DEFAULT '0',
PRIMARY KEY(`id`)
)ENGINE = InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of post_content
-- ----------------------------
BEGIN
;
INSERT INTO `post_content`(`id`, `content`, `post_id`)
VALUES (1,
'## Go Markdown 博客系统\r\n> 基于 Go 语言实现的 Markdown 博客系统\r\n\r\n### 技术栈\r\n\r\n* 前端框架：[Bootstrap v3.3.7](http://getbootstrap.com)\r\n* 语言：[go](https://go.dev/)\r\n* 网络库：标准库 net/http\r\n* 配置文件解析库 [Viper](https://github.com/spf13/viper)\r\n* 日志库：[zap](https://github.com/uber-go/zap)\r\n* 搜索引擎：[elasticsearch](https://github.com/olivere/elastic/v7)\r\n* 数据库：[mysql](https://github.com/go-sql-driver/mysql)\r\n* 缓存：[redis](https://github.com/go-redis/redis)\r\n* 文件存储：阿里云 oss、cdn\r\n* markdown 编辑器：[markdown editor](https://github.com/pandao/editor.md)\r\n* pprof 性能调优\r\n* 包管理工具 [Go Modules](https://github.com/golang/go/wiki/Modules)\r\n* 评论插件：[gitalk](https://github.com/gitalk/gitalk) \r\n* 后台登录：cookie \r\n* 使用 make 来管理 Go 工程\r\n* 使用 shell(startup.sh) 脚本来管理进程\r\n* 使用 YAML 文件进行多环境配置\r\n* 优雅退出\r\n* Http 请求 panic 异常捕获\r\n* 错误信息钉钉预警\r\n\r\n### 目录结构\r\n\r\n```shell\r\n├── Makefile                     # 项目管理文件\r\n├── conf                         # 配置文件统一存放目录\r\n├── internal                     # 业务目录\r\n│   ├── handler                  # http 接口\r\n│   ├── pkg                      # 内部应用程序代码\r\n│   └── routers                  # 业务路由\r\n├── logs                         # 存放日志的目录\r\n├── static                       # 存放静态文件的目录\r\n├── tpl                          # 存放模板的目录\r\n├── main.go                      # 项目入口文件\r\n├── pkg                          # 公共的 package\r\n├── tests                        # 单元测试\r\n└── startup.sh                   # 启动脚本\r\n```\r\n\r\n### 功能模块\r\n\r\n#### 后台\r\n* 文章管理：文章增删改查\r\n* 页面管理：页面增删改查，可自定义 markdown 页面\r\n* 分类管理：分类增删改查\r\n* 标签管理：标签列表\r\n  \r\n#### 前台\r\n* 文章列表：倒序展示文章、可置顶\r\n* 内容页面：markdown 内容展示\r\n* 标签页面：按标签文章数量排序\r\n* 关于页面：个人说明\r\n* 阅读清单：个人阅读书籍\r\n* 站内搜索：支持文章标题、描述、内容、分类、标签模糊搜索\r\n\r\n## 开发规范\r\n\r\n遵循: [Uber Go 语言编码规范](https://github.com/uber-go/guide/blob/master/style.md)\r\n\r\n### 常用命令\r\n\r\n- make help 查看帮助\r\n- make dep 下载 Go 依赖包\r\n- make build 编译项目\r\n- make tar 打包文件\r\n\r\n### 部署流程\r\n* 依赖环境：\r\n  \r\n   mysql、redis、elasticsearch\r\n   > elasticsearch 可通过配置开启关闭，redis主要考虑到后续加缓存\r\n  \r\n* 安装部署\r\n\r\n```\r\n# 下载安装，可以不用是 GOPATH\r\ngit clone https://github.com/convee/goblog.git\r\n\r\n# 进入到下载目录\r\ncd goblog\r\n\r\n# 生成环境配置文件\r\ncd conf\r\n\r\n# 修改 mysql、redis、elasticsearch 配置\r\n\r\n# 导入初始化 sql 结构\r\nmysql -u root -p\r\n> create database blog;\r\n> set names utf8mb4;\r\n> use blog;\r\n> source blog.sql;\r\n\r\n\r\n# 下载依赖\r\nmake dep\r\n\r\n# 编译\r\nmake build\r\n\r\n# 运行\r\n./goblog dev.yml\r\n\r\n# 后台运行\r\nnohup ./goblog dev.yml &\r\n```\r\n\r\n* supervisord 部署\r\n  \r\n```\r\n[program:goblog]\r\ndirectory = /data/modules/blog\r\ncommand = /data/modules/blog/goblog -c conf/prod.yml\r\nautostart = true\r\nautorestart = true\r\nstartsecs = 5\r\nuser = root\r\nredirect_stderr = true\r\nstdout_logfile = /data/modules/blog/supervisor.log\r\n```\r\n\r\n* 访问首页\r\n\r\nhttp://localhost:9091\r\n\r\n* 访问后台\r\n\r\nhttp://localhost:9091/admin\r\n  \r\n用户名：admin@convee.cn\r\n  \r\n密码：123456\r\n\r\n* 演示站：https://convee.cn',
1);
COMMIT
;

-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag`
(
`id` int(11
)
unsigned NOT NULL AUTO_INCREMENT,
`name` varchar
(
100
)
CHARACTER
SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`count` int(11)NOT NULL DEFAULT '0',
`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY(`id`)
)ENGINE = InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tag
-- ----------------------------
BEGIN
;
INSERT INTO `tag`(`id`, `name`, `count`, `created_at`, `updated_at`)
VALUES (28, 'go', 1, '2023-01-09 23:55:05', '2023-01-09 23:55:05');
INSERT INTO `tag`(`id`, `name`, `count`, `created_at`, `updated_at`)
VALUES (29, 'markdown', 1, '2023-01-09 23:55:05', '2023-01-09 23:55:05');
COMMIT
;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
`id` int(11
)
unsigned NOT NULL AUTO_INCREMENT,
`email` varchar
(
50
)
CHARACTER
SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`password` varchar(128)CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
PRIMARY KEY(`id`)
)ENGINE = InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of user
-- ----------------------------
BEGIN;
INSERT INTO `user` (`id`, `email`, `password`) VALUES (2, 'convee.cn@gmail.com', '$2a$08$dhShLwIUiikE8FLwUrZHXO/Qx2elkF00DyEO/F3sbHlnqqZLIfLYe');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

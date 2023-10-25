<p align="center">
	<strong>goblog 基于 Go 语言实现的 Markdown 博客系统</strong>
</p>
<p align="center">
   <a target="_blank" href="#">
      <img style="display: inline-block;" src="https://img.shields.io/badge/Go-1.17.13-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/Bootstrap-3.3.7-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/Mysql-5.7-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/GORM-v1.24.3-blue"/>
      <img style="display: inline-block;" src="https://img.shields.io/badge/Redis-6.0-red"/>
    </a>
</p>

[在线预览](#在线预览) | [项目介绍](#项目介绍) | [技术介绍](#技术介绍) | [目录结构](#目录结构) | [环境说明](#环境说明) | [快速开始](#快速开始) | [总结&鸣谢](#总结鸣谢)  | [后续计划](#后续计划) | [更新日志](#更新日志)


## 在线预览

* 博客前台链接：[convee.cn](https://www.convee.cn)

* 博客后台链接：[convee.cn/admin](https://www.convee.cn/admin)

### 项目介绍

#### 后台

* cookie 鉴权
* markdown 编辑器
* 文章管理：文章增删改查
* 页面管理：自定义 markdown 专题页面
* 分类管理：分类增删改查
* 标签管理：标签列表

#### 前台

* 文章列表：倒序展示文章、可置顶
* 内容页面：markdown 内容展示
* 标签页面：按标签文章数量排序
* 关于页面：个人说明
* 阅读清单：个人阅读书籍
* 站内搜索：支持文章标题、描述、内容、分类、标签模糊搜索


## 技术介绍

* 前端框架：[Bootstrap v3.3.7](http://getbootstrap.com)
* 语言：[go](https://go.dev/)
* 网络库：标准库 net/http
* 配置文件解析库 [Viper](https://github.com/spf13/viper)
* 日志库：[zap](https://github.com/uber-go/zap)
* 搜索引擎：[elasticsearch](https://github.com/olivere/elastic/v7)
* 数据库：[mysql](https://github.com/go-sql-driver/mysql)
* 缓存：[redis](https://github.com/go-redis/redis)
* 文件存储：阿里云 oss、cdn
* markdown 编辑器：[markdown editor](https://github.com/pandao/editor.md)
* pprof 性能调优
* 包管理工具 [Go Modules](https://github.com/golang/go/wiki/Modules)
* 评论插件：[gitalk](https://github.com/gitalk/gitalk)
* 后台登录鉴权：cookie
* 使用 make 来管理 Go 工程
* 使用 shell(startup.sh) 脚本来管理进程
* 使用 YAML 文件进行多环境配置
* 优雅退出
* Http 请求 panic 异常捕获
* 错误信息钉钉预警

## 目录结构

```shell
├── Makefile                     # 项目管理文件
├── conf                         # 配置文件统一存放目录
├── internal                     # 业务目录
│   ├── handler                  # http 接口
│   ├── pkg                      # 内部应用程序代码
│   └── routers                  # 业务路由
├── logs                         # 存放日志的目录
├── static                       # 存放静态文件的目录
├── tpl                          # 存放模板的目录
├── main.go                      # 项目入口文件
├── pkg                          # 公共的 package
├── tests                        # 单元测试
└── startup.sh                   # 启动脚本
```

## 环境说明

| 开发工具                          | 说明                  |
| ----------------------------- | ------------------- |
| Vscode   | Golang 后端 + Javascript 前端 |
| Navicat    | MySQL 远程连接工具        |
| RDM | Redis 远程连接工具        |

| 开发环境   | 版本   |
| ------ | ---- |
| Golang | 1.17.13 |
| MySQL  | 5.7  |
| Redis  | 6.x  |

## 开发规范

遵循: [Uber Go 语言编码规范](https://github.com/uber-go/guide/blob/master/style.md)

## 快速开始

### make命令

- make help 查看帮助
- make dep 下载 Go 依赖包
- make build 编译项目
- make tar 打包文件

### 部署流程

* 安装部署

```
# 下载安装，可以不用是 GOPATH
git clone https://github.com/convee/goblog.git

# 进入到下载目录
cd goblog

# 生成环境配置文件
cd conf

# 修改 mysql、redis、elasticsearch 配置

# 导入初始化 sql 结构
mysql -u root -p
> create database blog;
> set names utf8mb4;
> use blog;
> source blog.sql;

# 下载依赖
make dep

# 编译
make build

# 后台运行
nohup ./goblog dev.yml &

```

* supervisord 部署

```
[program:goblog]
directory = /data/modules/blog
command = /data/modules/blog/goblog -c conf/prod.yml
autostart = true
autorestart = true
startsecs = 5
user = root
redirect_stderr = true
stdout_logfile = /data/modules/blog/supervisor.log
```

* 访问首页

http://localhost:9091

* 访问后台

http://localhost:9091/admin

用户名：admin@convee.cn

密码：123456

## 总结鸣谢
该项目后端使用原生 Golang，前端使用 jQuery，代码风格简洁，注释完善，适合 Golang 初学者学习。这个过程中也参考了很多优秀的开源项目，感谢大家的开源让程序员的世界更加丰富。

## 后续计划

* 图片上传
* 权限管理

## 更新日志
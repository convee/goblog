## Go Markdown 博客系统
> 基于 Go 语言实现的 Markdown 博客系统

### 技术栈

* 前端框架：bootstrap
* 语言：go
* 网络库：标准库 net/http
* 配置文件解析库 [Viper](https://github.com/spf13/viper)
* 日志库zap：https://github.com/uber-go/zap
* elasticsearch：https://github.com/olivere/elastic/v7
* mysql：https://github.com/go-sql-driver/mysql
* redis：https://github.com/go-redis/redis
* 文件存储：阿里云 oss、cdn
* markdown editor：https://github.com/pandao/editor.md
* pprof 性能调优
* 包管理工具 [Go Modules](https://github.com/golang/go/wiki/Modules)
* 后台登录：cookie 
* 使用 make 来管理 Go 工程
* 使用 shell(startu.sh) 脚本来管理进程
* 使用 YAML 文件进行多环境配置

### 目录结构

```shell
├── Makefile                     # 项目管理文件
├── conf                         # 配置文件统一存放目录
├── docs                         # 框架相关文档
├── internal                     # 业务目录
│   ├── handler                  # http 接口
│   ├── pkg                      # 内部应用程序代码
│   ├── routers                  # 业务路由
│   └── task                     # 异步任务
├── logs                         # 存放日志的目录
├── main.go                      # 项目入口文件
├── pkg                          # 公共的 package
├── tests                        # 单元测试依赖的配置文件，主要是供docker使用的一些环境配置文件
└── build                        # 存放用于执行各种构建，安装，分析等操作的脚本
```

### 功能模块

#### 后台
* 文章管理
* 页面管理
* 分类管理
* 标签管理
  
#### 前台
* 文章列表
* 内容页面
* 标签页面
* 关于页面
* 站内搜索

### 常用命令

- make help 查看帮助
- make dep 下载 Go 依赖包
- make build 编译项目
- make tar 打包文件

### 部署流程
* 依赖环境：
  * * mysql 
  * * redis 
  * * elasticsearch
  
* 克隆仓库

```
# 下载安装，可以不用是 GOPATH
git clone https://github.com/convee/goblog.git

# 进入到下载目录
cd goblog

# 生成环境配置文件
cd conf

# 修改 mysql、redis、elasticsearch 配置
mysql -u root -p
> create database blog;
> use blog;
> source blog.sql;


# 下载依赖
make dep

# 编译
make build

# 运行
./goblog dev.yml
```

* 访问首页

http://localhost:9091

* 访问后台

http://localhost:9091/admin
  
用户名：convee.@admin.cn
  
密码：123456
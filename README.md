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
* 使用 shell(startup.sh) 脚本来管理进程
* 使用 YAML 文件进行多环境配置

### 目录结构

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

### 功能模块

#### 后台
* 文章管理：文章增删改查
* 页面管理：页面增删改查，可自定义 markdown 页面
* 分类管理：分类增删改查
* 标签管理：标签列表
  
#### 前台
* 文章列表：倒序展示文章、可置顶
* 内容页面：markdown 内容展示
* 标签页面：按标签文章数量排序
* 关于页面：个人说明
* 阅读清单：个人阅读书籍
* 站内搜索：支持文章标题、描述、内容、分类、标签模糊搜索

## 开发规范

遵循: [Uber Go 语言编码规范](https://github.com/uber-go/guide/blob/master/style.md)

### 常用命令

- make help 查看帮助
- make dep 下载 Go 依赖包
- make build 编译项目
- make tar 打包文件

### 部署流程
* 依赖环境：
  
   mysql、redis、elasticsearch
   > elasticsearch 可通过配置开启关闭，redis主要考虑到后续加缓存
  
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

# 运行
./goblog dev.yml

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
  
用户名：convee.@admin.cn
  
密码：123456

* 演示站：https://convee.cn
## Go Markdown 博客系统

### 技术架构
* 前端框架：bootstrap
* 语言：go
* go 依赖管理：go module
* 网络库：标准库 net/http
* 配置读取：https://github.com/spf13/viper
* 日志库zap：https://github.com/uber-go/zap
* elasticsearch：https://github.com/olivere/elastic/v7
* mysql：https://github.com/go-sql-driver/mysql
* redis：https://github.com/go-redis/redis
* 文件存储：阿里云 oss、cdn
* 压测：wrk
* markdown editor：https://github.com/pandao/editor.md
* pprof 性能调优
* 后台登录：cookie 

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


### 部署流程
* 依赖环境：
  * * mysql 
  * * redis 
  * * elasticsearch
  
* 克隆仓库
```
git clone https://github.com/convee/blog
```
* 导入 sql：
```
mysql -u root -p
> create database blog;
> use blog;
> source blog.sql;
```
* 修改配置
```
路径：conf/dev.yml
修改：项目根路径、CDN路径、mysql、redis、es 配置
```
* 下载依赖
```
go get
```
* 运行
```
go run main.go
```

* 访问首页

http://localhost:9090

* 访问后台

http://localhost:9090/admin
  
用户名：convee.@admin.cn
  
密码：123456
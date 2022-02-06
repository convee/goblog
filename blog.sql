-- MySQL dump 10.13  Distrib 8.0.27, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: blog
-- ------------------------------------------------------
-- Server version	8.0.27-0ubuntu0.20.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `category`
--

DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `category`
--

/*!40000 ALTER TABLE `category` DISABLE KEYS */;
INSERT INTO `category` VALUES (1,'技术','2021-08-15 06:09:06','2021-09-05 11:17:53'),(2,'产品','2021-08-15 06:09:31','2021-08-15 06:09:35'),(3,'生活笔记','2021-08-15 06:09:42','2021-08-15 06:09:42');
/*!40000 ALTER TABLE `category` ENABLE KEYS */;

--
-- Table structure for table `page`
--

DROP TABLE IF EXISTS `page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `page` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ident` varchar(20) NOT NULL DEFAULT '''''',
  `title` varchar(100) NOT NULL,
  `content` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `page`
--

/*!40000 ALTER TABLE `page` DISABLE KEYS */;
INSERT INTO `page` VALUES (1,'about','关于我','### 为梦想而努力！\r\n#### 一线互联网码农，热衷探究技术背后的原理。喜欢情景喜剧、相声、小品，阅读，终生学习者。\r\n\r\n### 其他平台\r\n* Github：https://github.com/convee\r\n* 知乎主页：https://www.zhihu.com/people/convee/activities\r\n* 微信公众号：穿西装的程序员\r\n\r\n### 联系我\r\n* Gmail: convee.cn@gmail.com'),(2,'read','阅读清单','### 为梦想而努力！\r\n| 序号  | 书籍  | 进度  |\r\n| ------------ | ------------ |\r\n| 1 | 《高性能MySQL》 |  20% |\r\n| 2 | 《PHP7内核剖析》  |  30% |\r\n| 3 | 《ElasticSearch源码解析与优化实践》  |  30% |\r\n| 4 | 《深入理解Kafka》  |  30% |\r\n| 5 | 《Go程序设计语言》  |  30% |\r\n| 6 | 《labuladong的算法小抄》  |  30% |\r\n');
/*!40000 ALTER TABLE `page` ENABLE KEYS */;

--
-- Table structure for table `post`
--

DROP TABLE IF EXISTS `post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `category_id` int NOT NULL DEFAULT '0',
  `is_top` tinyint(1) NOT NULL DEFAULT '0',
  `tag_ids` json NOT NULL,
  `views` int NOT NULL DEFAULT '0',
  `description` varchar(2000) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post`
--

/*!40000 ALTER TABLE `post` DISABLE KEYS */;
/*!40000 ALTER TABLE `post` ENABLE KEYS */;

--
-- Table structure for table `post_content`
--

DROP TABLE IF EXISTS `post_content`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `post_content` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `content` longtext,
  `post_id` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `post_content`
--

/*!40000 ALTER TABLE `post_content` DISABLE KEYS */;
INSERT INTO `post_content` VALUES (1,'## Go Markdown 博客系统\r\n> 基于 Go 语言实现的 Markdown 博客系统\r\n\r\n### 技术栈\r\n\r\n* 前端框架：bootstrap\r\n* 语言：go\r\n* 网络库：标准库 net/http\r\n* 配置文件解析库 [Viper](https://github.com/spf13/viper)\r\n* 日志库zap：https://github.com/uber-go/zap\r\n* elasticsearch：https://github.com/olivere/elastic/v7\r\n* mysql：https://github.com/go-sql-driver/mysql\r\n* redis：https://github.com/go-redis/redis\r\n* 文件存储：阿里云 oss、cdn\r\n* markdown editor：https://github.com/pandao/editor.md\r\n* pprof 性能调优\r\n* 包管理工具 [Go Modules](https://github.com/golang/go/wiki/Modules)\r\n* 后台登录：cookie \r\n* 使用 make 来管理 Go 工程\r\n* 使用 shell(startu.sh) 脚本来管理进程\r\n* 使用 YAML 文件进行多环境配置\r\n\r\n### 目录结构\r\n\r\n```shell\r\n├── Makefile                     # 项目管理文件\r\n├── conf                         # 配置文件统一存放目录\r\n├── docs                         # 框架相关文档\r\n├── internal                     # 业务目录\r\n│   ├── handler                  # http 接口\r\n│   ├── pkg                      # 内部应用程序代码\r\n│   ├── routers                  # 业务路由\r\n│   └── task                     # 异步任务\r\n├── logs                         # 存放日志的目录\r\n├── main.go                      # 项目入口文件\r\n├── pkg                          # 公共的 package\r\n├── tests                        # 单元测试依赖的配置文件，主要是供docker使用的一些环境配置文件\r\n└── build                        # 存放用于执行各种构建，安装，分析等操作的脚本\r\n```\r\n\r\n### 功能模块\r\n\r\n#### 后台\r\n* 文章管理\r\n* 页面管理\r\n* 分类管理\r\n* 标签管理\r\n  \r\n#### 前台\r\n* 文章列表\r\n* 内容页面\r\n* 标签页面\r\n* 关于页面\r\n* 站内搜索\r\n\r\n### 常用命令\r\n\r\n- make help 查看帮助\r\n- make dep 下载 Go 依赖包\r\n- make build 编译项目\r\n- make tar 打包文件\r\n\r\n### 部署流程\r\n* 依赖环境：\r\n  * * mysql \r\n  * * redis \r\n  * * elasticsearch\r\n  \r\n* 克隆仓库\r\n\r\n```\r\n# 下载安装，可以不用是 GOPATH\r\ngit clone https://github.com/convee/goblog.git\r\n\r\n# 进入到下载目录\r\ncd goblog\r\n\r\n# 生成环境配置文件\r\ncd conf\r\n\r\n# 修改 mysql、redis、elasticsearch 配置\r\nmysql -u root -p\r\n> create database blog;\r\n> use blog;\r\n> source blog.sql;\r\n\r\n\r\n# 下载依赖\r\nmake dep\r\n\r\n# 编译\r\nmake build\r\n\r\n# 运行\r\n./goblog dev.yml\r\n```\r\n\r\n* 访问首页\r\n\r\nhttp://localhost:9091\r\n\r\n* 访问后台\r\n\r\nhttp://localhost:9091/admin\r\n  \r\n用户名：convee.@admin.cn\r\n  \r\n密码：123456',1),(18,'### 1. 下载源码并解压\r\n\r\n```\r\nwget wget https://download.redis.io/releases/redis-6.2.6.tar.gz\r\n\r\ntar -zxvf redis-6.2.6.tar.gz\r\n\r\ncd redis-6.2.6\r\n\r\n```\r\n\r\n### 2. 编译\r\n\r\n```\r\nmake\r\n```\r\n\r\n### 3. 安装\r\n\r\n```\r\nsudo make PREFIX=/usr/local/redis install\r\n\r\nsudo mkdir /usr/local/redis/etc/\r\n\r\nsudo cp redis.conf /usr/local/redis/etc/\r\n\r\ncd /usr/local/redis/bin/\r\n\r\nsudo cp redis-benchmark redis-cli redis-server /usr/local/bin/\r\n\r\n```\r\n\r\n### 4. 更改配置\r\n\r\nsudo cp /usr/local/redis/etc/redis.conf /usr/local/redis/etc/6379.conf\r\nsudo vim /usr/local/redis/etc/redis.conf\r\n\r\n```\r\ndaemonize yes\r\n\r\ntimeout 300\r\n\r\nloglevel verbose\r\n\r\nlogfile stdout\r\n\r\nbind 127.0.0.1\r\n```\r\n\r\n### 5. 配置环境变量\r\n\r\n```\r\nsudo vim /etc/profile\r\n\r\nexport PATH=\"$PATH:/usr/local/redis/bin\"\r\n\r\nsource /etc/profile\r\n\r\n```\r\n\r\n### 6. 配置启动脚本\r\n\r\n基于编译目录 utils/redis_init_script 脚本修改\r\n\r\nsudo vim /etc/init.d/redis\r\n\r\n```\r\n#!/bin/sh\r\n#\r\n# Simple Redis init.d script conceived to work on Linux systems\r\n# as it does use of the /proc filesystem.\r\n\r\n### BEGIN INIT INFO\r\n# Provides:     redis_6379\r\n# Default-Start:        2 3 4 5\r\n# Default-Stop:         0 1 6\r\n# Short-Description:    Redis data structure server\r\n# Description:          Redis data structure server. See https://redis.io\r\n### END INIT INFO\r\n\r\nREDISPORT=6379\r\nEXEC=/usr/local/redis/bin/redis-server\r\nCLIEXEC=/usr/local/redis/bin/redis-cli\r\n\r\nPIDFILE=/usr/local/redis/var/redis_${REDISPORT}.pid\r\nCONF=\"/usr/local/redis/etc/redis/${REDISPORT}.conf\"\r\n\r\ncase \"$1\" in\r\n    start)\r\n        if [ -f $PIDFILE ]\r\n        then\r\n                echo \"$PIDFILE exists, process is already running or crashed\"\r\n        else\r\n                echo \"Starting Redis server...\"\r\n                $EXEC $CONF\r\n        fi\r\n        ;;\r\n    stop)\r\n        if [ ! -f $PIDFILE ]\r\n        then\r\n                echo \"$PIDFILE does not exist, process is not running\"\r\n        else\r\n                PID=$(cat $PIDFILE)\r\n                echo \"Stopping ...\"\r\n                $CLIEXEC -p $REDISPORT shutdown\r\n                while [ -x /proc/${PID} ]\r\n                do\r\n                    echo \"Waiting for Redis to shutdown ...\"\r\n                    sleep 1\r\n                done\r\n                echo \"Redis stopped\"\r\n        fi\r\n        ;;\r\n    *)\r\n        echo \"Please use start or stop as first argument\"\r\n        ;;\r\nesac\r\n```\r\n\r\n### 7. 开启自启动设置\r\n\r\n```\r\nsudo chmod +x /etc/init.d/redis\r\n```\r\n\r\n### 8. 启动测试\r\n\r\n```\r\nsudo /etc/init.d/redis start\r\n\r\nps -ef|grep redis\r\n\r\nnetstat -an|grep 6379\r\n```',11),(19,'## 1.安装依赖包\r\n\r\n```bash\r\nyum install -y gcc gcc-c++  make zlib zlib-devel pcre pcre-devel  libjpeg libjpeg-devel libpng libpng-devel freetype freetype-devel gdbm-devel db4-devel libXpm-devel libX11-devel gd-devel gmp gmp-devel expat-devel libxml2 libxml2-devel glibc glibc-devel glib2 glib2-devel bzip2 bzip2-devel ncurses ncurses-devel libcurl libcurl-devel curl curl-devel libmcrypt libmcrypt-devel libxslt libxslt-devel xmlrpc-c xmlrpc-c-devel libicu-devel libmemcached-devel libzip readline readline-devel e2fsprogs e2fsprogs-devel krb5 krb5-devel openssl openssl-devel openldap openldap-devel nss_ldap openldap-clients openldap-servers\r\n```\r\n\r\n## 2.下载 PHP\r\n\r\n```bash\r\nwget https://www.php.net/distributions/php-7.3.20.tar.gz\r\n```\r\n\r\n## 3.解压\r\n\r\n```bash\r\ntar -xzvf php-7.3.20.tar.gz\r\ncd php-7.3.20\r\n```\r\n\r\n## 4.配置\r\n\r\n```bash\r\n./configure --prefix=/usr/local/php7 --with-config-file-path=/usr/local/php7/etc --enable-inline-optimization --disable-debug --enable-fpm --with-fpm-user=www --with-fpm-group=www --disable-rpath --enable-soap --with-libxml-dir --with-xmlrpc --with-openssl  --with-mhash --with-pcre-regex --with-zlib --enable-bcmath --with-bz2 --enable-calendar --with-curl --enable-exif --with-pcre-dir --enable-ftp --with-gd --with-openssl-dir --with-jpeg-dir --with-png-dir --with-zlib-dir --with-freetype-dir --enable-gd-jis-conv --with-gettext --with-gmp --with-mhash --enable-mbstring --with-onig --enable-shared --enable-opcache --with-mysqli=mysqlnd --with-pdo-mysql=mysqlnd --with-readline --with-iconv --enable-pcntl --enable-shmop --enable-sysvmsg --enable-sysvsem --enable-sysvshm --enable-sockets  --enable-zip --enable-wddx --with-pear\r\n\r\nmake && make install\r\n```\r\n\r\n## 5.创建www用户\r\n\r\n```bash\r\ngroupadd www #添加 www 用户组\r\nuseradd -g www www #添加 www 用户到 www 用户组\r\n```\r\n\r\n## 6.初始化 php-fpm 配置\r\n\r\n```bash\r\n# 复制 php.ini\r\ncp php-7.2.20/php.ini-production /usr/local/php7/etc/php.ini\r\n\r\n# 增加执行权限\r\nchmod +x /etc/init.d/php-fpm\r\n# 配置 php-fpm 文件\r\ncd /usr/local/php7/etc/\r\ncp php-fpm.conf.default php-fpm.conf\r\n# 进入 php-fpm.conf , 并去除 pid = run/php-fpm.pid 的注释\r\nvim php-fpm.conf\r\n\r\n# 复制 www.conf 文件\r\ncp php-fpm.d/www.conf.default php-fpm.d/www.conf\r\n```\r\n\r\n## 7.php-fpm 启动脚本\r\n\r\n- /etc/init.d/php-fpm stop # 停止服务\r\n- /etc/init.d/php-fpm start # 启动服务\r\n\r\n- /etc/init.d/php-fpm restart # 重启服务\r\n\r\n```bash\r\n# 复制启动脚本到 init.d 目录\r\ncp php-7.2.20/sapi/fpm/init.d.php-fpm /etc/init.d/php-fpm\r\n```\r\n\r\n## 8.centos 管理服务化\r\n\r\n- systemctl enable xxxxxx # 配置自启动\r\n- systemctl stop xxxxx # 停止服务\r\n\r\n- systemctl start xxxx # 开启服务\r\n- systemctl status xxxx # 查看状态\r\n\r\n```bash\r\n# 在 centos 7 之后我们可以使用 systemctl 更好的管理系统服务\r\n# 所以我们也要让 php-fpm 支持\r\n# 因为 php 7.2 源码包里面含有 systemctl 所需要的脚本文件\r\n# 我们只要复制过去即可,我们来开始配置\r\n# 进入下载的 php源码包\r\n$ cd php-7.2.20/sapi/fpm\r\n# 复制其中的 php-fpm.service 到 /usr/lib/systemd/system/\r\n$ cp php-fpm.service /usr/lib/systemd/system/\r\n# 再次使用 systemctl enable php-fpm 进行配置自启动\r\n$ systemctl enable php-fpm\r\n# 重启测试一下看看自己服务器的 php-fpm 是否成功运行\r\n```\r\n\r\n## 9.php别名设置\r\n\r\n- 新增 /usr/bin/php 文件\r\n- chmod u+x /usr/bin/php\r\n\r\n```bash\r\n#!/bin/bash\r\n/usr/local/php7/bin/php $1\r\n```\r\n\r\n- ln -s /usr/local/php7/bin/php /usr/bin/php',12),(20,'## tcpdump 简介：\r\n\r\n\r\n```\r\ntcpdump是linux下的网络抓包工具。支持针对网络层，协议，主机，网络或者端口的过滤，并提供and or not等逻辑语句来去掉无用信息。\r\n```\r\n\r\n\r\n## 安装：\r\n\r\n\r\n- ubuntu:\r\n\r\n\r\n\r\n```\r\n    $ sudo apt-get install tcpdump\r\n```\r\n\r\n\r\n- centos:\r\n\r\n\r\n\r\n```\r\n    $ yum install tcpdump\r\n```\r\n\r\n\r\n> tips: 如果安装失败，换一下国内的源即可，或者vpn下。\r\n\r\n\r\n\r\n## 关键字：\r\n\r\n\r\n- 类型关键字：host, net, port\r\n- 传输方向的关键字：src, dst, dst or src, dst and src.这些 关键字指明了传输的方向。缺省是src or dst\r\n- 协议的关键字：主要包括fddi, ip, arp, rarp, tcp, udp等类型。\r\n- 其他关键字：gateway, broadcast, less, greater。\r\n- 逻辑运算：!, ||, AND\r\n\r\n\r\n\r\n查看可以监听的网络接口：\r\n\r\n\r\n```\r\n[root@iZuf69su1tsn6fy07g4262Z ~]# tcpdump -D\r\n1.eth0 [Up, Running]\r\n2.lo [Up, Running, Loopback]\r\n3.any (Pseudo-device that captures on all interfaces) [Up, Running]\r\n4.bluetooth-monitor (Bluetooth Linux Monitor) [none]\r\n5.nflog (Linux netfilter log (NFLOG) interface) [none]\r\n6.nfqueue (Linux netfilter queue (NFQUEUE) interface) [none]\r\n7.usbmon0 (All USB buses) [none]\r\n8.usbmon1 (USB bus number 1)\r\n```\r\n\r\n\r\n## tcpdump命令：\r\n\r\n\r\n```\r\n$ tcpdump option filter\r\n    option: 选项\r\n    filter: 过滤条件\r\n\r\n$ tcpdump -Xnlps0 -nn -iany  host  192.168.5.1\r\n  option: -Xnlps0 -nn -iany (选项)\r\n  filter: host 192.168.5.1 (参数)\r\n```\r\n\r\n\r\noption:\r\n\r\n\r\n-i [网络接口]： 监听哪个网络接口，默认监听的是-D显示的第一个网络接口。\r\n\r\n\r\n-n：显示主机ip和端口。\r\n\r\n\r\n-w [file]：把抓包的结果写入到file中。\r\n\r\n\r\n-C [size(MB)]：配合-w使用，表示file大小，当抓包结果>size,自动生成新的文件。\r\n\r\n\r\n-X：以16进制以及ASCII的形式打印出数据内容\r\n\r\n\r\n-x：除了打印出header外，还打印packat中的数据\r\n\r\n\r\n-xx：以16进制的形式打印header，packet里面的数据\r\n\r\n\r\n-A: 把每一个packet都用ASCII的形式打印出来\r\n\r\n\r\n-c：收到3个package就退出\r\n\r\n\r\n-e：把连接层的头打印出来\r\n\r\n\r\nfilter:\r\n\r\n\r\n1.  协议关键字：\r\n    fddi, ip, arp, rarp \r\n1.  源或目标主机：\r\n    src, src or dst, src and dst \r\n\r\n---\r\n\r\n实例说明\r\n\r\n\r\n本机环境说明：\r\n\r\n\r\nPHP5.6 docke，ip为192.168.5.1\r\n\r\n\r\nnginx  docker容器，ip为172.18.0.4\r\n\r\n\r\n其他服务, ip为192.168.5.1\r\n\r\n\r\n- 例子1： 访问本机一个页面，抓取数据包\r\n\r\n\r\n\r\n```\r\nxxxx@xxx:~/Documents$ sudo tcpdump -nn  host 192.168.5.1\r\ntcpdump: verbose output suppressed, use -v or -vv for full protocol decode\r\nlistening on br-ba99820d8dbf, link-type EN10MB (Ethernet), capture size 262144 bytes\r\n10:54:04.448134 IP 172.18.0.4.37242 > 172.18.0.2.9000: Flags [S], seq 2092336432, win 29200, options [mss 1460,sackOK,TS val 43377964 ecr 0,nop,wscale 7], length 0\r\n\r\n说明：\r\n10:54:04.448134 IP：代表时间 和 数据包类型\r\n172.18.0.4.37242 > 172.18.0.2.9000 表示数据包流向 \r\nFlags[S(建立连接)|P(发送数据)|F(完成)|R(重置)|.(none)]\r\nseq: 包序号\r\nack: 确认包\r\nwin: 数据窗口大小(接受窗口大小)\r\nOPTIONS：tcp包 中的选项字段\r\nMSS： 最大分段大小(发送数据的最大长度)\r\n```\r\n\r\n\r\n- 例子二：抓取整个请求到响应整个流程的包数据\r\n\r\n\r\n\r\n```\r\nxxxx@xxx:~/Documents$ tcpdump -nn  host 172.18.0.2 or 172.18.0.4\r\n说明： 抓取php56 和 nginx两台docker的数据包\r\n```\r\n\r\n\r\n## 常用的几个抓包命令：\r\n\r\n\r\n1. 获取指定IP上的所有请求，无包数据\r\n\r\n\r\n\r\n```\r\n$ tcpdump -nn host :IP\r\n```\r\n\r\n\r\n2. 获取指定IP上的所有请求，有包数据\r\n\r\n\r\n\r\n```\r\n$ tcpdump -nnX host :IP\r\n```\r\n\r\n\r\n3. 获取源IP发送到目的IP的所有包\r\n\r\n\r\n\r\n```\r\n$ tcpdump -nn src host :IP and \\( dst host :IP and port :NUM \\)\r\n```\r\n\r\n\r\n4. 获取 两个IP之间的相互通信\r\n\r\n\r\n\r\n```\r\n$ tcpdump -nn host :IP and host :IP\r\n```\r\n\r\n\r\n5. 获取IP发送了哪些包\r\n\r\n\r\n\r\n```\r\n$ tcpdump -nn src host :IP\r\n```',13),(21,'## PHP vs Go\r\n\r\n## 语言对比\r\n\r\n| 语言 | PHP | Go    |\r\n| ----|----|----- |\r\n| 特点 | 脚本语言、开发快、部署方便、性能一般 | 静态语言、部署方便、性能好 |\r\n| 部署  | 1、源码部署、FPM 依赖 Nginx 之类的服务器 2、对于长驻 Http Server 方面，Workerman并未普及，Swoole 贴近底层而面临未来版本升级、维护成本高 | 1、既可以直接 run 文件运行，也可以编译成二进制运行，既方便也安全  2、不依赖 Web 服务器运行，单文件就可以启动高性能 Http Server，资源消耗极少 | \r\n| 性能  | 中小型应用完全可以接受，但一旦服务器数量增多，性能差距就会直接暴漏出来 | 21世纪的C语言 |\r\n| 服务端编程  | 只有 Pcntl 用于多进程编程，比较简陋，Pthreads 多线程不安全，网络编程门槛高，多数在结合框架基础上才能够做一些工作 | 官方内置很多网络库，底层不依赖C/C++ 实现，协程编程模型的 CPU 效率高、易用 |\r\n\r\n## 性能对比\r\n\r\n服务器配置：8核 16G内存\r\n业务场景：获取订单200000 21年11月11日的曝光量，使用 redis hget 命令\r\n压测：使用 wrk 测试，开16个线程，1000个连接\r\n​\r\n\r\n软件：nginx+fpm\r\n代码：\r\n\r\n```php\r\n<?php\r\n\r\n\r\n$redis = new \\Redis();\r\n\r\n$redis->connect(\'127.0.0.1\', 6379);\r\n\r\n$pv = $redis->hGet(\'PV_COUNT_200000\', 211111);\r\n\r\necho $pv;\r\n\r\n```\r\n\r\n命令：wrk -t16 -c1000 -d10s --latency [http://192.168.5.100/pv.php](http://192.168.5.100/pv.php)\r\n结果：\r\n\r\n```bash\r\nroot@bj-ali-adserving-goproxy-test-5-32:~# wrk -t16 -c1000 -d10s --latency http://192.168.5.100/pv.php\r\nRunning 10s test @ http://192.168.5.100/pv.php\r\n  16 threads and 1000 connections\r\n  Thread Stats   Avg      Stdev     Max   +/- Stdev\r\n    Latency   101.57ms    4.86ms 138.95ms   91.92%\r\n    Req/Sec   611.39     53.87     1.27k    88.71%\r\n  Latency Distribution\r\n     50%  101.87ms\r\n     75%  103.61ms\r\n     90%  105.00ms\r\n     99%  108.48ms\r\n  97708 requests in 10.09s, 19.19MB read\r\nRequests/sec:   9682.76\r\nTransfer/sec:      1.90MB\r\n\r\n```\r\n\r\n软件：go\r\n代码：\r\n\r\n```go\r\npackage main\r\n\r\nimport (\r\n	\"context\"\r\n	\"log\"\r\n	\"net/http\"\r\n	\"strconv\"\r\n\r\n	\"github.com/go-redis/redis/v8\"\r\n)\r\n\r\nfunc pv(w http.ResponseWriter, r *http.Request) {\r\n	pv, _ := redisClient.HGet(context.Background(), \"PV_COUNT_200000\", \"211111\").Int()\r\n	log.Println(100)\r\n	w.Write([]byte(strconv.Itoa(pv)))\r\n}\r\n\r\n// RedisClient redis 客户端\r\nvar redisClient *redis.Client\r\n\r\nfunc init() {\r\n	redisClient = redis.NewClient(&redis.Options{\r\n		Addr:         \"127.0.0.1:6379\",\r\n		Password:     \"\",\r\n		DB:           0,\r\n		MinIdleConns: 200,\r\n		PoolSize:     100,\r\n	})\r\n\r\n	_, err := redisClient.Ping(context.Background()).Result()\r\n	if err != nil {\r\n		panic(err)\r\n	}\r\n}\r\n\r\nfunc main() {\r\n	http.HandleFunc(\"/pv.go\", pv)\r\n	log.Println(\"http listen 9091\")\r\n	err := http.ListenAndServe(\"0.0.0.0:9091\", nil)\r\n	if err != nil {\r\n		log.Fatalln(\"http listen failed\")\r\n	}\r\n}\r\n\r\n```\r\n\r\n命令：wrk -t16 -c1000 -d10s --latency [http://192.168.5.100:9091/pv.go](http://192.168.5.100:9091/pv.go)\r\n\r\n```bash\r\nroot@bj-ali-adserving-goproxy-test-5-32:~# wrk -t16 -c1000 -d10s --latency http://192.168.5.100:9091/pv.go\r\nRunning 10s test @ http://192.168.5.100:9091/pv.go\r\n  16 threads and 1000 connections\r\n  Thread Stats   Avg      Stdev     Max   +/- Stdev\r\n    Latency    11.14ms    3.26ms  54.48ms   85.59%\r\n    Req/Sec     5.61k   837.60    14.14k    78.43%\r\n  Latency Distribution\r\n     50%   10.20ms\r\n     75%   11.81ms\r\n     90%   15.29ms\r\n     99%   22.46ms\r\n  898298 requests in 10.10s, 101.95MB read\r\nRequests/sec:  88946.75\r\nTransfer/sec:     10.09MB\r\n```\r\n\r\n对比：\r\n\r\n| 语言       | PHP       | Go       |\r\n| ---------- | --------- | -------- |\r\n| QPS        | 9682.76   | 88946.75 |\r\n| CPU 利用率 |           |          |\r\n| 延时       | 平均101ms | 平均11ms |\r\n\r\n\r\n\r\n',14),(22,'\r\n\r\n## 一. 创建用户\r\n\r\n### 命令\r\n\r\n```sql\r\nCREATE USER \'username\'@\'host\' IDENTIFIED BY \'password\';\r\n```\r\n\r\n### 说明\r\n\r\n- username：用户名\r\n- host：可登录的主机\r\n- password：登录密码\r\n\r\n### 例子\r\n\r\n```sql\r\nCREATE USER \'username\'@\'localhost\' IDENTIFIED BY \'123456\';\r\nCREATE USER \'username\'@\'%\' IDENTIFIED BY \'123456\';\r\n```\r\n\r\n## 二. 授权\r\n\r\n### 命令\r\n\r\n```sql\r\nGRANT privileges ON dbname.tablename TO \'username\'@\'host\'\r\n```\r\n\r\n### 说明\r\n\r\n- privileges：用户的操作权限，如SELECT，INSERT，UPDATE等，如果要授予所的权限则使用ALL\r\n- dbname：数据库名\r\n- tablename：表名，如果要授予该用户对所有数据库和表的相应操作权限则可用*表示，如*.*\r\n\r\n### 例子\r\n\r\n```sql\r\nGRANT SELECT, INSERT ON test.user TO \'username\'@\'%\';\r\nGRANT ALL ON *.* TO \'username\'@\'%\';\r\n```\r\n\r\n### \r\n\r\n## 三.设置与更改用户密码\r\n\r\n### 命令\r\n\r\n```sql\r\nSET PASSWORD FOR \'username\'@\'host\' = PASSWORD(\'newpassword\');\r\n```\r\n\r\n如果是当前登陆用户用:\r\n\r\n```sql\r\nSET PASSWORD = PASSWORD(\"newpassword\");\r\n```\r\n\r\n### 例子\r\n\r\n```sql\r\nSET PASSWORD FOR \'username\'@\'%\' = PASSWORD(\"123456\");\r\n```\r\n\r\n## 四. 撤销用户权限\r\n\r\n### 命令\r\n\r\n```sql\r\nREVOKE privilege ON dbname.tablename FROM \'username\'@\'host\';\r\n```\r\n\r\n## 说明\r\n\r\nprivilege, dbname, tablename：同授权部分\r\n\r\n### 例子\r\n\r\n```sql\r\nREVOKE SELECT ON *.* FROM \'username\'@\'%\';\r\n```\r\n\r\n## 五.删除用户\r\n\r\n### 命令\r\n\r\n```sql\r\nDROP USER \'username\'@\'host\';\r\n```',15),(23,'### 1.安装 PHP\r\n\r\n```bash\r\nbrew install php@7.3\r\n```\r\n\r\n### 2.安装 Protobuf 扩展\r\n\r\n```bash\r\nwget https://github.com/allegro/php-protobuf/archive/master.zip\r\n \r\nunzip master.zip\r\n \r\ncd php-protobuf-master\r\n \r\nsudo /usr/local/opt/php@7.3/bin/phpize\r\n \r\nsudo ./configure --prefix=/usr/local/opt/php@7.3/bin/php --with-php-config=/usr/local/opt/php@7.3/bin/php-config\r\n \r\nmake && make install\r\n \r\n//然后在php.ini里面加一下extension = \"protobuf.so\"\r\n```\r\n\r\n### 3.安装 Composer\r\n\r\n```bash\r\ncd php-protobuf-master\r\n\r\ncurl -s http://getcomposer.org/installer | php\r\n\r\n/usr/local/opt/php@7.3/bin/php composer.phar install\r\n```\r\n\r\n### 4. Protobuf 使用\r\n\r\n```bash\r\n/usr/local/opt/php@7.3/bin/php ./php-protobuf-master/protoc-gen-php.php test.proto\r\n```',16),(24,'php',17);
/*!40000 ALTER TABLE `post_content` ENABLE KEYS */;

--
-- Table structure for table `tag`
--

DROP TABLE IF EXISTS `tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tag` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `count` int NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tag`
--

/*!40000 ALTER TABLE `tag` DISABLE KEYS */;
INSERT INTO `tag` VALUES (1,'go1',3,'2021-09-07 16:18:11','2021-11-24 16:01:48'),(2,'markdown',1,'2021-09-07 16:18:11','2021-11-13 15:59:44'),(20,'redis',1,'2021-11-06 13:46:18','2021-11-13 15:59:44'),(21,'go12',4,'2021-11-07 12:52:12','2021-11-24 16:09:12'),(22,'tcpdump',1,'2021-11-10 15:52:18','2021-11-13 15:59:44'),(24,'mysql',1,'2021-11-13 12:01:57','2021-11-13 15:59:44'),(25,'protobuf',1,'2021-11-13 14:46:37','2021-11-13 15:59:44');
/*!40000 ALTER TABLE `tag` ENABLE KEYS */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(50) NOT NULL DEFAULT '',
  `password` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'admin@convee.cn','$2a$08$kYT2DJz9W7Yv5AWv0swq6e0qgVyw9sOkKWQ5s20HBohXJLuYf4ZpG');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-02-06 15:58:28

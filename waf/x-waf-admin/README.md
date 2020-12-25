# X-WAF

X-WAF是一款适用中、小企业的云WAF系统，让中、小企业也可以非常方便地拥有自己的免费云WAF。

文档地址：[https://waf.xsec.io/](https://waf.xsec.io/)

# 主要特性

- 支持对常见WEB攻击的防御，如sql注入、xss、路径穿越，阻断扫描器的扫描等
- 对持对CC攻击的防御
- waf为反向模式，后端保护的服务器可直接用内网IP，不需暴露在公网中
- 支持IP、URL、Referer、User-Agent、Get、Post、Cookies参数型的防御策略
- 安装、部署与维护非常简单
- 支持在线管理waf规则
- 支持在线管理后端服务器
- 多台waf的配置可自动同步
- 跨平台，支持在linux、unix、mac和windows操作系统中部署

# 架构简介
x-waf由waf自身与Waf管理后台组成：

- [waf](https://github.com/xsec-lab/x-waf)：基于openresty + lua开发。
- [waf管理后台](https://github.com/xsec-lab/x-waf-admin)：采用golang + xorm + macrom开发的，支持二进制的形式部署。

waf和waf-admin必须同时部署在每一台云WAF服务器中。

# 下载安装
## waf安装
### centos平台

从[openresty](http://openresty.org/en/download.html)官方下载最新版本的源码包。

编译安装openresty：

```bash
yum -y install pcre pcre-devel
wget https://openresty.org/download/openresty-1.9.15.1.tar.gz
tar -zxvf openresty-1.9.15.1.tar.gz 
cd openresty-1.9.15.1
./configure 
gmake && gmake install

/usr/local/openresty/nginx/sbin/nginx  -t
nginx: the configuration file /usr/local/openresty/nginx/conf/nginx.conf syntax is ok
nginx: configuration file /usr/local/openresty/nginx/conf/nginx.conf test is successful
/usr/local/openresty/nginx/sbin/nginx 
```

### ubuntu平台安装

编译安装openresty：

```bash
apt-get install libreadline-dev libncurses5-dev libpcre3-dev libssl-dev perl make build-essential
sudo ln -s /sbin/ldconfig /usr/bin/ldconfig
wget https://openresty.org/download/openresty-1.9.15.1.tar.gz
tar -zxvf openresty-1.9.15.1.tar.gz
cd openresty-1.9.15.1
make && sudo make install
```

## 安装waf管理后台x-waf-admin

### 二进制安装

直接从github中下载对应操作系统的二进制版本，https://github.com/xsec-lab/x-waf-admin/releases

### 源码安装

-  首先需要搭建好go语言开发环境，可以参考[Go Web 编程](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.1.md)
- 安装依赖包
```bash
go get gopkg.in/macaron.v1
go get gopkg.in/ini.v1
go get github.com/go-sql-driver/mysql
go get github.com/go-xorm/xorm
go get github.com/xsec-lab/x-waf-admin
```

- 从github中下载最新的版本
- 执行go build server.go编译出二进制版本，然后将server、conf、publib和templates目录一起打包上传到服务器中即可运行。


## 后端服务器管理

当多台waf做负载均衡时，只需登录其中一台进行管理即可，多台waf的所有的配置信息会自动同步到所有的服务器中。

管理地址为：http://ip:5000/login/

管理后台的默认的账户及口令分别为：admin，x@xsec.io，请管理员部署系统后第一时间修改密码，防止被攻击者使用默认口令登入胡乱改动waf的配置。

### 新增站点

在`Site Manager`选项中，可以新增一个后端服务器，需要填写以下内容：

- Site Name，表示要加入waf的网站的域名
- 80表示该网站监听的端口
- Backend，表示有多少个后台app server，可以写多个（换行分割），例如：
```bash
1.1.1.1:80
8.8.8.8:80
```

- SSL Status，表示是否启用ssl，参数为on或off（如果要启用的话，需要在nginx中配置有效的证书）
- Debug Level，表示日志级别，可选的参数有`debug, info, notice, warn, error, crit, alert, emerg`

### 站点配置同步

新增站点后，需要在后台同步站点信息后方可生效，同步的方式有2种：

1. 全部同步
1. 针对某一新增的站点进行同步

## waf规则管理

在`waf Rules`选项中，可以修改waf的规则，修改完后可以点击“同步全部策略”按钮，将最新的规则同步到所有的服务器并让openresty重新加载。
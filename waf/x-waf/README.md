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

*  首先需要搭建好go语言开发环境，可以参考[Go Web 编程](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.1.md)
* 安装依赖包

```bash
go get gopkg.in/macaron.v1
go get gopkg.in/ini.v1
go get github.com/go-sql-driver/mysql
go get github.com/go-xorm/xorm
go get github.com/xsec-lab/x-waf-admin
```
* 从github中下载最新的版本
* 执行go build server.go编译出二进制版本，然后将server、conf、publib和templates目录一起打包上传到服务器中即可运行。

## 致谢

1. 感谢春哥开源的[openresty](https://openresty.org)
1. 感谢unixhot开源的[waf](https://github.com/unixhot/waf)
1. 感谢无闻开源的[macron](https://go-macaron.com/)和[peach](https://peachdocs.org/)
1. 感谢lunny开源的[xorm](https://github.com/go-xorm/xorm)

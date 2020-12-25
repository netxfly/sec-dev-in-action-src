## 代理蜜罐的开发与应用实战

本仓库为『代理蜜罐的开发与应用实战』的配套代码，原文地址：

### 代理蜜罐架构

![](http://docs.xsec.io/images/x-proxy//proxy_honeypot.png)

- 代理蜜罐Agent，提供代理服务，收集http请求与响应数据并发送到server集群
- 代理蜜罐Server（支持水平扩展），接收Agent传来的数据，对数据简单判断后入库
- 后端数据库（mongodb），存储代理蜜罐的数据
- 数据分析程序，对存数的数据进行加工处理，方便管理端展示
- 管理端，查看收集到的数据与数据分析结果

`server`与`manager`的运行依赖`mongodb`，可通过以下命令创建一个用户供`server`与`manager`使用：

```shell
db.createUser(
... {
...  user : "xproxy",
...  pwd : "xsec.io",
...  roles: [ { role : "readWrite", db : "xproxy" }]
... }
... )
```

管理端初次启动会添加一个默认用户，用户名为：`xproxy`，密码为：`x@xsec.io`


## 附录

### 参考资料

- [Is NordVPN a Honeypot?](http://vpnscam.com/is-nordvpn-a-honeypot/)
- [基于vpn和透明代理的web漏洞扫描器的实现思路及demo](https://github.com/netxfly/Transparent-Proxy-Scanner)

### Go语言学习资料

- [在2019成为一名Go开发者的路线图](https://github.com/Quorafind/golang-developer-roadmap-cn)
- [Go入门指南](https://github.com/Unknwon/the-way-to-go_ZH_CN)
- [Go WEB编程](https://github.com/astaxie/build-web-application-with-golang)
- [Go语言高级编程](https://github.com/chai2010/advanced-go-programming-book)
- [Go 语言学习资料与社区索引](https://github.com/Unknwon/go-study-index)
- [Go学习之路](https://github.com/developer-learning/learning-golang)
- [Go 101](https://gfw.go101.org/article/101.html)

### 用到的库与框架

- [goproxy](https://github.com/elazarl/goproxy)
- [cli](github.com/urfave/cli)
- [logrus](github.com/sirupsen/logrus)
- [macaron](https://github.com/go-macaron/macaron)
- [xorm](github.com/go-xorm/xorm)
- [upper.io](upper.io/db.v3)
- [mgo](gopkg.in/mgo.v2)

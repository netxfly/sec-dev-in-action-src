package vars

import (
	"github.com/patrickmn/go-cache"
	"gopkg.in/cheggaaa/pb.v2"

	"strings"
	"sync"
	"time"
)

var (
	// ip_list，用户名与密码字典
	IpList   = "ip_list.txt"
	UserDict = "user.dic"
	PassDict = "pass.dic"

	// 启动时间
	StartTime time.Time
	// 结果保存文件
	ResultFile = "password_crack.txt"

	// 超时时间
	TimeOut = 3 * time.Second
	// 默认协程数
	ScanNum = 5000

	DebugMode bool

	// 弱口令扫描进度条
	ProgressBar *pb.ProgressBar
	// 检测端口是否开放的进度条
	ProcessBarActive *pb.ProgressBar
)

var (
	// 扫描结果保存到一个cache中，该cache库支持内存数据落盘
	CacheService *cache.Cache

	PortNames = map[int]string{
		22:    "SSH",
		3306:  "MYSQL",
		6379:  "REDIS",
		1433:  "MSSQL",
		5432:  "POSTGRESQL",
		27017: "MONGODB",
	}

	// 标记特定服务的特定用户是否破解成功，成功的话不再尝试破解该用户
	SuccessHash sync.Map

	SupportProtocols map[string]bool
)

func init() {
	CacheService = cache.New(cache.NoExpiration, cache.DefaultExpiration)

	SupportProtocols = make(map[string]bool)
	for _, proto := range PortNames {
		SupportProtocols[strings.ToUpper(proto)] = true
	}

}

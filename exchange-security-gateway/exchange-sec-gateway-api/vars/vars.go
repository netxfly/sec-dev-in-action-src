package vars

import (
	"strings"

	"exchange_zero_trust_api/settings"

	"github.com/go-redis/redis"
)

type (
	RedisConfig struct {
		Host     string
		Port     int
		Db       int
		Password string
	}

	Prefix struct {
		UserPrefix    string // EAS的设备账户前缀
		AccountPrefix string // EAS的账户前缀
		CodePrefix    string // 激活码前缀

		DevicePrefix string
		DeviceInfo   string
		DeviceName   string

		EwsPrefix   string // ews设备前缀
		EwsUaPrefix string // ews ua前缀
	}
)

var (
	// redis配置
	RedisConf RedisConfig

	// redis instance
	RedisInstance *redis.Client

	// redis key prefix
	RedisKeyPrefix Prefix

	// 激活信息推送方式
	SendType string

	// 短信接口的配置
	SmsURL    string
	SmsKey    string
	SmsHeader string

	// 邮件配置
	MailServer string
	MailPort   int
	MailUser   string
	MailPass   string
	MailSuffix string

	// timeout
	Timeout = 3

	// 大象配置
	PubId        string // 公众号ID
	ClientId     string // 开发者信息的appKey
	ClientSecret string // 开发者信息的appToken
	Host         string // 服务器基地址
	UrlPath      string // url path

	//	微信配置
	CorpId     string // 企业ID
	AgentId    int    // 应用ID
	CorpSecret string // 应用secret
	CorpHost   string // 企业微信API的地址

	//	阿里钉钉配置
	OapiHost     string // OAPI
	AliAgentId   int    // 应用ID
	AliAppKey    string // app key
	AliAppSecret string // app secret
)

func init() {
	cfg := settings.Cfg

	SendType = cfg.Section("").Key("SEND_TYPE").MustString("neixin")
	SendType = strings.ToLower(SendType) // SendType的值统一为小写

	secRedis := settings.Cfg.Section("redis")
	RedisConf.Host = secRedis.Key("host").MustString("127.0.0.1")
	RedisConf.Port = secRedis.Key("port").MustInt(6379)
	RedisConf.Db = secRedis.Key("db").MustInt(0)
	RedisConf.Password = secRedis.Key("password").MustString("password")

	sec := cfg.Section("SMS")
	SmsURL = sec.Key("URL").MustString("")
	SmsHeader = sec.Key("HEADER").MustString("X-API-Token")
	SmsKey = sec.Key("KEY").MustString("")

	secMail := cfg.Section("MAIL")
	MailServer = secMail.Key("SERVER").MustString("mail.xsec.io")
	MailPort = secMail.Key("PORT").MustInt(25)
	MailUser = secMail.Key("USER").MustString("noreply")
	MailPass = secMail.Key("PASS").MustString("")
	MailSuffix = secMail.Key("SUFFIX").MustString("xsec.io")

	secNx := cfg.Section("NEIXIN")
	PubId = secNx.Key("PUB_ID").MustString("")
	ClientId = secNx.Key("CLIENT_ID").MustString("")
	ClientSecret = secNx.Key("CLIENT_SECRET").MustString("")
	Host = secNx.Key("HOST").MustString("")
	UrlPath = secNx.Key("URL_PATH").MustString("")

	secWeixin := cfg.Section("WEIXIN")
	CorpId = secWeixin.Key("CORP_ID").MustString("")
	AgentId = secWeixin.Key("AGENT_ID").MustInt(0)
	CorpSecret = secWeixin.Key("CORP_SECRET").MustString("")
	CorpHost = secWeixin.Key("CORP_HOST").MustString("")

	secDd := cfg.Section("DINGDING")
	OapiHost = secDd.Key("OAPI_HOST").MustString("")
	AliAgentId = secDd.Key("AGENT_ID").MustInt(0)
	AliAppKey = secDd.Key("APP_KEY").MustString("")
	AliAppSecret = secDd.Key("APP_SECRET").MustString("")

	prefixSec := cfg.Section("prefix")
	RedisKeyPrefix.UserPrefix = prefixSec.Key("user_prefix").MustString("exchange_user_")
	RedisKeyPrefix.AccountPrefix = prefixSec.Key("account_prefix").MustString("exchange_account_")
	RedisKeyPrefix.CodePrefix = prefixSec.Key("code_prefix").MustString("exchange_code_")
	RedisKeyPrefix.DevicePrefix = prefixSec.Key("device_prefix").MustString("exchange_device_")
	RedisKeyPrefix.DeviceInfo = prefixSec.Key("device_info").MustString("exchange_info_")
	RedisKeyPrefix.EwsPrefix = prefixSec.Key("ews_prefix").MustString("exchange_ews_")
	RedisKeyPrefix.EwsUaPrefix = prefixSec.Key("ews_ua_prefix").MustString("exchange_ews_ua_")
}

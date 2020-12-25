package plugins

import (
	"sec-dev-in-action-src/scanner/password_crack/models"
)

type ScanFunc func(service models.Service) (result models.ScanResult, err error)

var (
	ScanFuncMap map[string]ScanFunc
)

func init() {
	ScanFuncMap = make(map[string]ScanFunc)
	ScanFuncMap["FTP"] = ScanFtp
	ScanFuncMap["SSH"] = ScanSsh
	ScanFuncMap["MYSQL"] = ScanMysql
	ScanFuncMap["MSSQL"] = ScanMssql
	ScanFuncMap["REDIS"] = ScanRedis
	ScanFuncMap["POSTGRESQL"] = ScanPostgres
	ScanFuncMap["MONGODB"] = ScanMongodb
}

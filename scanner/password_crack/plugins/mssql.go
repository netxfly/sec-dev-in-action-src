package plugins

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"

	"sec-dev-in-action-src/scanner/password_crack/models"
)

func ScanMssql(service models.Service) (result models.ScanResult, err error) {
	result.Service = service

	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v", service.Ip,
		service.Port, service.Username, service.Password, "master")

	db, err := sql.Open("mssql", dataSourceName)
	if err == nil {
		defer func() {
			err = db.Close()
		}()

		err = db.Ping()
		if err == nil {
			result.Result = true
		}
	}

	return result, err
}

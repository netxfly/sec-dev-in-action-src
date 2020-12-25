package plugins

import (
	"fmt"

	"github.com/jlaffaye/ftp"

	"sec-dev-in-action-src/scanner/password_crack/models"
	"sec-dev-in-action-src/scanner/password_crack/vars"
)

func ScanFtp(s models.Service) (result models.ScanResult, err error) {
	result.Service = s
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", s.Ip, s.Port), vars.TimeOut)
	if err == nil {
		err = conn.Login(s.Username, s.Password)
		if err == nil {
			defer func() {
				err = conn.Logout()
			}()
			result.Result = true
		}
	}
	return result, err
}

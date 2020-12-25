package plugins

import (
	"fmt"

	"gopkg.in/mgo.v2"

	"sec-dev-in-action-src/scanner/password_crack/models"
	"sec-dev-in-action-src/scanner/password_crack/vars"
)

func ScanMongodb(s models.Service) (result models.ScanResult, err error) {
	result.Service = s
	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", s.Username, s.Password, s.Ip, s.Port, "test")
	session, err := mgo.DialWithTimeout(url, vars.TimeOut)

	if err == nil {
		defer session.Close()
		err = session.Ping()
		if err == nil {
			result.Result = true
		}
	}

	return result, err
}

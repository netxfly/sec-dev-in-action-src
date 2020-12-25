package models

import (
	"sec-dev-in-action-src/waf/x-waf-admin/modules/util"

	"log"
)

type User struct {
	Id       int64
	UserName string `xrom:"unique"`
	Password string `xrom:"varchar(50) notnull"`
}

func ListUser() (users []User, err error) {
	users = make([]User, 0)
	err = Engine.Find(&users)
	log.Println(err, users)
	return users, err
}

func NewUser(username, password string) (err error) {
	encryptpass := util.EncryptPass(password)
	_, err = Engine.Insert(&User{UserName: username, Password: encryptpass})
	return err
}

func UpdateUser(id int64, username, password string) (err error) {
	user := new(User)
	_, err = Engine.Id(id).Get(user)
	user.UserName = username
	user.Password = util.EncryptPass(password)
	_, err = Engine.Id(id).Update(user)

	return err
}

func DelUser(id int64) (err error) {
	_, err = Engine.Delete(&User{Id: id})
	return err
}

func Auth(username, password string) (ret bool, err error) {
	encryptpass := util.EncryptPass(password)
	log.Println(username, password, encryptpass)
	userAuth := &User{UserName: username, Password: encryptpass}
	has, err := Engine.Get(userAuth)
	log.Println(has, err, userAuth)

	if has {
		ret = true
	}
	return ret, err
}

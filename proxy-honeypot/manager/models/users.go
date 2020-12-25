/*

Copyright (c) 2018 sec.lu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package models

import (
	"gopkg.in/mgo.v2/bson"

	"sec-dev-in-action-src/proxy-honeypot/manager/util"
)

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	UserName string
	Password string
}

func ListUser() (users []User, err error) {
	err = collAdmin.Find(nil).All(&users)
	return users, err
}

func GetUserById(id string) (user User, err error) {
	err = collAdmin.FindId(bson.ObjectIdHex(id)).One(&user)
	//log.Println(id, user, err)
	return user, err

}

func NewUser(username, password string) (err error) {
	encryptPass := util.EncryptPass(password)
	err = collAdmin.Insert(&User{Id: bson.NewObjectId(), UserName: username, Password: encryptPass})
	return err
}

func UpdateUser(id string, username, password string) (err error) {
	user := new(User)
	err = collAdmin.FindId(bson.ObjectIdHex(id)).One(user)
	user.UserName = username
	user.Password = util.EncryptPass(password)
	err = collAdmin.UpdateId(bson.ObjectIdHex(id), user)
	return err
}

func DelUser(id string) (err error) {
	err = collAdmin.RemoveId(bson.ObjectIdHex(id))
	return err
}

func Auth(username, password string) (result bool, err error) {
	encryptPass := util.EncryptPass(password)
	userAuth := User{}
	err = collAdmin.Find(bson.M{"username": username, "password": encryptPass}).One(&userAuth)
	if err == nil && userAuth.UserName == username && userAuth.Password == encryptPass {
		result = true
	}
	return result, err
}
